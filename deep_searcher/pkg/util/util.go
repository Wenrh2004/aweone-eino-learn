package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
)

func DeduplicateResults(documents []*schema.Document) []*schema.Document {
	// Deduplicate the documents based on their content and keep only content
	uniqueDocs := make(map[string]*schema.Document)
	for _, doc := range documents {
		if _, exists := uniqueDocs[doc.Content]; !exists {
			uniqueDocs[doc.Content] = doc
		}
	}
	// Convert the map back to a slice
	uniqueDocsSlice := make([]*schema.Document, 0, len(uniqueDocs))
	for _, doc := range uniqueDocs {
		uniqueDocsSlice = append(uniqueDocsSlice, doc)
	}
	return uniqueDocsSlice
}

// LiteralEval 解析字符串响应为Go对象
func LiteralEval(responseContent string) ([]interface{}, error) {
	// 去除字符串前后的空白字符
	responseContent = strings.TrimSpace(responseContent)

	// 移除<think>和</think>之间的内容
	if strings.Contains(responseContent, "<think>") && strings.Contains(responseContent, "</think>") {
		thinkEnd := strings.Index(responseContent, "</think>") + len("</think>")
		if thinkEnd > 0 && thinkEnd <= len(responseContent) {
			responseContent = responseContent[thinkEnd:]
			responseContent = strings.TrimSpace(responseContent)
		}
	}

	// 处理代码块格式
	if strings.HasPrefix(responseContent, "```") && strings.HasSuffix(responseContent, "```") {
		if strings.HasPrefix(responseContent, "```json") {
			responseContent = responseContent[7 : len(responseContent)-3]
		} else {
			return nil, errors.New("invalid code block format")
		}
		responseContent = strings.TrimSpace(responseContent)
	}

	// 尝试解析为JSON对象
	var result []interface{}
	err := sonic.Unmarshal([]byte(responseContent), &result)
	if err == nil {
		return result, nil
	}

	// 如果解析失败，尝试提取JSON结构
	jsonContent, err := extractJSONStructure(responseContent)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonContent), &result)
	if err != nil {
		return nil, errors.New("parsing the extracted content failed: " + err.Error() + "and the content is: " + jsonContent)
	}

	return result, nil
}

// extractJSONStructure 从字符串中提取第一个有效的JSON结构（列表或对象）
func extractJSONStructure(content string) (string, error) {
	firstChar := -1
	for i, c := range content {
		if c == '[' || c == '{' {
			firstChar = i
			break
		}
	}

	if firstChar == -1 {
		return "", errors.New("响应中未找到有效的JSON结构")
	}

	startChar := content[firstChar]
	endChar := byte(']')
	if startChar == '{' {
		endChar = '}'
	}

	// 使用栈匹配括号，同时处理字符串中的括号
	stack := 1
	inString := false
	escaped := false

	for i := firstChar + 1; i < len(content); i++ {
		c := content[i]

		if escaped {
			escaped = false
			continue
		}

		if c == '\\' {
			escaped = true
			continue
		}

		if c == '"' {
			inString = !inString
			continue
		}

		if !inString {
			if c == startChar {
				stack++
			} else if c == endChar {
				stack--
				if stack == 0 {
					return content[firstChar : i+1], nil
				}
			}
		}
	}

	return "", errors.New("no matching closing parentheses were found")
}

func GetSupportedDocuments(retrievedResults []*schema.Document, supportedDocIndices []interface{}) []*schema.Document {
	var supportedRetrievedResults []*schema.Document

	for _, indexVal := range supportedDocIndices {
		// 将索引转换为整数
		var index int
		switch v := indexVal.(type) {
		case float64:
			index = int(v)
		case int:
			index = v
		case int64:
			index = int(v)
		default:
			// 如果无法转换为整数，跳过此索引
			continue
		}

		// 检查索引是否有效
		if index >= 0 && index < len(retrievedResults) {
			supportedRetrievedResults = append(supportedRetrievedResults, retrievedResults[index])
		}
	}

	return supportedRetrievedResults
}

func IntermediateAnswereModifier(_ context.Context, input map[string]any) (output map[string]any, err error) {
	followupQuery, ok := input["sub_query"]
	if !ok {
		return nil, errors.New("sub_query is not found")
	}
	subQuery, ok := followupQuery.(*schema.Message)
	if !ok {
		return nil, errors.New("sub_query is not a Message")
	}
	answer, ok := input["answer"]
	if !ok {
		return nil, errors.New("answer is not found")
	}
	intermediateAnswer, ok := answer.(*schema.Message)
	if !ok {
		return nil, errors.New("answer is not a Message")
	}
	answerChecker, err := intermediateAnswerChecker(intermediateAnswer)
	if err != nil {
		return nil, err
	}
	output["intermediate_contexts"] = fmt.Sprintf("Intermediate query: %s\nIntermediate answer: %s", subQuery.Content, intermediateAnswer.Content)
	output["answer"] = intermediateAnswer
	output["check_status"] = answerChecker
	return output, nil
}

func intermediateAnswerChecker(input *schema.Message) (output *schema.Message, err error) {
	if strings.Contains(input.Content, "No relevant information found") {
		return nil, nil
	}
	return input, nil
}

func GetCollateOutput(_ context.Context, input map[string]any) (output map[string]any, err error) {
	retrievedResults, ok := input["all_retrieved_results"]
	if !ok {
		return nil, errors.New("all_retrieved_results is not found")
	}
	documents, ok := retrievedResults.([]*schema.Document)
	if !ok {
		return nil, errors.New("all_retrieved_results is not match")
	}
	intermediateContexts, ok := input["intermediate_contexts"]
	if !ok {
		return nil, errors.New("intermediate_contexts is not found")
	}
	contexts, ok := intermediateContexts.([]string)
	if !ok {
		return nil, errors.New("intermediate_contexts is not match")
	}
	return map[string]any{
		"all_retrieved_results": documents,
		"additional_info":       contexts,
	}, nil
}

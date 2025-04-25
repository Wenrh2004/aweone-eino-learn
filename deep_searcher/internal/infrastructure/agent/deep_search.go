package agent
const (
	SubQueryPrompt = `To answer this question more comprehensively, please break down the original question into up to four sub-questions. Return as list of str.
If this is a very simple question and no decomposition is necessary, then keep the only one original question in the python code list.

Original Question: {original_query}


<EXAMPLE>
Example input:
"Explain deep learning"

Example output:
[
"What is deep learning?",
"What is the difference between deep learning and machine learning?",
"What is the history of deep learning?"
]
</EXAMPLE>

Provide your response in a python code list of str format:`
	
	RerankPrompt = `Based on the query questions and the retrieved chunk, to determine whether the chunk is helpful in answering any of the query question, you can only return "YES" or "NO", without any other information.

Query Questions: {query}
Retrieved Chunk: {retrieved_chunk}

Is the chunk helpful in answering the any of the questions?`
	
	ReflectPrompt = `Determine whether additional search queries are needed based on the original query, previous sub queries, and all retrieved document chunks. If further research is required, provide a Python list of up to 3 search queries. If no further research is required, return an empty list.

If the original query is to write a report, then you prefer to generate some further queries, instead return an empty list.

Original Query: {question}

Previous Sub Queries: {mini_questions}

Related Chunks:
{mini_chunk_str}

Respond exclusively in valid List of str format without any other text.`
	
	SummaryPrompt = `You are a AI content analysis expert, good at summarizing content. Please summarize a specific and detailed answer or report based on the previous queries and the retrieved document chunks.

Original Query: {question}

Previous Sub Queries: {mini_questions}

Related Chunks:
{mini_chunk_str}`
)
chatTemplate := prompt.FromMessages(schema.FString,
&schema.Message{
Role:    schema.User,
Content: SubQueryPrompt,
})

variables := map[string]any{
"original_query": originalQuery,
}

messages, err := chatTemplate.Format(ctx, variables)
if err != nil {
return err
}

package graph

// prompt const for chain of rag
var (
	FollowupQueryPrompt = `You are using a search tool to answer the main query by iteratively searching the database. Given the following intermediate queries and answers, generate a new simple follow-up question that can help answer the main query. You may rephrase or decompose the main query when previous answers are not helpful. Ask simple follow-up questions only as the search tool may not understand complex questions.

## Previous intermediate queries and answers
{intermediate_contexts}

## Main query to answer
{query}

Respond with a simple follow-up question that will help answer the main query, do not explain yourself or output anything else.
`

	IntermediateAnswerPrompt = `Given the following documents, generate an appropriate answer for the query. DO NOT hallucinate any information, only use the provided documents to generate the answer. Respond "No relevant information found" if the documents do not contain useful information.

## Documents
{retrieved_documents}

## Query
{sub_query}

Respond with a concise answer only, do not explain yourself or output anything else.`

	ReflectionPrompt = `Given the following intermediate queries and answers, judge whether you have enough information to answer the main query. If you believe you have enough information, respond with "Yes", otherwise respond with "No".

## Intermediate queries and answers
{intermediate_contexts}

## Main query
{query}

Respond with "Yes" or "No" only, do not explain yourself or output anything else.`

	GetSupportedDocsPrompt = `Given the following documents, select the ones that are support the Q-A pair.

## Documents
{retrieved_documents}

## Q-A Pair
### Question
{sub_query}
### Answer
{answer}

Respond with a json array of indices of the selected documents.`
)

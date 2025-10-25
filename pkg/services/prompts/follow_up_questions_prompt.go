package prompts

const FollowUpQuestionsPrompt = `
You are an expert in investing! Your mission is given a conversation between a user and an AI assistant about investing to respond
with %d follow up questions that the user can ask given the context of the conversation.

## CONVERSATION
%+v

## RESPONSE FORMAT
- Your response MUST BE a json parsable string with a key named 'follow_up_questions' and value an array of strings that will contain
the follow up questions.

Example response:
{
	"follow_up_questions":[
		"follow up question",
		"follow up question",
		"follow up question"
	]
}
`

package prompts

const FollowUpQuestionsPrompt = `
You are an expert in investing! Your mission is given a conversation about investing to respond
with %d follow up questions that the user can ask ask given the context of the conversation.

## CONVERSATION
%+v

## RESPONSE FORMAT
Your answer MUST have the following format:
<Follow up question>
<Follow up question>
.
.
<Follow up question>
`

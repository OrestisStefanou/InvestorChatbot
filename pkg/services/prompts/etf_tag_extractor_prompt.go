package prompts

const EtfTagExtractorPrompt = `
Given a conversation about etfs your mission is to understand for which etf symbol the conversation is about.

## Etf names and symbols
%+v

Some context of the user asking the question is given below. You should take this into consideration.
## User context
%+v

## Response instructions
- Your response MUST BE a json parsable string with a key named 'etf_symbols' and value an array of strings that will contain
the etf symbols the conversation is about. In case the question is generic and not for a specific etf then return an empty 
array for 'etf_symbols' key.

- Your answer should focus on the last question of the conversation, for example if the first 5 messages are about Vanguard S&P 500 ETF
but the last question is about Vanguard Growth ETF then your answer should contain the Vanguard Growth ETF symbol.

Example response if the conversation is about Vanguard S&P 500 ETF:
{"etf_symbols":["VOO"]}

Example response if the conversation is about Vanguard S&P 500 ETF and Vanguard Growth ETF:
{"etf_symbols":["VOO", "VUG"]}

If the conversation is about the user portfolio then you should use the user context above for your response.

# Conversation
%+v
`

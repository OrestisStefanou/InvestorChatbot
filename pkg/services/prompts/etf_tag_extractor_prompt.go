package prompts

const EtfTagExtractorPrompt = `
Given a conversation about etfs your mission is to understand for which etf symbol the conversation is about.

## Etf names and symbols
%+v

## Response instructions
- Your response MUST BE a json parsable string with a key named 'etf_symbol' and value the string that will contain
the etf symbol the conversation is about. In case the question is generic and not for a specific etf then return an empty 
string for 'etf_symbol' key.

- Your answer should focus on the last question of the conversation, for example if the first 5 messages are about Vanguard S&P 500 ETF
but the last question is about Vanguard Growth ETF then your answer should contain the Vanguard Growth ETF symbol.

Example response:
{"etf_symbol":["VOO"]}

# Conversation
%+v
`

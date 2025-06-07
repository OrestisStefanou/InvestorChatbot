package prompts

const StockOverviewTagExtractorPrompt = `
Given a conversation about stocks your mission is to understand for which stock symbols the conversation is about.

## Stock names and symbols
%+v

## Response instructions
- Your response MUST BE a json parsable string with a key named 'stock_symbols' and value an array of strings that will contain
the stock symbols the conversation is about. In case the question is generic and not for a specific stock then return an empty 
array for 'stock_symbols' key.

For example if the conversation is about the microsoft stock the response should look like this:
{"stock_symbols":["MSFT"]}

If the conversation is about Microsoft and Apple then the response should look like this:
{"stock_symbols":["MSFT", "AAPL"]}

# Conversation
%+v
`

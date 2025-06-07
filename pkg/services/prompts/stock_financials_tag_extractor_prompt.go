package prompts

const StockFinancialsTagExtractorPrompt = `
Given a conversation about stocks financials your mission is to to extract the following information:
1. Which stock symbols the conversation is about.
2. If the conversation is about balance sheets.
3. If the conversation is about cash flows.
4. If the conversation is about income statements.

## Stock names and symbols
%+v

## Response instructions
Your response MUST BE a json parsable string with the following keys:
- stock_symbols: A list of strings that will contain the symbols of the stocks that the conversation is about.
In case the question is generic and not for a specific stock then return an empty array for 'stock_symbols' key.
- balance_sheet: A boolean that must be set to true in case the conversation is about balance sheets and false otherwise.
- income_statement: A boolean that must be set to true in case the conversation is about income statements and false otherwise.
- cash_flow: A boolean that must be set to true in case the conversation is about cash flows and false otherwise.

For example if the conversation is about Microsoft balance sheets the response should look like this:
{"stock_symbols":["MSFT"], "balance_sheet": true, "income_statement": false, "cash_flow": false}

If the conversation is about Microsoft and Apple income statements and cash flows then the response should look like this:
{"stock_symbols":["MSFT", "AAPL"], "balance_sheet": false, "income_statement": true, "cash_flow": true}

# Conversation
%+v
`

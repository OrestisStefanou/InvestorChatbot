package prompts

const TopicExtractorPrompt = `
# Objective
Given a conversation about investing your mission is to categorize it into ONE of the following topics:
- education
- sectors
- stock_overview
- stock_financials
- etfs
- news

Below are some examples for each topic:
## education
- What are the differences between short-term and long-term investing?
- What is the difference between active and passive investing?
- What are index funds, and why do investors use them?
- What is diversification, and why is it important?

## sectors
- What are the main stock market sectors, and how do they differ?
- Which are the top performing sectors?
- How should a beginner decide which sectors to invest in?
- What are some good defensive sectors for long-term investors?

## stock_overview
- What are the key financial ratios for this stock, and what do they indicate?
- Can you give me a high-level overview of this stock?
- Can you do a quick valuation of the stock?
- What does the current ratio tell me about a company's ability to pay short-term liabilities?
- What is the current P/E ratio of this stock, and what does it tell me?

## stock_financials
- Can you give me an overview of the cash flow?
- Can you give me an overview of the income statement?
- Can you give me an overview of the balance sheet?

## etfs
- Which are the different ETF asset classes and what is the difference between them?
- What are the key advantages of investing in ETFs?
- How do ETFs provide diversification benefits to investors?

## news
- What are the latest market news?
- What are the latest news of Apple stock?

## General guidance on how to choose a topic
- education: Anything that has to do with investing education falls under this topic
- sectors: Anything that is related to stock sectors falls under this topic
- stock_overview: Anything that is related to a stock but is not specifically about balance sheets, income statemets or 
cash flows falls under this category
- stock_financials: If the conversation is specifically about income statement or cash flow or balance sheet then it falls under this category
- etfs: Anything that is related to ETFs falls under this category
- news: Anything that is related to market or stock news falls under this category

# Response instructions
Your response MUST BE just the topic and nothing else.

# Conversation to categorize
%+v 
`

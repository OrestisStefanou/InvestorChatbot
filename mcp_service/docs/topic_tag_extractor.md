# How topic and tag extraction works
Topic and tag extraction endpoint is broken down into two steps
1. Use an llm to identify the topic of the conversation. You can look at the prompt that is used for this task in `pkg/services/prompts/topic_extractor_prompt.go` file.
2. Based on the topic that we extracted from step 1 use another llm to extract the tags. For example if the topic extracted was education then
there is no need to make a second llm call since education topic needs no tags. If the topic extracted was stock_overview then we use a second 
llm to extract the stock symbols from the conversation. You can find the prompt here `pkg/services/prompts/stock_overview_tag_extractor_prompt.go`.
Each topic that needs tag extraction has it's own prompt.

## Advantages of this approach
1. Having smaller and simpler tasks(one for topic extraction and one for tag extraction) improves the performance of the llm.
2. Tag extraction for each topic has different complexity(for example extracting the sector name is easier than extracting the 
stock symbols) so having a separate prompt for each topic can improve the performance of the llm.

## Disadvantages of this approach
1. More expensive since we have to make two llm calls. Especially for topics that we have to extract the symbols because we pass all the 
available symbols in the prompt.
2. If topic extraction is wrong then the second call for tag extraction is basically a waste of money.


### What can be done better
- Instead of passing all the symbols in the prompt find a smarter way to limit the symbols that are needed. One possible solution for this 
is to use embeddings and vector similarity. 
    1. Store the symbols along with company or etf name in a vector database using embeddings. Keep in mind that we have to store the actual text 
    in the database as well since we will need it to pass it the prompt.
    2. Transform user's query to a vector using embeddings.
    3. Perform a distance query on the database and retrieve the symbols that are closer to the user's query.
    4. Include only those in the prompt. 
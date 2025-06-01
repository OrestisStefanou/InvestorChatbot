package prompts

const SectorTagExtractorPrompt = `
# Objective
Given a conversation about stock sectors your mission is to understand which sector the conversation is about.

## Sectors
%s

## Response instructions
- Focus on the last question of the conversation, for example in the first messages are about about the technology sector 
but the last question is about the energy sector then in your response you should have the energy sector.
- Your response MUST BE a json parsable string with a key named 'sector_name' and value one of the sectors above. In case the question is 
sector generic and not for a specific sector then return an empty string as a value for the 'sector' key.

# Conversation
%+v
`

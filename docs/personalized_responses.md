# How to personalize chatbot responses

## Goal
Make the responses of the chatbot personalized based on the user profile. For example the chatbot should be able to take into consideration the following:
- User's portfolio
- User profile
    - Risk appetite
    - Age
    - Investing goal
    - Name?

## How will this work
- Expose an endpoint POST /user-context where the caller will be able to add user information to the backend of this service
- The body of the request will look like this:
{
    "user_id" : "some_user_id",  # must be given
    "user_profile": {   # caller can pass whatever information they want here and it will bne passed as it to the llm
        "name" : "Orestis",
        "age": 27,
        "risk_apettite": "medium"
    },
    "user_portfolio": [ # Optional
        {
            "symbol": "NVDA",   # if name is not given this must not be empty
            "name": "Nvidia",   # if symbol is not given this must not be empty
            "asset_class": "stock", # Must be given?
            "portfolio_percentage": 0.5 # Optional
            "quantity" : 10 # Optional
        }
    ]
}

- User context will be passed in the prompt so that the llm can take it into consideration. The existing chat endpoint should be updated so that the user_id can be passed in the body of the request.

- User context should also be passed in the topic and tag extraction prompt so that endpoint should also be updated to accept the user_id in the body of the request. 
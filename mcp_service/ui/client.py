import streamlit as st
import requests
import json
import random

BASE_URL = "http://localhost:1323"
USER_ID = ""    # Add user id here


def create_session() -> str:
    response = requests.post(f"{BASE_URL}/session")
    session_id = response.json()["session_id"]
    return session_id


def extract_topic_and_tags(session_id: str, question: str) -> dict:
    payload = {
        "session_id": session_id,
        "question": question,
        "user_id": USER_ID,
    }

    response = requests.post(f"{BASE_URL}/chat/extract_topic_and_tags", json=payload)
    if response.status_code != 200:
        print(response.json())
        raise Exception("Failed to extract topic and tags")
    
    return response.json()


def fetch_follow_ups(session_id, number_of_questions=5):
    try:
        payload = {"session_id": session_id, "number_of_questions": number_of_questions}
        response = requests.post(f"{BASE_URL}/follow_up_questions", json=payload)
        return response.json().get("follow_up_questions", [])
    except Exception as e:
        return [f"Error fetching follow-up questions: {str(e)}"]


def fetch_faqs(topic: str, topic_tags: dict) -> list[str]:
    topic_to_faq_topic = {
        "education": "education",
        "sectors": "sectors",
        "stock_overview": "stock_overview",
        "etfs": "etfs",
        "stock_financials": "stock_overview",
    }
    
    faq_topic = topic_to_faq_topic.get(topic)

    if topic == "stock_financials":
        if topic_tags["balance_sheet"]:
            faq_topic = "balance_sheet"
        if topic_tags["income_statement"]:
            faq_topic = "income_statement"
        if topic_tags["cash_flow"]:
            faq_topic = "cash_flow"

    if not faq_topic:
        return []

    try:
        response = requests.get(f"{BASE_URL}/faq", params={"faq_topic": faq_topic})
        return response.json().get("faq", [])
    except Exception as e:
        return [f"Error fetching FAQs: {str(e)}"]


def get_prompt_placeholder() -> str:
    last_topic = "education"    # Default value in case this is the very first message
    if "last_topic" in st.session_state:
        last_topic = st.session_state["last_topic"]
    
    last_topic_tags = {}
    if "last_topic_tags" in st.session_state:
        last_topic_tags = st.session_state["last_topic_tags"]
    
    faqs = fetch_faqs(last_topic, last_topic_tags)
    if not faqs:
        return ""
    
    # Return a random question from the faqs
    return random.choice(faqs) 


def stream_chat_response(
    question: str,
    session_id: str, 
    topic: str,
    topic_tags: dict=None,
):
    url = f"{BASE_URL}/chat"
    payload = {
        "question": question,
        "topic": topic,
        "session_id": session_id
    }
    
    if topic_tags:
        payload["topic_tags"] = topic_tags

    headers = {"Content-Type": "application/json"}

    try:
        response = requests.post(url, json=payload, headers=headers, stream=True)

        response.raise_for_status()  # Raise error for non-2xx responses

        for chunk in response.iter_lines(decode_unicode=True):
            if chunk:  # filter out keep-alive new lines
                yield chunk

    except requests.exceptions.RequestException as e:
        yield f"[ERROR] {e}"


def response_generator(session_id: str, question: str):
    topic_and_tags = extract_topic_and_tags(
        session_id=session_id,
        question=question,
    )

    for chunk in stream_chat_response(
        question=question,
        session_id=session_id,
        topic=topic_and_tags["topic"],
        topic_tags=topic_and_tags["topic_tags"],
    ):
        try:
            # Parse each chunk from JSON string literal to clean text
            text = json.loads(chunk)
            if text.strip():  # skip empty chunks
                yield text
        except Exception as e:
            print(f"Chunk parse error: {e} | raw: {repr(chunk)}")
            continue

    st.session_state.last_topic = topic_and_tags["topic"]
    st.session_state.last_topic_tags = topic_and_tags["topic_tags"]

    # Update prompt placeholder based on the latest topic and tags
    st.session_state.prompt_placeholder = get_prompt_placeholder()



st.title("Stock Analysis chatbot")

# Initialize chat history
if "messages" not in st.session_state:
    st.session_state.messages = []

# Initialize session
if "session_id" not in st.session_state:
    st.session_state.session_id = create_session()

# Display chat messages from history on app rerun
for message in st.session_state.messages:
    with st.chat_message(message["role"]):
        st.markdown(message["content"])


# Set the prompt placeholder once per session
if "prompt_placeholder" not in st.session_state:
    st.session_state.prompt_placeholder = get_prompt_placeholder()


# Get user input or use follow-up as prompt
prompt = st.chat_input()

if prompt:
    # Add user message to chat history
    st.session_state.messages.append({"role": "user", "content": prompt})
    with st.chat_message("user"):
        st.markdown(prompt)

    # Display assistant response in chat message container
    with st.chat_message("assistant"):
        response_placeholder = st.empty()
        full_response = ""

        for chunk in response_generator(
            session_id=st.session_state.session_id, question=prompt
        ):
            full_response += chunk
            response_placeholder.markdown(full_response)

    # Add assistant's reply to history
    st.session_state.messages.append({"role": "assistant", "content": full_response})

    # Show follow-up questions
    follow_ups = fetch_follow_ups(st.session_state.session_id)
    if follow_ups:
        st.markdown("#### Follow-up Questions")
        for i, follow_up in enumerate(follow_ups):
            st.markdown(f"- {follow_up}")

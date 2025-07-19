import streamlit as st
import requests
import json


BASE_URL = "http://localhost:1323"
USER_ID = "orestis_user_id"


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


st.title("Investor Assistant chatbot")

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


# Get user input or use follow-up as prompt
prompt = st.chat_input("What is up?")

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

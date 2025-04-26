import gradio as gr
import requests
import json

BASE_URL = "http://localhost:1323"  # Change this to your actual API server if needed

# ========== Backend API functions ==========

def fetch_topics():
    try:
        response = requests.get(f"{BASE_URL}/topics")
        return response.json().get("topics", [])
    except:
        return []


def create_session():
    try:
        response = requests.post(f"{BASE_URL}/session")
        return response.json().get("session_id", "")
    except:
        return ""


def fetch_faq(topic, include_bs, include_is, include_cf):
    topic_to_faq_topic = {
        "education": "education",
        "sectors": "sectors",
        "stock_overview": "stock_overview",
        "etfs": "etfs",
        "stock_financials": "stock_overview",
    }
    
    faq_topic = topic_to_faq_topic.get(topic)

    if topic == "stock_financials":
        if include_bs:
            faq_topic = "balance_sheet"
        if include_is:
            faq_topic = "income_statement"
        if include_cf:
            faq_topic = "cash_flow"

    if not faq_topic:
        return []

    print(f"Fetching FAQs for topic: {faq_topic}")
    try:
        response = requests.get(f"{BASE_URL}/faq", params={"faq_topic": faq_topic})
        return response.json().get("faq", [])
    except Exception as e:
        return [f"Error fetching FAQs: {str(e)}"]


def fetch_follow_ups(session_id, number_of_questions=5):
    try:
        payload = {"session_id": session_id, "number_of_questions": number_of_questions}
        response = requests.post(f"{BASE_URL}/follow_up_questions", json=payload)
        return response.json().get("follow_up_questions", [])
    except Exception as e:
        return [f"Error fetching follow-up questions: {str(e)}"]


# ========== Chat + Follow Up Logic ==========
def ask_question(question, topic, sector, industry, tickers, include_bs, include_is, include_cf, etf):
    session_id = create_session()
    topic_tags = {}

    if sector:
        topic_tags["sector_name"] = sector
    if industry:
        topic_tags["industry_name"] = industry
    if tickers:
        topic_tags["stock_symbols"] = [s.strip().upper() for s in tickers.split(",")]
    if include_bs:
        topic_tags["balance_sheet"] = True
    if include_is:
        topic_tags["income_statement"] = True
    if include_cf:
        topic_tags["cash_flow"] = True
    if etf:
        topic_tags["etf_symbol"] = etf.upper()

    payload = {
        "question": question,
        "topic": topic,
        "session_id": session_id,
        "topic_tags": topic_tags
    }

    buffer = ""
    try:
        response = requests.post(f"{BASE_URL}/chat", json=payload, stream=True)
        if response.status_code == 200:
            for chunk in response.iter_lines():
                if chunk:
                    try:
                        part = json.loads(chunk.decode("utf-8"))
                    except json.JSONDecodeError:
                        part = chunk.decode("utf-8")
                    part = part.replace('""', '"')
                    buffer += part
                    yield buffer, gr.update(value="")  # Empty FAQs while answering
        else:
            yield f"Error: {response.status_code} - {response.json().get('error', 'Unknown error')}", gr.update(value="")
    except Exception as e:
        yield str(e), gr.update(value="")

    # After chat, fetch follow-ups
    follow_ups = fetch_follow_ups(session_id)
    if follow_ups:
        follow_ups_text = "### Follow-Up Questions:\n" + "\n".join([f"- {q}" for q in follow_ups])
    else:
        follow_ups_text = ""
    yield buffer, follow_ups_text


# ========== FAQ update logic ==========
def update_faq(selected_topic, include_bs, include_is, include_cf):
    if selected_topic:
        faqs = fetch_faq(selected_topic, include_bs, include_is, include_cf)
        faq_text = "### FAQs\n" + "\n".join([f"- {q}" for q in faqs])
        return faq_text
    return ""


# ========== Build Gradio UI ==========
topics = fetch_topics()

with gr.Blocks() as demo:
    gr.Markdown("## 💬 Financial Chat Assistant")

    with gr.Row():
        topic_dropdown = gr.Dropdown(choices=topics, label="Select Topic", interactive=True)
        question_input = gr.Textbox(label="Your Question")

    faq_box = gr.Markdown()

    with gr.Accordion("Optional Tags", open=False):
        sector_input = gr.Textbox(label="Sector (e.g. Technology)")
        industry_input = gr.Textbox(label="Industry (e.g. Semiconductors)")
        tickers_input = gr.Textbox(label="Stock Symbols (comma-separated, e.g. AAPL, MSFT)")
        etf_input = gr.Textbox(label="ETF Symbol (e.g. SPY)")
        include_bs = gr.Checkbox(label="Include Balance Sheet")
        include_is = gr.Checkbox(label="Include Income Statement")
        include_cf = gr.Checkbox(label="Include Cash Flow")

    submit_btn = gr.Button("Ask")

    output_box = gr.Markdown()

    # Update FAQ on topic selection
    topic_dropdown.change(
        fn=update_faq,
        inputs=[topic_dropdown, include_bs, include_is, include_cf],
        outputs=faq_box
    )

    # On Submit: Ask + Get Response + Follow-Ups
    submit_btn.click(
        fn=ask_question,
        inputs=[
            question_input, topic_dropdown, sector_input, industry_input,
            tickers_input, include_bs, include_is, include_cf, etf_input
        ],
        outputs=[output_box, faq_box],
    )

demo.launch()

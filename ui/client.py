import gradio as gr
import requests
import json
import re

BASE_URL = "http://localhost:1323"  # Change to your API base URL

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

    try:
        response = requests.post(f"{BASE_URL}/chat", json=payload, stream=True)
        if response.status_code == 200:
            buffer = ""
            for chunk in response.iter_lines():
                if chunk:
                    try:
                        part = json.loads(chunk.decode("utf-8"))
                    except json.JSONDecodeError:
                        part = chunk.decode("utf-8")
                    part = part.replace('""', '"')
                    buffer += part
                    yield buffer
        else:
            yield f"Error: {response.status_code} - {response.json().get('error', 'Unknown error')}"
    except Exception as e:
        yield str(e)

# Initialize topics
topics = fetch_topics()

# Gradio UI
with gr.Blocks() as demo:
    gr.Markdown("## 💬 Financial Chat Assistant")

    with gr.Row():
        topic_dropdown = gr.Dropdown(choices=topics, label="Select Topic", interactive=True)
        question_input = gr.Textbox(label="Your Question")

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

    submit_btn.click(
        fn=ask_question,
        inputs=[
            question_input, topic_dropdown, sector_input, industry_input,
            tickers_input, include_bs, include_is, include_cf, etf_input
        ],
        outputs=output_box
    )

demo.launch()

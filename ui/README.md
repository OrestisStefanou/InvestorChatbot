### 📘 `README.md`

# 🧪 Demo UI

This is a **demo UI built with Gradio** for testing and interacting with the backend in this project. It provides a simple frontend to send questions, select topics, and optionally include financial context.

---

## ⚙️ Dependencies

Make sure you have Python 3.10+ installed, then install required packages:

```bash
pip install gradio requests
```

---

## ▶️ Running the Demo

1. Ensure the backend API is running and accessible (default is `http://localhost:8000`).
2. Run the demo UI:

```bash
python financial_chat_ui.py
```

3. Open the link provided by Gradio, usually:

```
http://localhost:7860
```

---

## 💡 Notes

- This UI streams the chat response as it's generated and renders it as Markdown.
- Topics are dynamically fetched from the `/topics` endpoint.
- A new session is created automatically for each chat.
- Optional tags like `sector`, `industry`, stock `tickers`, and `ETF` symbol can be included to add context.

---
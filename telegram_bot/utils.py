import re

def split_ai_response_message(ai_response_message: str) -> list[str]:
    """
    Splits the AI response message into chunks to handle the character limit 
    of the Telegram API (4096 characters per message)
    """
    return ai_response_message.split("\n")


def markdown_to_telegram_markdown(text):
    """
    Convert basic Markdown to Telegram MarkdownV2 format.
    Handles: headers, bold, italic
    """

    # Step 1: Convert markdown syntax
    # Headers -> bold
    text = re.sub(r'^#{1,6}\s+(.+)$', r'*\1*', text, flags=re.MULTILINE)

    # Bold: **text** or __text__ -> *text*
    text = re.sub(r'\*\*(.+?)\*\*', r'*\1*', text)
    text = re.sub(r'__(.+?)__', r'*\1*', text)

    # Italic: *text* or _text_ -> _text_
    text = re.sub(r'(?<!\*)\*([^\*]+?)\*(?!\*)', r'_\1_', text)
    text = re.sub(r'(?<!_)_([^_]+?)_(?!_)', r'_\1_', text)

    # Step 2: Escape all special characters
    special_chars = ['_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!']

    for char in special_chars:
        text = text.replace(char, '\\' + char)

    # Step 3: Unescape formatting markers
    # Bold: \*text\* -> *text*
    text = re.sub(r'\\\*(.+?)\\\*', r'*\1*', text)

    # Italic: \_text\_ -> _text_
    text = re.sub(r'\\_(.+?)\\_', r'_\1_', text)

    return text

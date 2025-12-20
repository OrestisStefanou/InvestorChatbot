import logging

from telegram import Update
from telegram.ext import (
    Application, 
    CommandHandler, 
    MessageHandler, 
    filters, 
    ContextTypes
)

from config import settings

# Enable logging
logging.basicConfig(
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    level=logging.INFO
)
logger = logging.getLogger(__name__)

async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Send a message when the command /start is issued."""
    # TODO: Send an introduction message here
    # TODO: Create user context in agent service
    # TODO: Create a session in agent service
    user = update.effective_user
    logger.info(f"User {user.id} ({user.username}) sent /start")
    await update.message.reply_text(
        f'Hi {user.first_name}! I am your test bot. Send me any message and I will echo it back!'
    )

async def handle_message(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Echo the user message and log user info."""
    user = update.effective_user
    message_text = update.message.text
    
    # Log the message and user info
    logger.info(f"User {user.id} (@{user.username}) sent: {message_text}")
    
    # Respond to the user
    response = f"You said: {message_text}\n\nYour user ID: {user.id}"
    if user.username:
        response += f"\nYour username: @{user.username}"
    
    await update.message.reply_text(response)

async def error_handler(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Log errors caused by updates."""
    logger.error(f"Update {update} caused error {context.error}")

def main():
    application = Application.builder().token(settings.TELEGRAM_BOT_TOKEN).build()
    
    application.add_handler(CommandHandler("start", start))
    application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, handle_message))
    application.add_error_handler(error_handler)
    
    logger.info("Bot is starting with webhook...")
    application.run_webhook(
        listen="0.0.0.0",
        port=settings.TELEGRAM_WEBHOOK_PORT,
        url_path="webhook",
        webhook_url=f"{settings.TELEGRAM_WEBHOOK_URL}/webhook"
    )

if __name__ == '__main__':
    main()

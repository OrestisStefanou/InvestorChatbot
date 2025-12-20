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
from agent_service_client import AgentServiceClient

# Enable logging
logging.basicConfig(
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    level=logging.INFO
)
logger = logging.getLogger(__name__)

async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Send a message when the command /start is issued."""
    logger.info(f"User {update.effective_user.id} ({update.effective_user.username}) sent /start")
    agent_service_client = AgentServiceClient()
    
    user_id = f"telegram:{update.effective_user.id}"
    session_id = f"telegram_session:{update.effective_user.id}"
    try:
        await agent_service_client.create_user_context(
            user_id=user_id,
            user_profile={
                "first_name": update.effective_user.first_name,
            }
        )
    except Exception as e:
        logger.error(f"Failed to create user context: {e}")

    try:
        await agent_service_client.create_session(
            user_id=user_id,
            session_id=session_id,
        )
    except Exception as e:
        logger.error(f"Failed to create session: {e}")

    try:
        ai_message = await agent_service_client.generate_ai_response(
            session_id=session_id,
            message=f"Hey, I am your new client {update.effective_user.first_name}!",
        )
    except Exception as e:
        logger.error(f"Failed to generate AI response: {e}")
        ai_message = "I'm sorry, something went wrong. Please try again later."
    
    await update.message.reply_text(ai_message)


async def handle_message(update: Update, context: ContextTypes.DEFAULT_TYPE):
    agent_service_client = AgentServiceClient()
    session_id = f"telegram_session:{update.effective_user.id}"

    user = update.effective_user
    message_text = update.message.text
    
    # Log the message and user info
    logger.info(f"User {user.id} (@{user.username}) sent: {message_text}")
    
    try:
        ai_message = await agent_service_client.generate_ai_response(
            session_id=session_id,
            message=message_text,
        )
    except Exception as e:
        logger.error(f"Failed to generate AI response: {e}")
        ai_message = "I'm sorry, something went wrong. Please try again later."
    
    await update.message.reply_text(ai_message)

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

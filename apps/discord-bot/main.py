import discord
from discord.ext import commands
import logging
from dotenv import load_dotenv
import os
import pathlib
import asyncio
from typing import Optional
from discord.app_commands import CheckFailure, AppCommandError


# Set up logging configuration
log_dir: pathlib.Path = pathlib.Path(__file__).parent / 'logs'
handler: logging.FileHandler = logging.FileHandler(
    log_dir / 'discord-bot.log',
    encoding='utf-8',
    mode='w'
)
logging.basicConfig(
    level=logging.INFO,
    #level=logging.DEBUG,
    handlers=[handler],
)

# Load environment variables
load_dotenv()
token: Optional[str] = os.getenv('DISCORD_TOKEN')

# Configure bot intents
# Note: These must be enabled in both code and Discord Developer Portal
intents: discord.Intents = discord.Intents.default()
intents.message_content = True
intents.members = True

# Initialize bot with command prefix and intents
bot: commands.Bot = commands.Bot(command_prefix='!', intents=intents)


@bot.tree.error
async def on_app_command_error(interaction: discord.Interaction, error: AppCommandError) -> None:
    """
    Error handler for app slash commands such as invalid permissions or unexpected errors

    Args:
        interaction (discord.Interaction): _description_
        error (AppCommandError): _description_

    Raises:
        error: Error based on CheckFailure (invalid permissions) or unexpected errors
    """
    if isinstance(error, CheckFailure):
        await interaction.response.send_message(
            "🚫 You do not have permission to use this command.",
            ephemeral=True
        )
    else:
        await interaction.response.send_message(
            "⚠️ An unexpected error occurred. Please contact an admin.",
            ephemeral=True
        )
        raise error


@bot.event
async def on_ready() -> None:
    """Event triggered when the bot is ready and connected to Discord
    
    This event:
    1. Prints a ready message to console
    2. Logs the ready status to the log file
    """
    await bot.tree.sync()
    print(f"Bot: {bot.user} is ready to go.")
    logging.info(f"Bot: {bot.user} is ready to go.")


async def load_cogs() -> None:
    """Load all cog modules from the cogs directory
    
    This function:
    1. Finds all .py files in the cogs directory
    2. Attempts to load each as a cog
    3. Logs and prints success or failure for each cog
    
    Note:
        Cogs must be valid Python modules with a setup function
    """
    print("Starting to load cogs...")
    cogs_dir: pathlib.Path = pathlib.Path(__file__).parent / "cogs"
    print(f"Looking for cogs in: {cogs_dir}")
    
    if not cogs_dir.exists():
        print(f"ERROR: Cogs directory does not exist at {cogs_dir}")
        return
        
    files: list[str] = os.listdir(cogs_dir)
    print(f"Found files in cogs directory: {files}")
    
    for filename in files:
        if filename.endswith(".py"):
            try:
                await bot.load_extension(f"cogs.{filename[:-3]}")
                print(f"Successfully loaded cog: {filename[:-3]}")
                logging.info(f"Loaded cog: {filename[:-3]}")
            except Exception as e:
                print(f"Failed to load cog {filename[:-3]}: {str(e)}")
                logging.error(f"Failed to load cog {filename[:-3]}: {e}")
    print("Finished loading cogs")


async def main() -> None:
    """Main function for the bot
    
    This function calls functions to load cogs and start the bot given the api token from environment variables
    """
    await load_cogs()
    await bot.start(token)


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logging.info("Bot shutting down...")
        print("Bot shutting down...")

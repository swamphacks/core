import discord
from discord.ext import commands
import logging
from dotenv import load_dotenv
import os
import pathlib
import asyncio


# logging setup
log_dir = pathlib.Path(__file__).parent / 'logs'
handler = logging.FileHandler(log_dir / 'discord-bot.log', encoding='utf-8', mode='w')
logging.basicConfig(
        level=logging.DEBUG,
        handlers=[handler],
        )

# load config   
load_dotenv()
token = os.getenv('DISCORD_TOKEN')

# intents (all permissions are here, manually enable in discord developer portal and in code)
intents = discord.Intents.default()
intents.message_content = True
intents.members = True

# create a bot instance, e.g. !hello can be used to call the hello command
bot = commands.Bot(command_prefix='!', intents=intents) 

@bot.event
async def on_ready():
    logging.info(f"Bot: {bot.user} is ready to go.")
    logging.info(f"Bot is in {len(bot.guilds)} guilds")


async def load_cogs():
    for filename in os.listdir(pathlib.Path(__file__).parent / "cogs"):
        if filename.endswith(".py"):
            try:
                # remove the .py extension and use it as the cog name
                await bot.load_extension(f"cogs.{filename[:-3]}")
                logging.info(f"Loaded cog: {filename[:-3]}")
            except Exception as e:
                logging.error(f"Failed to load cog {filename[:-3]}: {e}")

async def main():
        await load_cogs()
        await bot.start(token)


# run the bot
if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logging.info("Bot shutting down...")

import discord
from discord.ext import commands
import logging
from dotenv import load_dotenv
import os
import time
import pathlib

baseDirectory = pathlib.Path(__file__).resolve().parent
logsFile = baseDirectory / 'logs' / 'discord-bot.log'
start = time.perf_counter()

load_dotenv() # loads environment variables from .env file
token = os.getenv('DISCORD_TOKEN')


handler = logging.FileHandler(filename=str(logsFile), encoding='utf-8', mode='w')

# intents (all permissions are here, manually enable in discord developer portal and in code)
intents = discord.Intents.default()
intents.message_content = True
intents.members = True

# create a bot instance
bot = commands.Bot(command_prefix='!', intents=intents) # e.g. !hello can be used to call the hello command

@bot.event
async def on_ready():
    print(f"Bot: {bot.user.name} is ready to go!")
    

@bot.event
async def on_message(message):
    if message.author == bot.user:
        return
    # put cool bot commands for messages here :>
    if "pickle" in message.content.lower():
        await message.channel.send("PICKLEBALLLLLLLL!")

    await bot.process_commands(message) # allows handling of all other messages sent in the server anywhere else (must include!)

# run the bot
bot.run(token, log_handler=handler, log_level=logging.DEBUG)

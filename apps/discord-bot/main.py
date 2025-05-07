import discord
from discord.ext import commands
import logging
from dotenv import load_dotenv
import os

load_dotenv() # loads environment variables from .env file
token = os.getenv('DISCORD_TOKEN')
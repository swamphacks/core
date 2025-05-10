import discord
from discord.ext import commands
import re

# regex for potential spam messages
SPAM_PATTERNS = [
    # selling tickets
    re.compile(r"""
        (?i)
        \bselling\b
        .{0,60}?
        (tickets?|passes?|spots?|seats?)
        (.{0,20}?\bfor\b.{0,20}?\$?\d{2,5})?
    """, re.VERBOSE),

    # subleasing
    re.compile(r"""
        (?i)
        \b(sublease|subleasing|sublet|lease)\b
        (.{0,20}?\bfor\b.{0,20}?\$?\d{2,5})?
    """, re.VERBOSE),

    # DM if interested
    re.compile(r"""
        (?i)
        dm\s+
        (me\s+)?(if\s+)?interested\b
    """, re.VERBOSE),
]


class AntiSpam(commands.Cog):
    def __init__(self, bot):
        self.bot = bot
        self.ignore_channels = set()

    @commands.Cog.listener()
    async def on_message(self, message):
        
        # ignore a bot message
        if message.author.bot:
            return
        
        # ignore messages in ignored channels
        if message.channel.id in self.ignore_channels:
            return

        # check for spam patterns
        content = message.content.lower()
        for pattern in SPAM_PATTERNS:
            if pattern.search(content):
                try:
                    original_message = message.content
                    await message.delete()
                    # send message in channel
                    await message.channel.send(f"{message.author.mention} Your message has been deleted for potential scam/spam.", delete_after=5)
                    
                    # send DM to user
                    try:
                        dm_message = (
                            f"Your message was automatically deleted for potential scam/spam content:\n"
                            f"```{original_message}```\n"
                            f"If you believe this was a mistake, please contact server staff for review."
                        )
                        await message.author.send(dm_message)
                    except discord.Forbidden:
                        print(f"Could not send DM to user {message.author}")
                    
                except discord.Forbidden:
                    print("bot lacks permissions to delete messages")
                    pass
                break

async def setup(bot):
    await bot.add_cog(AntiSpam(bot))

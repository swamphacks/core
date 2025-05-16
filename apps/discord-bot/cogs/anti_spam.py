import discord
from discord.ext import commands
import re
from typing import List, Pattern, Set
from utils.messaging import send_channel_message, send_dm



# Regular expressions for detecting spam patterns
SPAM_PATTERNS: List[Pattern] = [
    # Pattern for detecting ticket sales
    re.compile(r"""
        (?i)
        \bselling\b
        .{0,60}?
        (tickets?|passes?|spots?|seats?)
        (.{0,20}?\bfor\b.{0,20}?\$?\d{2,5})?
    """, re.VERBOSE),

    # Pattern for detecting subleasing advertisements
    re.compile(r"""
        (?i)
        \b(sublease|subleasing|sublet|lease)\b
        (.{0,20}?\bfor\b.{0,20}?\$?\d{2,5})?
    """, re.VERBOSE),

    # Pattern for detecting "DM if interested" messages
    re.compile(r"""
        (?i)
        dm\s+
        (me\s+)?(if\s+)?interested\b
    """, re.VERBOSE),
]


class AntiSpam(commands.Cog):
    """A cog that detects and handles spam messages in Discord channels.
    
    This cog monitors all messages in the server and automatically removes
    messages that match predefined spam patterns. It also notifies users
    when their messages are removed.
    """

    def __init__(self, bot: commands.Bot) -> None:
        """Initialize the AntiSpam cog
        
        Args:
            bot: Discord bot instance
        """
        self.bot: commands.Bot = bot
        self.ignore_channels: Set[int] = set()

    
    def is_spam(self, message: discord.Message) -> bool:
        """Check if a message matches any spam patterns
        
        Args:
            message: Discord message to check
            
        Returns:
            bool: True if the message is determined to be spam, False otherwise
        """
        content: str = message.content

        for pattern in SPAM_PATTERNS:
            if pattern.search(content):
                return True
            
        # TODO: Implement additional spam detection methods:
        # - Check for repeated messages
        # - Detect excessive mentions
        # - Check for suspicious links
        # - Monitor message frequency
        
        return False

    async def handle_spam(self, message: discord.Message) -> None:
        """Handle a detected spam message
        
        This method:
        1. Deletes the spam message
        2. Sends a notification in the channel by calling send_channel_message
        3. Sends a DM to the user explaining why their message was removed by calling send_dm
        
        Args:
            message: The spam message to handle
            
        Note:
            If the bot lacks permissions to delete messages or send DMs,
            the error will be logged but not raised.
        """
        try:
            await message.delete()
            content: str = f"{message.author.mention} Your message has been deleted for potential scam/spam."
            await send_channel_message(message.channel, content, delete_after=5)

            content = (
                f"Your message was automatically deleted for potential scam/spam content:\n"
                f"```{message.content}```\n"
                f"If you believe this was a mistake, please contact server staff for review."
            )
            await send_dm(message.author, content)
        except discord.Forbidden:
            print("Bot lacks permissions to delete messages")

    @commands.Cog.listener()
    async def on_message(self, message: discord.Message) -> None:
        """Event listener for new messages
        
        This method is called for every message sent in the server.
        1. Ignores messages from bots
        2. Ignores messages in ignored channels
        3. Calls is_spam function to determine if message is spam
        
        Args:
            message: The message that triggered the event
        """
        if message.author.bot:
            return
        
        if message.channel.id in self.ignore_channels:
            return

        if self.is_spam(message):
            await self.handle_spam(message)


async def setup(bot: commands.Bot) -> None:
    """Add the AntiSpam cog to the bot
    
    Args:
        bot: Discord bot instance
    """
    await bot.add_cog(AntiSpam(bot))

import discord
from discord.ext import commands
import re
from typing import List, Pattern, Set
from utils.messaging import send_channel_message, send_dm
from collections import defaultdict, deque
import time
from discord.utils import utcnow
import datetime
from collections import Counter


# Regular expressions for detecting spam patterns
SPAM_PATTERNS: List[Pattern] = [
    # # Pattern for detecting ticket sales
    # re.compile(r"""
    #     (?i)
    #     \bselling\b
    #     (.{0,20}?\bfor\b.{0,20}?\$?\d{2,5})?
    # """, re.VERBOSE),

    # # Pattern for detecting subleasing advertisements
    # re.compile(r"""
    #     (?i)
    #     \b(sublease|subleasing|sublet|lease)\b
    #     (.{0,20}?\bfor\b.{0,20}?\$?\d{2,5})?
    # """, re.VERBOSE),

    # # Pattern for detecting "DM if interested" messages
    # re.compile(r"""
    #     (?i)
    #     dm\s+
    #     (me\s+)?(if\s+)?interested\b
    # """, re.VERBOSE),
    
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
        self.user_message_history: dict[int, deque[tuple[str, float]]] = defaultdict(deque)
        self.repetition_window_sec = 30
        self.repetition_threshold = 5
        
    def is_repeated_message(self, message: discord.Message) -> tuple[bool, str | None]:
        """Check if a message is repeated
        
        Args:
            message: Discord message to check
        """
        user_id = message.author.id
        content = message.content.strip().lower()
        now = time.time()
        
        history = self.user_message_history[user_id]
        
        new_history = deque(maxlen=self.repetition_threshold)
        for msg, timestamp in history:
            if now - timestamp < self.repetition_window_sec:
                new_history.append((msg, timestamp))
            
        new_history.append((content, now))
        self.user_message_history[user_id] = new_history
        
        repeat_count = 0
        for msg, _ in new_history:
            if msg == content:
                repeat_count += 1
        
        if repeat_count >= self.repetition_threshold:
            return True, "Your message has been deleted for repeated content and spam."
        return False, None

    def is_excessive(self, message: discord.Message) -> tuple[bool, str | None]:
        """Check if a message has excessive spam or mentions
        
        Args:
            message: Discord message to check
        """
        words = message.content.split()
        words_count = Counter(words)
        for word, count in words_count.items():
            if word.startswith("<@") and count > 3:
                return True, f"Please refrain from excessively mentioning another user."
            if count > 10:
                return True, f"Your message has been deleted for excessive spamming {word}."
            
        return False, None
    
    def is_spam(self, message: discord.Message) -> tuple[bool, str | None]:
        """Check if a message matches any spam patterns
        
        Args:
            message: Discord message to check
            
        Returns:
            bool: True if the message is determined to be spam, False otherwise
        """

        # TODO: Implement additional spam detection methods:
        # use soft delete method (immediately halt the user and then asynchronously delete the messages)
        # - Check for repeated messages
        # - Detect excessive mentions
        # - Check for suspicious links
        # - Monitor message frequency
        
        # Check repeated messages
        repeated, reason = self.is_repeated_message(message)
        if repeated:
            return True, reason

        # Check excessive mentions
        mentions, reason = self.is_excessive(message)
        if mentions:
            return True, reason
        
        return False, None
    
    async def timeout_user(self, member: discord.Member, reason: str = "Spam", duration_seconds: int = 30) -> None:
        try:
            until = utcnow() + datetime.timedelta(seconds=duration_seconds)
            await member.timeout(until, reason=reason)
        except discord.Forbidden:
            print("Bot does not have permission to timeout members.")
        except discord.HTTPException as e:
            print(f"Failed to timeout member: {e}")
        except TypeError as e:
            print(f"Type error when applying timeout: {e}")

    async def handle_spam(self, message: discord.Message, reason: str | None = None) -> None:
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
            if isinstance(message.author, discord.Member):
                await self.timeout_user(message.author, reason)
            
            await message.delete()
            content: str = f"{message.author.mention} {reason}"
            await send_channel_message(message.channel, content, delete_after=15)
            
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
        if message.author.bot or message.channel.id in self.ignore_channels:
            return

        is_spam, reason = self.is_spam(message)
        if is_spam:
            await self.handle_spam(message, reason)


async def setup(bot: commands.Bot) -> None:
    """Add the AntiSpam cog to the bot
    
    Args:
        bot: Discord bot instance
    """
    await bot.add_cog(AntiSpam(bot))

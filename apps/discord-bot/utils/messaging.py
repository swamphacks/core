import discord
from typing import Union, Optional


async def send_channel_message(
    channel: discord.TextChannel,
    content: str,
    delete_after: Optional[int] = None
) -> None:
    """Send a message to a Discord channel with optional auto-deletion after some time
    
    Args:
        channel: Discord text channel to send the message to
        content: Message content to send
        delete_after: Number of seconds after which to delete the message, or None to keep it
        
    Note:
        If the bot lacks permissions to send messages, the error will be logged but not raised
    """
    try:
        await channel.send(content, delete_after=delete_after)
    except discord.Forbidden:
        print(f"[Error] Missing permissions to send message in #{channel.name}")
    except discord.HTTPException as e:
        print(f"[Error] Failed to send message in #{channel.name}: {e}")


async def send_dm(
    user: Union[discord.User, discord.Member],
    content: str
) -> None:
    """Send a direct message to a Discord user
    
    Args:
        user: Discord user or member to send the DM to
        content: The message content to send
        
    Note:
        If the bot cannot send DMs to the user, the error will be logged but not raised
    """
    try:
        await user.send(content)
    except discord.Forbidden:
        print(f"Cannot DM {user}")
    except discord.HTTPException as e:
        print(f"Failed to send DM: {e}")

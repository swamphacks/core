import discord

def get_next_thread_name(channel: discord.TextChannel) -> str:
    """
    Function to get the next thread name for a given channel and return the next thread name
    
    Args:
        channel: Discord text channel object
        
    Returns:
        str: The next thread name
    """
    # logic to find next support-<n> name
    existing_threads =  [t for t in channel.threads if t.name.startswith(f"{channel.name}-")]
    used_numbers = []
    
    for t in existing_threads:
        try:
            suffix = int(t.name.split("-")[1])
            used_numbers.append(suffix)
        except (IndexError, ValueError):
            continue
    
    next_number = max(used_numbers, default=0) + 1
    thread_name = f"{channel.name}-{next_number}"
    return thread_name
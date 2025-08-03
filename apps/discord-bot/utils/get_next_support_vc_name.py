import discord


def get_next_support_vc_name(category: discord.CategoryChannel) -> str :
    """    
    Generate the next available support voice channel name in a given category.

    Args:
        category (discord.CategoryChannel): The category where the voice channels are located.

    Returns:
        str: The next available voice channel name in the format "VC-<number>".
    """
    existing_channels = [channel.name for channel in category.voice_channels if channel.name.startswith("VC-")]
    
    used_numbers = []
    
    for vc in existing_channels:
        try:
            suffix = int(vc.split("-")[1])
            used_numbers.append(suffix)
        except (IndexError, ValueError):
            continue
    
    next_number = max(used_numbers, default=0) + 1
    vc_name = f"VC-{next_number}"
    return vc_name
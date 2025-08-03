import discord

def get_available_mentors(mentor_role: discord.Role) -> list[discord.Member]:
    """
    Get a list of available mentors from the mentor role.

    Args:
        mentor_role: The mentor role to check for available mentors.

    Returns:
        A list of discord.Member objects representing available mentors.
    """
    available_mentors = []
    for mentor in mentor_role.members:
        if "Available Mentor" in [role.name for role in mentor.roles]:
            available_mentors.append(mentor)
    return available_mentors
    
async def set_all_mentors_available(mentor_role: discord.Role) -> None:
    """
    Set all mentors in the mentor role to available.

    Args:
        mentor_role: The mentor role containing the mentors.
    
    Returns:
        None
    """
    available_role = discord.utils.get(mentor_role.guild.roles, name="Available Mentor")
    if not available_role:
        return
    for mentor in mentor_role.members:
        # print(f"Checking mentor: {mentor.name}")
        if "Available Mentor" not in [role.name for role in mentor.roles]:
            
            await mentor.add_roles(available_role)
        else:
            # print(f"{mentor.name} is already an available mentor")
            return


async def set_busy_mentor(mentor: discord.Member, busy: bool) -> None:
    """
    Set a mentor's status to busy or available.
    
    Args:
        mentor (discord.Member): The mentor to set the status for.
        busy (bool): True to set the mentor as busy, False to set them as available.
    Returns:
        None
    """
    busy_role = discord.utils.get(mentor.guild.roles, name="Busy Mentor")
    if not busy_role:
        return
    if busy == True:
        await mentor.add_roles(busy_role)
    else:
        if busy_role in mentor.roles:
            await mentor.remove_roles(busy_role)
        else:
            print(f"No **Busy Mentor** role to remove from {mentor.name}")
            return
        

async def set_available_mentor(mentor: discord.Member, available: bool) -> None:
    """
    Set a mentor's status to available or not available.
    
    Args:
        mentor (discord.Member): The mentor to set the status for.
        available (bool): True to set the mentor as available, False to set them as not available.
    
    Returns:
        None"""
    available_role = discord.utils.get(mentor.guild.roles, name="Available Mentor")
    if not available_role:
        print(f"No **Available Mentor** role found for {mentor.name}")
        return
    if available == True:
        await mentor.add_roles(available_role)
    else:
        
        if available_role in mentor.roles:
            await mentor.remove_roles(available_role)
        else:
            print(f"No **Available Mentor** role to remove from {mentor.name}")
            return
        
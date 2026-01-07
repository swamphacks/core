import discord
from utils.roles_config import get_acceptable_roles, RoleNames, get_role_id

def get_available_mentors(guild: discord.Guild) -> list[discord.Member]:
    """
    Get a list of available mentors from all acceptable roles.

    Args:
        guild: The Discord guild to check for available mentors.

    Returns:
        A list of discord.Member objects representing available mentors.
    """
    acceptable_roles = get_acceptable_roles()
    available_mentors = []
    available_mentor_role_name = RoleNames.AVAILABLE_MENTOR
    
    for role_name in acceptable_roles:
        # Try to get role by ID first, then by name
        role_id = get_role_id(role_name)
        if role_id:
            role = guild.get_role(int(role_id))
        else:
            role = discord.utils.get(guild.roles, name=role_name)
        
        if role:
            for mentor in role.members:
                if available_mentor_role_name in [role.name for role in mentor.roles]:
                    # Avoid duplicates if user has multiple acceptable roles
                    if mentor not in available_mentors:
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
    available_role_name = RoleNames.AVAILABLE_MENTOR
    available_role_id = get_role_id(available_role_name)
    
    if available_role_id:
        available_role = mentor_role.guild.get_role(int(available_role_id))
    else:
        available_role = discord.utils.get(mentor_role.guild.roles, name=available_role_name)
    
    if not available_role:
        return
    
    for mentor in mentor_role.members:
        # print(f"Checking mentor: {mentor.name}")
        if available_role_name not in [role.name for role in mentor.roles]:
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
    busy_role_name = RoleNames.BUSY_MENTOR
    busy_role_id = get_role_id(busy_role_name)
    
    if busy_role_id:
        busy_role = mentor.guild.get_role(int(busy_role_id))
    else:
        busy_role = discord.utils.get(mentor.guild.roles, name=busy_role_name)
    
    if not busy_role:
        return
    
    if busy == True:
        await mentor.add_roles(busy_role)
    else:
        if busy_role in mentor.roles:
            await mentor.remove_roles(busy_role)
        else:
            print(f"No **{busy_role_name}** role to remove from {mentor.name}")
            return
        

async def set_available_mentor(mentor: discord.Member, available: bool) -> None:
    """
    Set a mentor's status to available or not available.
    
    Args:
        mentor (discord.Member): The mentor to set the status for.
        available (bool): True to set the mentor as available, False to set them as not available.
    
    Returns:
        None"""
    available_role_name = RoleNames.AVAILABLE_MENTOR
    available_role_id = get_role_id(available_role_name)
    
    if available_role_id:
        available_role = mentor.guild.get_role(int(available_role_id))
    else:
        available_role = discord.utils.get(mentor.guild.roles, name=available_role_name)
    
    if not available_role:
        print(f"No **{available_role_name}** role found for {mentor.name}")
        return
    
    if available == True:
        await mentor.add_roles(available_role)
    else:
        if available_role in mentor.roles:
            await mentor.remove_roles(available_role)
        else:
            print(f"No **{available_role_name}** role to remove from {mentor.name}")
            return
        
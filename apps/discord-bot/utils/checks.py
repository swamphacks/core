import os
from discord import app_commands, Interaction, Permissions
import json
from typing import Callable, Coroutine, Any
from utils.roles_config import get_acceptable_roles

def has_bot_full_access() -> Callable[[Interaction], Coroutine[Any, Any, bool]]:
    """
    Check if the user has one of the acceptable roles or is an administrator for slash commands.

    Returns:
        bool: True if the user has one of the acceptable roles or is an administrator, False otherwise
    """
    async def predicate(interaction: Interaction):
        # Get acceptable roles from config
        acceptable_roles = get_acceptable_roles()
        
        # Ensure interaction is in a channel and a user exists
        if not interaction.guild or not interaction.user:
            return False
        
        member = interaction.guild.get_member(interaction.user.id)
        if not member:
            return False
        
        # get all the permissions for the member
        perms: Permissions = member.guild_permissions
        
        # check if user has moderator privileges
        # for now it just checks if they are an admin but adjust this later
        if any(role.name in acceptable_roles for role in member.roles):
            # print(f"User {member.name} has moderator privileges.")
            return True
        else:
            # print(f"User {member.name} does not have moderator privileges.")
            return False
       
    return app_commands.check(predicate)

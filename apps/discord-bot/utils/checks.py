import os
from discord import app_commands, Interaction, Permissions
import discord
import json
from typing import Callable, Coroutine, Any
from utils.roles_config import get_acceptable_roles, RoleNames

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

def is_mod_slash() -> Callable[[Interaction], Coroutine[Any, Any, bool]]:
    """
    Check if the user has moderator role for slash commands.
    Uses the same logic as has_bot_full_access() for consistency.

    Returns:
        bool: True if the user has moderator privileges, False otherwise
    """
    # Reuse the same logic as has_bot_full_access
    return has_bot_full_access()

def requires_admin() -> Callable[[Interaction], Coroutine[Any, Any, bool]]:
    """
    Check if the user has the Admin role for slash commands.
    This is a stricter check than has_bot_full_access() and requires the Admin role specifically.

    Returns:
        bool: True if the user has the Admin role, False otherwise
    """
    async def predicate(interaction: Interaction):
        # Ensure interaction is in a guild and a user exists
        if not interaction.guild or not interaction.user:
            return False
        
        member = interaction.guild.get_member(interaction.user.id)
        if not member:
            return False
        
        # Check if user has Admin role
        admin_role = discord.utils.get(member.roles, name=RoleNames.ADMIN)
        if admin_role:
            return True
        
        return False
       
    return app_commands.check(predicate)

def has_bot_full_access_or_hacker() -> Callable[[Interaction], Coroutine[Any, Any, bool]]:
    """
    Check if the user has bot full access (moderator, mentor, bot, staff, admin) OR Hacker (XI) role.
    This allows Hacker (XI) to use specific commands like create_vc, add_to_thread, grant_vc_access
    while still restricting them from other moderator-only commands.

    Returns:
        bool: True if the user has bot full access or Hacker (XI) role, False otherwise
    """
    async def predicate(interaction: Interaction):
        # Ensure interaction is in a guild and a user exists
        if not interaction.guild or not interaction.user:
            return False
        
        member = interaction.guild.get_member(interaction.user.id)
        if not member:
            return False
        
        # First check if user has bot full access (moderator, mentor, bot, staff, admin)
        acceptable_roles = get_acceptable_roles()
        if any(role.name in acceptable_roles for role in member.roles):
            return True
        
        # Then check if user has Hacker (XI) role
        hacker_role = discord.utils.get(member.roles, name=RoleNames.HACKER_XI)
        if hacker_role:
            return True
        
        return False
       
    return app_commands.check(predicate)

def requires_admin_or_moderator() -> Callable[[Interaction], Coroutine[Any, Any, bool]]:
    """
    Check if the user has Admin or Moderator role for dangerous commands.
    This is a stricter check than has_bot_full_access() and only allows Admin and Moderator roles.

    Returns:
        bool: True if the user has Admin or Moderator role, False otherwise
    """
    async def predicate(interaction: Interaction):
        # Ensure interaction is in a guild and a user exists
        if not interaction.guild or not interaction.user:
            return False
        
        member = interaction.guild.get_member(interaction.user.id)
        if not member:
            return False
        
        # Check if user has Admin role
        admin_role = discord.utils.get(member.roles, name=RoleNames.ADMIN)
        if admin_role:
            return True
        
        # Check if user has Moderator role
        moderator_role = discord.utils.get(member.roles, name=RoleNames.MODERATOR)
        if moderator_role:
            return True
        
        return False
       
    return app_commands.check(predicate)

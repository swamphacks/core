"""
Role configuration for the Discord bot.

This module centralizes all role names and IDs used throughout the bot.
Role IDs are optional - if not provided, the bot will search by role name.

To configure role IDs:
1. Get the role ID from Discord (right-click role > Copy ID, with Developer Mode enabled)
2. Update the ROLE_IDS dictionary with the role ID as a string
3. Example: ROLE_IDS["Moderator"] = "123456789012345678"

To add/remove acceptable roles:
1. Modify the ACCEPTABLE_ROLES list
2. Add corresponding entries to ROLE_IDS if you want to use IDs
3. Update RoleNames class if adding new role types
"""

from typing import Optional

# Roles that are allowed to use all bot commands
ACCEPTABLE_ROLES: list[str] = ["Moderator", "Mentor (XI)", "Bot", "Staff (XI), Admin,"]

# Roles that can be set as available mentors (can be different from ACCEPTABLE_ROLES)
ACCEPTABLE_MENTOR_ROLES: list[str] = ["Mentor (XI)"]

# Optional: Map role names to role IDs for faster lookups
# If a role ID is None fallback to search by name
# Format: "Role Name": Optional[role_id_as_string]
ROLE_IDS: dict[str, Optional[str]] = {
    "Moderator": None,
    "Mentor": None,
    "Mentor (XI)": None,
    "Available Mentor": None,
    "Busy Mentor": None,
}

# Role names used throughout the bot
class RoleNames:
    """Centralized role name constants."""
    MODERATOR = "Moderator"
    MENTOR = "Mentor"
    MENTOR_XI = "Mentor (XI)"
    ADMIN = "Admin"
    AVAILABLE_MENTOR = "Available Mentor"
    BUSY_MENTOR = "Busy Mentor"


def get_role_id(role_name: str) -> Optional[str]:
    """
    Get the role ID for a given role name.
    
    Args:
        role_name: The name of the role
        
    Returns:
        The role ID as a string, or None if not configured
    """
    return ROLE_IDS.get(role_name)


def set_role_id(role_name: str, role_id: Optional[str]) -> None:
    """
    Set the role ID for a given role name.
    
    Args:
        role_name: The name of the role
        role_id: The role ID as a string, or None to search by name
    """
    ROLE_IDS[role_name] = role_id


def get_acceptable_roles() -> list[str]:
    """
    Get the list of acceptable roles for bot commands.
    
    Returns:
        List of role names
    """
    return ACCEPTABLE_ROLES.copy()


def set_acceptable_roles(roles: list[str]) -> None:
    """
    Set the list of acceptable roles for bot commands.
    
    Args:
        roles: List of role names
    """
    global ACCEPTABLE_ROLES
    ACCEPTABLE_ROLES = roles.copy()


def get_acceptable_mentor_roles() -> list[str]:
    """
    Get the list of roles that can be set as available mentors.
    
    Returns:
        List of role names
    """
    return ACCEPTABLE_MENTOR_ROLES.copy()


def set_acceptable_mentor_roles(roles: list[str]) -> None:
    """
    Set the list of roles that can be set as available mentors.
    
    Args:
        roles: List of role names
    """
    global ACCEPTABLE_MENTOR_ROLES
    ACCEPTABLE_MENTOR_ROLES = roles.copy()


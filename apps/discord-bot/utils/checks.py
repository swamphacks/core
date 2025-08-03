import os
from discord import app_commands, Interaction, Permissions
import json
from typing import Callable, Coroutine, Any

# def load_config():
#     config_path = os.path.join(os.path.dirname(__file__), "..", "config.json")
#     try:
#         with open(config_path) as f:
#             return json.load(f)
#     except FileNotFoundError:
#         raise FileNotFoundError(f"Config file not found at {config_path}")
#     except json.JSONDecodeError as e:
#         raise ValueError(f"Invalid JSON in config file: {e}")

# config = load_config()

def is_mod_slash() -> Callable[[Interaction], Coroutine[Any, Any, bool]]:
    """
    Check if the user has the moderator role for slash commands. The role is set in the config.json file.

    Returns:
        bool: True if the user has the moderator role, False otherwise
    """
    async def predicate(interaction: Interaction):
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
        if perms.administrator:
            return True
        else:
            return False
       
    return app_commands.check(predicate)

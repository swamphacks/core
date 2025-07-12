import os
from discord import app_commands, Interaction
import json
from typing import Callable, Coroutine, Any

def load_config():
    config_path = os.path.join(os.path.dirname(__file__), "..", "config.json")
    try:
        with open(config_path) as f:
            return json.load(f)
    except FileNotFoundError:
        raise FileNotFoundError(f"Config file not found at {config_path}")
    except json.JSONDecodeError as e:
        raise ValueError(f"Invalid JSON in config file: {e}")

config = load_config()

def is_mod_slash() -> Callable[[Interaction], Coroutine[Any, Any, bool]]:
    """
    Check if the user has the moderator role for slash commands. The role is set in the config.json file.

    Returns:
        bool: True if the user has the moderator role, False otherwise
    """
    async def predicate(interaction: Interaction):
        try:
            mod_role_id = int(config["roles"]["moderator"])
        except KeyError:
            raise ValueError("Moderator role ID not found in config")
        except ValueError:
            raise ValueError("Invalid moderator role ID in config")

        role_ids = [role.id for role in interaction.user.roles]

        has_permission = mod_role_id in role_ids
        return has_permission

    return app_commands.check(predicate)

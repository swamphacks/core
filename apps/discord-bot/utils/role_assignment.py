import logging
import uuid
import aiohttp
from typing import List, Tuple, Optional
from utils.db import Database

GET_ATTENDEES_WITH_DISCORD = """
SELECT
    a.account_id as discord_id,
    u.id as user_id,
    u.name,
    u.email
FROM auth.users u
JOIN event_roles er ON u.id = er.user_id
JOIN auth.accounts a ON u.id = a.user_id
WHERE er.event_id = %s
    AND er.role = 'attendee'
    AND a.provider_id = 'discord'
"""

logger = logging.getLogger(__name__)

async def get_attendees_for_event(event_id: str) -> List[Tuple[str, str, str, Optional[str]]]:
    """Query database for all attendees with Discord accounts for given event"""
    
    try:
        event_uuid = uuid.UUID(event_id)
    except ValueError:
        raise ValueError(f"Invalid event ID format: {event_id}")
        
    conn = Database.get_connection()
    try:
        cursor = conn.cursor()
        cursor.execute(GET_ATTENDEES_WITH_DISCORD, (str(event_uuid),))
        attendees = cursor.fetchall()
        cursor.close()
        return attendees
    finally:
        Database.return_connection(conn)
        
async def assign_role_via_webhook(
    webhook_url: str, 
    discord_id: str, 
    role_name: str, 
    guild_id: str,
    session: Optional[aiohttp.ClientSession] = None
) -> Tuple[str, Optional[str]]:
    """Assign a role to a user via webhook
    
    Args:
        webhook_url: URL of webhook to send requests to
        discord_id: Discord ID of user to assign role to
        role_name: Name of role to assign
        guild_id: Discord guild ID
        session: Optional aiohttp session to reuse (creates new one if not provided)
        
    Returns:
        Tuple of (status: str, error_message: Optional[str])
        status can be: "newly_assigned", "already_had", or "failed"
    """
    
    try:
        payload = {
            "action": "assign_role",
            "user_id": str(discord_id),
            "role_name": role_name,
            "guild_id": str(guild_id)
        }
        
        # Use provided session or create a new one
        if session:
            async with session.post(webhook_url, json=payload) as response:
                if response.status == 200:
                    text = await response.text()
                    if "already has the role" in text.lower():
                        return ("already_had", None)
                    elif "assigned successfully" in text.lower():
                        return ("newly_assigned", None)
                    else:
                        return ("unknown", text)
                else:
                    text = await response.text()
                    return ("failed", text)
        else:
            # Create a new session if none provided (for backwards compatibility)
            async with aiohttp.ClientSession() as new_session:
                async with new_session.post(webhook_url, json=payload) as response:
                    if response.status == 200:
                        text = await response.text()
                        if "already has the role" in text.lower():
                            return ("already_had", None)
                        elif "assigned successfully" in text.lower():
                            return ("newly_assigned", None)
                        else:
                            return ("unknown", text)
                    else:
                        text = await response.text()
                        return ("failed", text)
    except Exception as e:
        logger.error(f"Error assigning role via webhook: {e}")
        return ("failed", str(e))


async def assign_roles_to_attendees(webhook_url: str, attendees: List[Tuple[str, str, str, Optional[str]]], role_name: str, guild_id: str) -> Tuple[int, int, int, List[str]]:
    """Assign roles to attendees via webhook
    
    Args:
        webhook_url: URL of webhook to send requests to
        attendees: List of attendees to assign roles to
        role_name: Name of role to assign
        guild_id: Discord guild ID
        
    Returns:
        Tuple of (newly_assigned: int, already_had: int, failed: int, errors: List[str])
    """
    
    newly_assigned = 0
    already_had = 0
    failed = 0
    errors = []
    
    # Create a single session to reuse for all requests (better performance)
    async with aiohttp.ClientSession() as session:
        for discord_id, user_id, name, email in attendees:
            status, error_msg = await assign_role_via_webhook(webhook_url, discord_id, role_name, guild_id, session)  # Pass session
            
            if status == "newly_assigned":
                newly_assigned += 1
            elif status == "already_had":
                already_had += 1
            elif status == "failed":
                failed += 1
                error_detail = f"User {name} (Discord ID: {discord_id}, User ID: {user_id})"
                if error_msg:
                    error_detail += f": {error_msg}"
                errors.append(error_detail.strip())
            else:  # Handle "unknown" status
                failed += 1
                error_detail = f"User {name} (Discord ID: {discord_id}, User ID: {user_id}): Unknown status from webhook"
                if error_msg:
                    error_detail += f" - {error_msg}"
                errors.append(error_detail.strip())
    
    return (newly_assigned, already_had, failed, errors)

def format_assignment_summary(
    total_attendees: int,
    newly_assigned: int,
    already_had: int,
    failed: int,
    errors: List[str],
    max_errors_displayed: int = 20
) -> str:
    """Format assignment summary for logging
    
    Args:
        total_attendees: Total number of attendees queried
        newly_assigned: Number of roles newly assigned
        already_had: Number of users who already had the role
        failed: Number of failed assignments
        errors: List of error messages
        max_errors_displayed: Maximum number of errors to show
    """
    message = f"**Role Assignment Complete**\n\n"
    message += f"**Summary:**\n"
    message += f"- Total attendees in database: {total_attendees}\n"
    message += f"- Newly assigned: {newly_assigned}\n"
    message += f"- Already had role: {already_had}\n"
    message += f"- Failed assignments: {failed}\n"
    
    synced = newly_assigned + already_had
    sync_percentage = (synced / total_attendees * 100) if total_attendees > 0 else 0
    message += f"\n**Sync Status:**\n"
    message += f"- Up to date: {synced}/{total_attendees} ({sync_percentage:.1f}%)\n"
    
    if errors:
        message += f"\n**Errors ({len(errors)}):**\n"
        for error in errors[:max_errors_displayed]:
            message += f"- {error}\n"
        if len(errors) > max_errors_displayed:
            message += f"- ... and {len(errors) - max_errors_displayed} more errors\n"
    
    return message
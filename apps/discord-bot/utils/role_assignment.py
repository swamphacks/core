import logging
import asyncio
import aiohttp
from typing import List, Tuple, Optional

logger = logging.getLogger(__name__)

async def get_attendees_for_event(api_url: str, session_cookie: str, event_id: str) -> List[Tuple[str, str, str, Optional[str]]]:
    """Get attendees with Discord IDs for an event from API
    
    Args:
        api_url: Base URL of the API
        session_cookie: Session cookie for authentication
        event_id: Event ID (UUID)
        
    Returns:
        List of tuples (discord_id, user_id, name, email)
    """
    
    try:
        async with aiohttp.ClientSession() as session:
            headers = {"Cookie": f"sh_session_id={session_cookie}"}
            async with session.get(
                f"{api_url}/discord/event/{event_id}/attendees",
                headers=headers
            ) as response:
                if response.status == 200:
                    data = await response.json()
                    attendees = []
                    for attendee in data:
                        attendees.append((
                            attendee.get("discord_id", ""),
                            attendee.get("user_id", ""),
                            attendee.get("name", ""),
                            attendee.get("email")
                        ))
                    return attendees
                elif response.status == 404:
                    return []
                else:
                    text = await response.text()
                    logger.error(f"API error: {response.status} - {text}")
                    raise Exception(f"API error: {response.status}")
    except Exception as e:
        logger.error(f"Error fetching attendees from API: {e}")
        raise
        
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


async def assign_roles_to_attendees(webhook_url: str, attendees: List[Tuple[str, str, str, Optional[str]]], role_name: str, guild_id: str, chunk_size: int = 20, progress_callback=None, test_mode: bool = False) -> Tuple[int, int, int, List[str]]:
    """Assign roles to attendees via webhook in chunks
    
    Args:
        webhook_url: URL of webhook to send requests to
        attendees: List of attendees to assign roles to
        role_name: Name of role to assign
        guild_id: Discord guild ID
        chunk_size: Number of users to process at a time (default: 20)
        progress_callback: Optional async function to call with progress updates (current, total)
        test_mode: If True, stop on first error (default: False)
        
    Returns:
        Tuple of (newly_assigned: int, already_had: int, failed: int, errors: List[str])
    """
    newly_assigned = 0
    already_had = 0
    failed = 0
    errors = []
    
    total_attendees = len(attendees)
    
    # Create a single session to reuse for all requests (better performance)
    async with aiohttp.ClientSession() as session:
        # Process attendees in chunks
        for chunk_start in range(0, total_attendees, chunk_size):
            chunk_end = min(chunk_start + chunk_size, total_attendees)
            chunk = attendees[chunk_start:chunk_end]
            
            # Process this chunk
            for discord_id, user_id, name, email in chunk:
                status, error_msg = await assign_role_via_webhook(webhook_url, discord_id, role_name, guild_id, session)
                
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
                    
                    # In test mode, stop on first error
                    if test_mode:
                        logger.warning(f"Test mode: Stopping on first error - {error_detail}")
                        return (newly_assigned, already_had, failed, errors)
                else:  # Handle "unknown" status
                    failed += 1
                    error_detail = f"User {name} (Discord ID: {discord_id}, User ID: {user_id}): Unknown status from webhook"
                    if error_msg:
                        error_detail += f" - {error_msg}"
                    errors.append(error_detail.strip())
                    
                    # In test mode, stop on first error
                    if test_mode:
                        logger.warning(f"Test mode: Stopping on first error - {error_detail}")
                        return (newly_assigned, already_had, failed, errors)
            
            # Call progress callback if provided
            if progress_callback:
                await progress_callback(chunk_end, total_attendees)
            
            # Small delay between chunks to avoid rate limits (except after last chunk)
            if chunk_end < total_attendees:
                await asyncio.sleep(0.5)  # 500ms delay between chunks
    
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
    DISCORD_MAX_LENGTH = 2000
    
    message = f"**Role Assignment Complete**\n\n"
    message += f"**Summary:**\n"
    message += f"- Total attendees: {total_attendees}\n"
    message += f"- Newly assigned: {newly_assigned}\n"
    message += f"- Already had role: {already_had}\n"
    message += f"- Failed assignments: {failed}\n"
    
    synced = newly_assigned + already_had
    sync_percentage = (synced / total_attendees * 100) if total_attendees > 0 else 0
    message += f"\n**Sync Status:**\n"
    message += f"- Up to date: {synced}/{total_attendees} ({sync_percentage:.1f}%)\n"
    
    if errors:
        message += f"\n**Errors ({len(errors)}):**\n"
        
        # Calculate how many errors we can fit
        base_length = len(message)
        remaining_chars = DISCORD_MAX_LENGTH - base_length - 50  # Reserve 50 chars for "... and X more"
        
        errors_to_show = []
        current_length = 0
        
        for i, error in enumerate(errors[:max_errors_displayed]):
            error_line = f"- {error}\n"
            if current_length + len(error_line) > remaining_chars:
                break
            errors_to_show.append(error)
            current_length += len(error_line)
        
        for error in errors_to_show:
            # Truncate individual error messages if they're too long
            if len(error) > 150:
                error = error[:147] + "..."
            message += f"- {error}\n"
        
        if len(errors) > len(errors_to_show):
            remaining = len(errors) - len(errors_to_show)
            message += f"- ... and {remaining} more errors\n"
    
    # Final safety check - truncate if still too long
    if len(message) > DISCORD_MAX_LENGTH:
        message = message[:DISCORD_MAX_LENGTH - 3] + "..."
    
    return message
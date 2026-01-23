from discord.ext import commands
from discord import app_commands
import discord
import aiohttp
import logging
from typing import Literal, Optional
from utils.checks import is_mod_slash, has_bot_full_access, requires_admin, has_bot_full_access_or_hacker, requires_admin_or_moderator
import re
from typing import Literal
from utils.mentor_functions import set_all_mentors_available
from utils.role_assignment import (get_attendees_for_event, format_assignment_summary, assign_roles_to_attendees)
import os
from utils.get_next_support_vc_name import get_next_support_vc_name
from chatbot.llm import llm_response

class General(commands.Cog):
    """A cog containing general utility commands for the server
    
    This cog includes commands for:
    - Basic server interactions
    - Role management
    - Fun commands
    """
    def __init__(self, bot: commands.Bot) -> None:
        """Initialize the General cog
        
        Args:
            bot: Discord bot instance
        """
        self.bot: commands.Bot = bot
    
    def get_role(self, guild: discord.Guild, role_name: str) -> Optional[discord.Role]:
        """Helper to get a role by name from a guild."""
        return discord.utils.get(guild.roles, name=role_name)
    
    @commands.command()
    async def test(self, ctx: commands.Context) -> None:
        """Send a test message
        
        Args:
            ctx: The command context
        """
        await ctx.send("Testing")

    @app_commands.command(name="delete", description="Delete X amount of messages based on the number you provide")
    @app_commands.describe(
        amount="The amount of messages to delete"
    )
    @requires_admin_or_moderator()
    async def delete(
        self,
        interaction: discord.Interaction,
        amount: int
    ) -> None:
        """Delete X amount of messages based on the number you provide
        
        Args:
            interaction: The interaction that triggered this command
            amount: The amount of messages to delete
        """
        await interaction.response.defer(ephemeral=True)
        deleted = await interaction.channel.purge(limit=amount)
        await interaction.followup.send(
            f"Deleted {len(deleted)} messages.",
            ephemeral=True
        )
    
    @app_commands.command(name="delete_all_threads", description="Delete all threads in a specified channel")
    @requires_admin_or_moderator()
    async def delete_all_threads(self, interaction: discord.Interaction, channel: discord.TextChannel, delete_archived: bool = False) -> None:
        """Delete all threads in a specified channel
        
        Args:
            interaction: The interaction that triggered this command
        """
        guild = interaction.guild
        if channel not in guild.text_channels:
            await interaction.response.send_message("Error: Could not find the specified channel.", ephemeral=True)
            return
        
        for thread in channel.threads:
            # this only iterates over active threads so archived threads will not be deleted
            await thread.delete()
        
        if delete_archived:
            # delete archived public or private threads
            async for thread in channel.archived_threads(private=False):
                await thread.delete()
            async for thread in channel.archived_threads(private=True):
                await thread.delete()

        await interaction.response.send_message(
            f"All threads in {channel.mention} {'including archived ones ' if delete_archived else ''}have been deleted.",
            ephemeral=True
        )
        
    @app_commands.command(name="delete_all_vcs", description="Delete all voice channels in a specified category")
    @requires_admin_or_moderator()
    async def delete_all_vcs(self, interaction: discord.Interaction, category: discord.CategoryChannel) -> None:
        """Delete all voice channels in a specified category
        
        Args:
            interaction: The interaction that triggered this command
            category: The category from which to delete all voice channels
        """
        guild = interaction.guild
        if category not in guild.categories:
            await interaction.response.send_message("Error: Could not find the specified category.", ephemeral=True)
            return
        
        for channel in category.voice_channels:
            await channel.delete()
        
        await interaction.response.send_message("All voice channels in the specified category have been deleted.", ephemeral=True)
    
    @app_commands.command(
        name="role",
        description="Assign or remove a role from yourself"
    )
    @app_commands.describe(
        action="Whether to assign or remove the role",
        role="The role to assign or remove"
    )
    @has_bot_full_access()
    async def manage_role(
        self,
        interaction: discord.Interaction,
        action: Literal["assign", "remove"],
        role: discord.Role,
        member: discord.Member
    ) -> None:
        """Manage roles for the user who triggered the command
        
        Args:
            interaction: The interaction that triggered this command
            action: Whether to assign or remove the role
            role: The role to assign or remove
            
        Note:
            This command will:
            1. Check if the user already has/doesn't have the role
            2. Assign or remove the role if conditions are met
            3. Send appropriate feedback messages
        """
        # fetch the member to give the role to
        member = await interaction.guild.fetch_member(member.id)

        has_role = role in member.roles

        if action == "assign":
            if has_role:
                await interaction.response.send_message(
                    f"{member.mention} already has the **{role.name}** role.",
                    ephemeral=True
                )
                return
                
            try:
                await member.add_roles(role)
                await interaction.response.send_message(
                    f"Assigned **{role.name}** to {member.mention}.",
                    ephemeral=True
                )
                await interaction.followup.send(
                    f"{interaction.user.mention} assigned **{role.name}** role to {member.mention}.",
                    delete_after=5
                )
            except discord.Forbidden:
                await interaction.response.send_message(
                    "I don't have permission to assign roles! Please adjust my permissions.",
                    ephemeral=True
                )
        elif action == "remove":
            if not has_role:
                await interaction.response.send_message(
                    f"{member.mention} does not have the **{role.name}** role.",
                    ephemeral=True
                )
                return
                
            try:
                await member.remove_roles(role)
                await interaction.response.send_message(
                    f"Removed **{role.name}** from {member.mention}.",
                    ephemeral=True
                )
                await interaction.followup.send(
                    f"{interaction.user.mention} removed **{role.name}** role from {member.mention}.",
                    delete_after=5
                )
            except discord.Forbidden:
                await interaction.response.send_message(
                    "I don't have permission to assign roles! Please adjust my permissions.",
                    ephemeral=True
                )
        else:
            # This will not be reached but just wanted to show add and remove for commands
            await interaction.response.send_message(
                "Invalid action. Please use 'assign' or 'remove'.",
                ephemeral=True
            )

    @app_commands.command(name="set_available_mentors", description="Set available mentors in the server")
    @requires_admin_or_moderator()
    async def set_all_mentors_available(self, interaction: discord.Interaction) -> None:
        """
        Set all users in the server with acceptable mentor roles to "Available Mentor"
        
        Args:
            interaction: The interaction that triggered this command
        """
        from utils.roles_config import get_acceptable_mentor_roles, RoleNames, get_role_id
        
        await interaction.response.defer(ephemeral=True)
        
        acceptable_mentor_roles = get_acceptable_mentor_roles()
        available_role_name = RoleNames.AVAILABLE_MENTOR
        available_role_id = get_role_id(available_role_name)
        
        if available_role_id:
            available_role = interaction.guild.get_role(int(available_role_id))
        else:
            available_role = discord.utils.get(interaction.guild.roles, name=available_role_name)
        
        if not available_role:
            await interaction.followup.send(f"Error: Could not find the **{available_role_name}** role. Please create it before using this command.", ephemeral=True)
            return
        
        total_updated = 0
        roles_found = []
        
        for role_name in acceptable_mentor_roles:
            role_id = get_role_id(role_name)
            if role_id:
                role = interaction.guild.get_role(int(role_id))
            else:
                role = self.get_role(interaction.guild, role_name)
            
            if role:
                roles_found.append(role_name)
                for member in role.members:
                    if available_role_name not in [r.name for r in member.roles]:
                        await member.add_roles(available_role)
                        total_updated += 1
        
        if not roles_found:
            await interaction.followup.send(f"Error: Could not find any of the acceptable mentor roles ({', '.join(acceptable_mentor_roles)}). Please create them before using this command.", ephemeral=True)
            return
        
        await interaction.followup.send(f"Successfully set {total_updated} members from roles {', '.join(roles_found)} to available.", ephemeral=True)
        
        
    @app_commands.command(name="add_to_thread", description="Add a user to the support thread")
    @app_commands.describe(user="The user to add to the thread")
    @has_bot_full_access()
    async def add_to_thread(self, interaction: discord.Interaction, user: discord.Member) -> None:
        """
        Add a specified user to a support thread
        
        Args:
            interaction: The interaction that triggered this command
            user: The user to add to the thread
        """
        
        # first ensure command is being executed in a thread
        if not isinstance(interaction.channel, (discord.Thread,)):
            await interaction.response.send_message(
                "This command can only be used in a support thread.", ephemeral=True
            )
            return
        # next ensure the thread is in the support channel specifically (check if the thread's parent exists as well)
        if not interaction.channel.parent or interaction.channel.parent.name != "support":
            if not interaction.channel.category or interaction.channel.category.name != "--- SwampHacks XI (Support-VCs) ---":
                await interaction.response.send_message('This command can only be used in the "--- SwampHacks XI (Support-VCs) ---" category.', ephemeral=True)
            return
        
        # check if user is already in the thread
        try:
            await interaction.channel.fetch_member(user.id)
            await interaction.response.send_message(f"{user.mention} is already in this thread.", ephemeral=True)
            return
        except discord.NotFound:
            pass
        
        try:
            await interaction.channel.add_user(user)
            await interaction.response.send_message(f"{user.mention} has been added to the thread.", ephemeral=True)
        except discord.Forbidden:
            await interaction.response.send_message("I don't have permission to add users to this thread.", ephemeral=True)
        except Exception as e:
            await interaction.response.send_message(f"An error occurred: {str(e)}", ephemeral=True)
            
    @app_commands.command(name="create_vc", description="Creates a voice channel for support")
    @has_bot_full_access_or_hacker()
    async def create_vc(self, interaction: discord.Interaction) -> None:
        """
        Create a voice channel for support, requires a --- SwampHacks XI (Support-VCs) --- category and a Mentor role
        Can only be used in a thread, and only one VC per thread is allowed.

        Args:
            interaction: The interaction that triggered this command
        """
        from utils.roles_config import RoleNames, get_role_id
        from components.ticket_state import thread_vc_mapping
        
        # Check if command is being used in a thread
        if not isinstance(interaction.channel, discord.Thread):
            await interaction.response.send_message(
                "❌ This command can only be used inside a support thread.",
                ephemeral=True
            )
            return
        
        thread = interaction.channel
        
        # Check if thread already has a VC associated
        if thread.id in thread_vc_mapping:
            existing_vc_id = thread_vc_mapping[thread.id]
            existing_vc = interaction.guild.get_channel(existing_vc_id)
            if existing_vc:
                await interaction.response.send_message(
                    f"❌ This thread already has an associated voice channel: {existing_vc.mention}",
                    ephemeral=True
                )
                return
            else:
                # VC was deleted, remove from mapping
                del thread_vc_mapping[thread.id]
        
        # Extract thread number from thread name (e.g., "support-5" -> "5")
        thread_number = None
        try:
            if "-" in thread.name:
                thread_number = thread.name.split("-")[-1]
                # Verify it's a number
                int(thread_number)
        except (ValueError, IndexError):
            await interaction.response.send_message(
                "❌ Could not determine thread number from thread name. Thread name should be in format 'support-<number>'.",
                ephemeral=True
            )
            return
        
        mentor_role_name = RoleNames.MENTOR_XI
        mentor_role_id = get_role_id(mentor_role_name)
        
        if mentor_role_id:
            mentor_role = interaction.guild.get_role(int(mentor_role_id))
        else:
            mentor_role = discord.utils.get(interaction.guild.roles, name=mentor_role_name)
        
        category = discord.utils.get(interaction.guild.categories, name="--- SwampHacks XI (Support-VCs) ---")
        vc_author = interaction.user
        if not mentor_role:
            await interaction.response.send_message(
                f"Error: Could not find the **{mentor_role_name}** role. Please create it before using this command.",
                ephemeral=True
            )
            return
        if not category:
            await interaction.response.send_message(
                "Error: Could not find the --- SwampHacks XI (Support-VCs) --- category. Please create it before using this command.",
                ephemeral=True
            )
            return
        
        # Check if VC with this number already exists
        # Use "thread-VC-{number}" format to avoid conflicts with directly created VCs
        vc_name = f"thread-VC-{thread_number}"
        existing_vc = discord.utils.get(category.voice_channels, name=vc_name)
        if existing_vc:
            await interaction.response.send_message(
                f"❌ A voice channel with name {vc_name} already exists. Please delete it first or use a different thread.",
                ephemeral=True
            )
            return
        
        overwrites = {
            interaction.guild.default_role: discord.PermissionOverwrite(view_channel=False, connect=False),
            mentor_role: discord.PermissionOverwrite(view_channel=True, connect=True),
            vc_author: discord.PermissionOverwrite(view_channel=True, connect=True),
        }
        # Create the voice channel and get the channel object
        voice_channel = await interaction.guild.create_voice_channel(
            name=vc_name,
            category=category,
            user_limit=4,
            reason="Support VC created for mentor and inquirer",
            overwrites=overwrites
        )
        
        text_channel_embed = discord.Embed(
            title=f"Voice channel created: {vc_name}",
            description=f"Here is the voice channel you requested.",
            color=discord.Color.green(),
        )
        
        # Try to send a message in the voice channel's chat (if available)
        text_channel = interaction.guild.get_channel(voice_channel.id)
        if text_channel:
            await text_channel.send(content=f"{vc_author.mention}")
            await text_channel.send(embed=text_channel_embed)
        else:
            print("Voice channel does not have an associated text channel.")
            return

        # Store the association between thread and VC
        thread_vc_mapping[thread.id] = voice_channel.id
        
        # ping the user who created the thread
        await interaction.response.send_message(
            f"Voice channel created: {voice_channel.mention}",
            ephemeral=True
        )

    @app_commands.command(name="grant_vc_access", description="Grant a user access to a voice channel")
    @app_commands.describe(user="Grant a user access to a voice channel")
    @has_bot_full_access_or_hacker()
    async def grant_vc_access(self, interaction: discord.Interaction, user: discord.Member) -> None:
        """
        Grant a user access to a voice channel, can only be used in a voice channel under the --- SwampHacks XI (Support-VCs) --- category
        
        Args:
            interaction: The interaction that triggered this command
            user: The user to grant access to
        """
        # TODO: We may want to allow this comamnd to be used in any channel since only mods can use it, but then we need to add a parameter to specify the voice channel
        # first ensure command is being executed in a voice channel
        if not isinstance(interaction.channel, discord.VoiceChannel):
            await interaction.response.send_message(
                "This command can only be used in a voice channel.", ephemeral=True
            )
            return
        # next ensure the channel is in the support category specifically
        if not interaction.channel.category or interaction.channel.category.name != "--- SwampHacks XI (Support-VCs) ---":
            await interaction.response.send_message('This command can only be used in the "--- SwampHacks XI (Support-VCs) ---" category.', ephemeral=True)
            return
        
        # check if user already has access to the voice channel
        overwrites = interaction.channel.overwrites_for(user)
        if overwrites.connect is True:
            await interaction.response.send_message(f"{user.mention} already has access to this voice channel.", ephemeral=True)
            return
        
        # try to grant access to the voice channel
        try:
            await interaction.channel.set_permissions(user, connect=True, view_channel=True)
            await interaction.channel.send(f"{user.mention} has been granted access to this voice channel.")
        except discord.Forbidden:
            await interaction.response.send_message("I don't have permission to grant access to this voice channel.", ephemeral=True)
        except Exception as e:
            await interaction.response.send_message(f"An error occurred: {str(e)}", ephemeral=True)
            
    @app_commands.command(name="ask", description="Ask LLM a question")
    @app_commands.describe(prompt="Your question for LLM")
    async def ask(self, interaction: discord.Interaction, prompt: str):
        """Ask LLM a question and get a response."""
        await interaction.response.defer(thinking=True)
        answer = llm_response(prompt)
        await interaction.followup.send(answer, ephemeral=True)
    
    @app_commands.command(name="create_announcement", description="Create an announcement message that the bot will send")
    @app_commands.describe(
        message="The announcement message to send (supports Discord markdown formatting). Use @username or @rolename to mention users/roles.",
        attachment="Optional image or file attachment to include with the announcement"
    )
    @has_bot_full_access()
    async def create_announcement(
        self, 
        interaction: discord.Interaction, 
        message: str,
        attachment: discord.Attachment = None
    ) -> None:
        """Create an announcement message that the bot will send in the current channel
        
        Args:
            interaction: The interaction that triggered this command
            message: The announcement message to send (supports Discord markdown formatting)
            attachment: Optional image or file attachment to include with the announcement
        """
        await interaction.response.defer(ephemeral=True)
        
        file = None
        if attachment:
            file = await attachment.to_file()
        
        processed_message = message.encode().decode('unicode_escape')
        
        if interaction.guild:
            mention_pattern = r'@(\w+)'
            
            def replace_mention(match):
                name = match.group(1).lower()
                member = discord.utils.find(lambda m: m.name.lower() == name, interaction.guild.members)
                if member:
                    return member.mention
                member = discord.utils.find(lambda m: m.display_name.lower() == name, interaction.guild.members)
                if member:
                    return member.mention
                role = discord.utils.find(lambda r: r.name.lower() == name, interaction.guild.roles)
                if role:
                    return role.mention
                return match.group(0)
            
            processed_message = re.sub(mention_pattern, replace_mention, processed_message)
        
        allowed_mentions = discord.AllowedMentions(everyone=True, users=True, roles=True)
        
        if file:
            await interaction.channel.send(processed_message, file=file, allowed_mentions=allowed_mentions)
        else:
            await interaction.channel.send(processed_message, allowed_mentions=allowed_mentions)
        
        await interaction.delete_original_response()
        
    
    @app_commands.command(
        name="assign_hacker_roles",
        description="Assign hacker role to all attendees from API using webhook"
    )
    @app_commands.describe(
        event_id="UUID of event",
        role="Discord role to assign to attendees",
        test_mode="Test mode: only process first 5 users and stop on first error (default: False)"
    )
    @requires_admin_or_moderator()
    async def assign_hacker_roles(self, interaction: discord.Interaction, event_id: str, role: discord.Role, test_mode: bool = False) -> None:
        """Assign role to all attendees from API using webhook
        
        Args:
            interaction: The interaction that triggered this command
            event_id: UUID of event
            role: Discord role to assign to attendees
        """
        
        await interaction.response.defer(ephemeral=True)
        
        try:
            guild_id = interaction.guild.id if interaction.guild else None
            if not guild_id:
                await interaction.followup.send("Error: Could not determine guild.", ephemeral=True)
                return
            
            api_url = os.getenv("API_URL", "https://api.swamphacks.com")
            
            session_cookie = os.getenv("SESSION_COOKIE")
            if not session_cookie:
                await interaction.followup.send("Error: SESSION_COOKIE is not set.", ephemeral=True)
                return
            
            webhook_url = os.getenv("WEBHOOK_URL")
            if not webhook_url:
                await interaction.followup.send("Error: WEBHOOK_URL is not set.", ephemeral=True)
                return
            
            attendees = await get_attendees_for_event(api_url, session_cookie, event_id)
            if not attendees:
                await interaction.followup.send("Error: No attendees found for event.", ephemeral=True)
                return
            
            # Limit to 5 users in test mode
            if test_mode:
                attendees = attendees[:5]
                await interaction.followup.send(f"🧪 **TEST MODE**: Processing first 5 attendees only. Will stop on first error.", ephemeral=True)
            
            total_attendees = len(attendees)
            
            # Send initial progress message
            mode_text = "🧪 TEST MODE: " if test_mode else ""
            await interaction.followup.send(f"{mode_text}🔄 Processing {total_attendees} attendees in chunks of 20...", ephemeral=True)
            
            # Progress callback to send updates
            async def progress_update(current: int, total: int):
                if current % 20 == 0 or current == total:  # Update every 20 or at the end
                    await interaction.followup.send(f"⏳ Progress: {current}/{total} processed...", ephemeral=True)
            
            newly_assigned, already_had, failed, errors = await assign_roles_to_attendees(
                webhook_url, 
                attendees, 
                role.name, 
                str(guild_id),
                chunk_size=20,
                progress_callback=progress_update,
                test_mode=test_mode
            )
            
            summary = format_assignment_summary(total_attendees, newly_assigned, already_had, failed, errors)
            if test_mode and errors:
                summary += f"\n\n⚠️ **Test stopped early due to error.** Fix the issue before running full assignment."
            await interaction.followup.send(summary, ephemeral=True)
            
        except Exception as e:
            await interaction.followup.send(f"An error occurred: {str(e)}", ephemeral=True)
            
    @app_commands.command(name="remove_role_from_all", description="Remove a specific role from all members in the server")
    @app_commands.describe(role="The role to remove from all members")
    @requires_admin_or_moderator()
    async def remove_role_from_all(self, interaction: discord.Interaction, role: discord.Role) -> None:
        """Remove a specific role from all members in the server
        """
        await interaction.response.defer(ephemeral=True)
        try:
            guild = interaction.guild
            if not guild:
                await interaction.followup.send("Error: Could not determine guild.", ephemeral=True)
                return
            
            await interaction.followup.send(
                f"Fetching all members and removing **{role.name}** role... This may take a moment.",
                ephemeral=True
            )
            
            members_with_role = [member for member in guild.members if role in member.roles]
            
            if not members_with_role:
                await interaction.followup.send(
                    f"No members found with the **{role.name}** role.",
                    ephemeral=True
                )
                return
            
            removed = 0
            failed = 0
            errors = []
            
            for member in members_with_role:
                try:
                    await member.remove_roles(role, reason=f'Role removal via command by {interaction.user}')
                    removed += 1
                except discord.Forbidden:
                    failed += 1
                    errors.append(f"Permission denied for {member.mention}")
                except discord.HTTPException as e:
                    failed += 1
                    errors.append(f"Error removing role from {member.mention}: {str(e)}")
                except Exception as e:
                    failed += 1
                    errors.append(f"Unexpected error for {member.mention}: {str(e)}")
            
            message = f"**Role Removal Complete**\n\n"
            message += f"**Summary:**\n"
            message += f"- Total members with **{role.name}** role: {len(members_with_role)}\n"
            message += f"- Roles removed successfully: {removed}\n"
            message += f"- Failed removals: {failed}\n"
            
            if errors:
                message += f"\n**Errors ({len(errors)}):**\n"
                for error in errors[:10]:
                    message += f"- {error}\n"
                if len(errors) > 10:
                    message += f"- ... and {len(errors) - 10} more errors\n"
            
            await interaction.followup.send(message, ephemeral=True)
        except Exception as e:
            await interaction.followup.send(
                f"An error occurred: {str(e)}",
                ephemeral=True
            )
    
    @commands.Cog.listener()
    async def on_member_join(self, member: discord.Member) -> None:
        """Automatically assign roles wheen a member joins the server"""
        
        logger = logging.getLogger(__name__)
        guild_id = member.guild.id
        if not guild_id:
            return
        
        api_url = os.getenv("API_URL", "https://api.swamphacks.com")
        session_cookie =  os.getenv("SESSION_COOKIE")
        event_id = os.getenv("EVENT_ID")
        
        if not session_cookie:
            logger.error("SESSION_COOKIE is not set")
            return
        
        if not event_id:
            logger.error("EVENT_ID is not set")
            return
        
        try:
            async with aiohttp.ClientSession() as session:
                headers = {"Cookie": f"sh_session_id={session_cookie}"}
                async with session.get(
                    f"{api_url}/events/{event_id}/discord/{member.id}",
                    headers=headers
                ) as response:
                    if response.status == 200:
                        data = await response.json()
                        event_role = data.get("role")
                        discord_role_name = "Hacker"
                        if event_role == "attendee":
                            role = discord.utils.get(member.guild.roles, name=discord_role_name)
                            if role:
                                await member.add_roles(role, reason="Auto assigned: User has attendee role")
                                logger.info(f"Auto assigned {discord_role_name} role to {member.name} ({member.id})")
                    elif response.status == 404:
                        pass
                    else:
                        logger.error(f"Unexpected response status {response.status} when checking role for {member.name} ({member.id})")
                        

        except Exception as e:
            logger.error(f"Error assigning roles: {str(e)}")
async def setup(bot: commands.Bot) -> None:
    """Add the General cog to the bot
    
    Args:
        bot: Discord bot instance
    """
    await bot.add_cog(General(bot))
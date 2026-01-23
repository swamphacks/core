from discord.ui import View, Button
from discord import ButtonStyle, Interaction
import discord
from discord.errors import NotFound
from utils.mentor_functions import set_busy_mentor, set_available_mentor
from components.ticket_state import claimed_tickets
from chatbot.llm import summarize_text

class SupportThreadButtons(View):
    """
    View for support thread buttons which is used in ThreadSupportModal to add buttons onto the embed sent to the reports channel.

    """
    def __init__(self, thread: discord.Thread, description_input: discord.ui.TextInput) -> None:
        """
        Initializes the SupportThreadButtons view with the given thread and description input.

        Args:
            thread (discord.Thread): The support thread to which the buttons will be added.
            description_input (discord.ui.TextInput): The text input containing the description of the issue.
        """
        super().__init__(timeout=None)
        self.thread = thread
        self.description_input = description_input
        self.add_item(ClaimThreadButton(thread, description_input))
        self.add_item(CloseThreadButton(thread, description_input))


class CloseThreadButton(Button):
    """Button to close a support thread, archive it, and lock it."""
    def __init__(self, thread: discord.Thread, description_input: discord.ui.TextInput, thread_author=None, threadID=None):
        """
        Initializes the CloseThreadButton with the given thread and description input.

        Args:
            thread (discord.Thread): The support thread to be closed.
            description_input (discord.ui.TextInput): The text input containing the description of the issue.
            thread_author (discord.Member, optional): The user who created the thread. If provided, only creator and mentors can close.
            threadID: Optional thread ID parameter.
        """
        super().__init__(label="Close Thread", style=ButtonStyle.primary, custom_id="close_thread", emoji="❌")
        self.thread = thread
        self.threadID = threadID
        self.description_input = description_input
        self.thread_author = thread_author

    async def callback(self, interaction: Interaction):
        # Check if user is authorized to close (creator or mentor)
        if self.thread_author:
            from utils.roles_config import get_acceptable_roles
            
            is_creator = interaction.user.id == self.thread_author.id
            is_mentor = False
            
            # Check if user has mentor/acceptable role
            if interaction.guild and interaction.user:
                member = interaction.guild.get_member(interaction.user.id)
                if member:
                    acceptable_roles = get_acceptable_roles()
                    is_mentor = any(role.name in acceptable_roles for role in member.roles)
            
            if not is_creator and not is_mentor:
                await interaction.response.send_message(
                    "❌ Only the thread creator or mentors can close this thread.",
                    ephemeral=True
                )
                return
        
        claimed_tickets.pop(self.thread.id, None)
        # print(claimed_tickets)
        
        # get channels
        reports_channel = discord.utils.get(interaction.guild.channels, name="reports")
        if not reports_channel:
            await interaction.response.send_message("❌ Reports channel not found.", ephemeral=True)
            return
        archived_threads_channel = discord.utils.get(interaction.guild.channels, name="archived-support-threads")
        if not archived_threads_channel:
            await interaction.response.send_message("❌ Archived threads channel not found.", ephemeral=True)
            return
        bot_avatar_url = interaction.client.user.avatar.url if interaction.client.user.avatar else None
        
        # Send response first to avoid timeout
        await interaction.response.defer(ephemeral=True)
        
        try:
            # rename the thread to get new title
            prefix = "archived-"
            title = ""
            if interaction.message.embeds:
                embed_title = interaction.message.embeds[0].title
                # Handle different embed title formats
                if embed_title.startswith("💬 New Thread Request: "):
                    # Reports channel embed format
                    title = embed_title[22:]  # Remove "💬 New Thread Request: "
                elif embed_title.startswith("Request: "):
                    # Thread embed format
                    title = embed_title[9:]  # Remove "Request: "
                else:
                    # Fallback: use the full title
                    title = embed_title
                
                # Trim to fit Discord's thread name limit (100 chars) minus prefix
                max_length = 100 - len(prefix)
                if len(title) > max_length:
                    title = title[:max_length]
            else:
                title = self.thread.name
            new_name = f"archived-{title}"
            
            ### FUNCTIONALITY FOR LLM SUMMARIZATION ###
            # fetch all messages
            messages = []
            async for msg in self.thread.history(limit=None):
                if msg.content:
                    messages.append(msg.content)
            combined_text = "\n".join(messages)
            # get summary of thread
            summary = await summarize_text(combined_text)
            embed = discord.Embed(
                title=f"💡 Summary of thread: {self.thread.mention}",
                description=f"**Title**: {title}\n{summary}",
                color=discord.Color.blue()
            )
            embed.set_footer(
                text="Powered by SwampHacks",
                icon_url=bot_avatar_url
            )

            # send the summary to the archived threads channel
            await archived_threads_channel.send(embed=embed)

            # Delete associated VC if it exists
            from components.ticket_state import thread_vc_mapping
            try:
                if self.thread.id in thread_vc_mapping:
                    vc_id = thread_vc_mapping[self.thread.id]
                    vc = interaction.guild.get_channel(vc_id)
                    if vc and isinstance(vc, discord.VoiceChannel):
                        try:
                            await vc.delete(reason=f"Thread {self.thread.id} was closed")
                        except Exception as e:
                            print(f"Failed to delete VC {vc_id} associated with thread {self.thread.id}: {e}")
                    # Remove from mapping even if VC doesn't exist anymore
                    del thread_vc_mapping[self.thread.id]
            except Exception as e:
                # If VC deletion fails, log but continue with thread closure
                print(f"Error handling VC deletion for thread {self.thread.id}: {e}")
            
            try:
                # Fetch members explicitly and collect all member IDs
                members_to_remove = set()
                
                # Add all members from thread.members
                for thread_member in self.thread.members:
                    members_to_remove.add(thread_member.id)
                
                # Also explicitly include the interaction user and thread author
                if interaction.user:
                    members_to_remove.add(interaction.user.id)
                if self.thread_author:
                    members_to_remove.add(self.thread_author.id)
                
                # Try to fetch members if the list seems incomplete
                if len(members_to_remove) == 0:
                    try:
                        async for member in self.thread.fetch_members():
                            members_to_remove.add(member.id)
                    except Exception as fetch_error:
                        print(f"Failed to fetch members: {fetch_error}")
                
                for member_id in members_to_remove:
                    try:
                        member = interaction.guild.get_member(member_id)
                        if member:
                            await self.thread.remove_user(member)
                    except Exception as remove_error:
                        print(f"Failed to remove member {member_id} from thread {self.thread.id}: {remove_error}")
            except Exception as e:
                print(f"Error removing members from thread {self.thread.id}: {e}")
            
            try:
                await self.thread.edit(name=new_name)
                await self.thread.edit(archived=True, locked=True)
            except Exception as e:
                print(f"Failed to archive thread {self.thread.id}: {e}")
                try:
                    await self.thread.edit(archived=True, locked=True)
                except Exception as e2:
                    print(f"Failed to archive thread {self.thread.id} (fallback attempt): {e2}")
                    await interaction.followup.send(f"⚠️ Thread archiving encountered an error, but the thread should be closed.", ephemeral=True)
                    return
            
            await interaction.followup.send(f"Thread has been archived and locked.", ephemeral=True)

            
            # Set mentor status - only mark as available if they have no more tickets
            mentor_ticket_count = sum(1 for mentor_id in claimed_tickets.values() if mentor_id == interaction.user.id)
            if mentor_ticket_count == 0:
                await set_available_mentor(interaction.user, True)
                await set_busy_mentor(interaction.user, False)
            
            # trim description
            description = self.description_input.value
            shortened_description = ""
            if len(description) > 200:
                shortened_description = description[:200] + "..."
            else:
                shortened_description = description
            
            # Create disabled view for updating messages
            new_view = SupportThreadButtons(self.thread, self.description_input)
            # disable all buttons in the view
            for item in new_view.children:
                item.disabled = True
            
            # Update the thread message (if button was clicked from thread)
            message = interaction.message
            if message:
                embed = message.embeds[0] if message.embeds else None
                if embed:
                    new_embed = embed.copy()
                    new_embed.description = f"Issue: {shortened_description}\n\nActions: Thread closed by {interaction.user.display_name}."
                    new_embed.color = discord.Color.red()
                    try:
                        await message.edit(embed=new_embed, view=new_view, allowed_mentions=discord.AllowedMentions.none())
                    except Exception as e:
                        print(f"Failed to update thread message: {e}")
                else:
                    try:
                        await message.edit(view=new_view, allowed_mentions=discord.AllowedMentions.none())
                    except Exception as e:
                        print(f"Failed to update thread message view: {e}")
            
            # Also update the reports channel message
            # Find the reports channel message by searching for messages with the thread mention or matching title
            try:
                reports_channel = discord.utils.get(interaction.guild.channels, name="reports")
                if reports_channel:
                    # Search for the message in reports channel
                    # Look for messages that mention this thread or have matching embed title
                    async for reports_message in reports_channel.history(limit=100):
                        if reports_message.embeds:
                            embed_title = reports_message.embeds[0].title
                            # Check if this is the reports channel message for this thread
                            if (embed_title.startswith("💬 New Thread Request: ") and 
                                title in embed_title):
                                # Found the reports channel message
                                reports_embed = reports_message.embeds[0].copy()
                                reports_embed.description = f"Issue: {shortened_description}\n\nActions: {interaction.user.mention} closed {self.thread.name}."
                                reports_embed.color = discord.Color.red()
                                try:
                                    await reports_message.edit(embed=reports_embed, view=new_view)
                                except Exception as e:
                                    print(f"Failed to update reports channel message: {e}")
                                break
            except Exception as e:
                print(f"Failed to find/update reports channel message: {e}")
        except NotFound:
            await interaction.followup.send(
                "This support thread no longer exists.",
                ephemeral=True
            )
        except Exception as e:
            await interaction.followup.send(f"Failed to archive the support thread. Error: {e}", ephemeral=True)

class JoinThreadButton(Button):
    """Button to join a claimed support thread for additional mentor assistance."""
    def __init__(self, thread: discord.Thread):
        super().__init__(label="Join Thread", style=ButtonStyle.primary, custom_id="join_thread", emoji="➡️")
        self.thread = thread

    async def callback(self, interaction: Interaction):
        try:
            await self.thread.add_user(interaction.user)
            await interaction.response.send_message(f"You've joined the thread {self.thread.mention}", ephemeral=True)
        except NotFound:
            await interaction.response.send_message("❌ This thread no longer exists.", ephemeral=True)
        except Exception as e:
            await interaction.response.send_message(f"❌ Error joining thread: {str(e)}", ephemeral=True)

class JoinThreadButton(Button):
    """Button to join a claimed support thread for additional mentor assistance."""
    def __init__(self, thread: discord.Thread):
        super().__init__(label="Join Thread", style=ButtonStyle.primary, custom_id="join_thread", emoji="➡️")
        self.thread = thread

    async def callback(self, interaction: Interaction):
        try:
            await self.thread.add_user(interaction.user)
            await interaction.response.send_message(f"You've joined the thread {self.thread.mention}", ephemeral=True)
        except NotFound:
            await interaction.response.send_message("❌ This thread no longer exists.", ephemeral=True)
        except Exception as e:
            await interaction.response.send_message(f"❌ Error joining thread: {str(e)}", ephemeral=True)

class ClaimThreadButton(Button):
    """Button to claim a support thread and add the mentor to it."""
    def __init__(self, thread: discord.Thread, description_input: discord.ui.TextInput):
        """
        Initializes the ClaimThreadButton with the given thread and description input.

        Args:
            thread (discord.Thread): The support thread to be claimed.
            description_input (discord.ui.TextInput): The text input containing the description of the issue.
        """
        super().__init__(label="Claim Thread", style=ButtonStyle.primary, custom_id="claim_thread", emoji="📥")
        self.thread = thread
        self.description_input = description_input

    async def callback(self, interaction: Interaction):
        try:
            # Check if the thread is already claimed
            if claimed_tickets.get(self.thread.id):
                await interaction.response.send_message(
                    "This thread has already been claimed by another mentor.",
                    ephemeral=True
                )
                return
            
            # Mark as claimed
            claimed_tickets[self.thread.id] = interaction.user.id
            # print(claimed_tickets)
            
            
            await self.thread.add_user(interaction.user)
            await interaction.response.send_message(f"You've been added to the thread: {self.thread.mention}", ephemeral=True)
            
            await set_available_mentor(interaction.user, False)
            await set_busy_mentor(interaction.user, True)
            
            
            # Edit the original message to show disabled claim button, close button, and join button
            message = interaction.message
            new_view = View(timeout=None)
            
            # Add disabled claim button
            claim_button = ClaimThreadButton(self.thread, self.description_input)
            claim_button.disabled = True
            new_view.add_item(claim_button)
            
            # Add close button in the middle
            new_view.add_item(CloseThreadButton(self.thread, self.description_input))
            
            # Add join button on the right
            new_view.add_item(JoinThreadButton(self.thread))

            # Copy the original embed and update its description
            embed = message.embeds[0] if message.embeds else None
            if embed:
                new_embed = embed.copy()
                new_embed.description = f"Issue: {self.description_input.value}\n\nActions: {interaction.user.mention} claimed the thread. Please join the thread to assist the member."
                await message.edit(embed=new_embed, view=new_view)
            else:
                await message.edit(view=new_view)
            
        except NotFound:
            await interaction.response.send_message(
                "This support thread no longer exists.",
                ephemeral=True
            )
        except Exception as e:
            await interaction.response.send_message(f"Failed to notify you about the support thread. Error: {e}", ephemeral=True)
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
    def __init__(self, thread: discord.Thread, description_input: discord.ui.TextInput, threadID=None):
        """
        Initializes the CloseThreadButton with the given thread and description input.

        Args:
            thread (discord.Thread): The support thread to be closed.
            description_input (discord.ui.TextInput): The text input containing the description of the issue.
        """
        super().__init__(label="Close Thread", style=ButtonStyle.primary, custom_id="close_thread", emoji="❌")
        self.thread = thread
        self.threadID = threadID
        self.description_input = description_input

    async def callback(self, interaction: Interaction):
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
        bot_avatar_url = interaction.client.user.avatar.url if interaction.client.user.avatar.url else None
        
        try:
            # rename the thread to get new title
            prefix = "archived-"
            title = ""
            if interaction.message.embeds:
                title = interaction.message.embeds[0].title
                trimmed_title = title[22:100 - len(prefix)]
                title = trimmed_title
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

            # archive and lock the thread
            await interaction.response.send_message(f"Thread: {self.thread.mention} has been archived and locked.", ephemeral=True)
            await self.thread.edit(name=new_name,archived=True, locked=True)

            
            # Set mentor status - only mark as available if they have no more tickets
            mentor_ticket_count = sum(1 for mentor_id in claimed_tickets.values() if mentor_id == interaction.user.id)
            if mentor_ticket_count == 0:
                await set_available_mentor(interaction.user, True)
                await set_busy_mentor(interaction.user, False)
            
            # edit original message to disable claim button
            message = interaction.message
            if not message:
                await interaction.response.send_message(
                    "Message not found.",
                    ephemeral=True
                )
                return
            new_view = SupportThreadButtons(self.thread, self.description_input)
            # disable all buttons in the view
            for item in new_view.children:
                item.disabled = True
            # copy the original embed and update its title and description
            embed = message.embeds[0] if message.embeds else None
            
            # trim description
            description = self.description_input.value
            shortened_description = ""
            if len(description) > 200:
                shortened_description = description[:200] + "..."
            else:
                shortened_description = description
            if embed:
                new_embed = embed.copy()
                new_embed.description = f"Issue: {shortened_description}\n\nActions: {interaction.user.mention} closed {self.thread.name}."
                new_embed.color = discord.Color.red()
                await message.edit(embed=new_embed, view=new_view)
            else:
                await message.edit(view=new_view)
        except NotFound:
            await interaction.response.send_message(
                "This support thread no longer exists.",
                ephemeral=True
            )
        except Exception as e:
            await interaction.response.send_message(f"Failed to archive the support thread. Error: {e}", ephemeral=True)

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
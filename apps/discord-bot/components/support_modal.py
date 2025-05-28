import discord
from discord.ui import Modal, TextInput
from discord import TextStyle
from components.join_thread_button import JoinThreadButton
from components.close_thread_button import CloseThreadButton
from utils.get_next_thread_name import get_next_thread_name

class SupportModal(Modal, title="Support Inquiry"):
    """
    Support modal for creating a support thread and embeds.
    
    This modal provides a form interface for users to submit support requests.
    It creates a private thread in the support channel and notifies staff members
    through the reports channel.
    """
    def __init__(self) -> None:
        """
        Initialize the support modal with title and description input fields.
        
        Creates two text input fields:
        - Title: A short title for the support request (max 100 characters)
        - Description: A detailed description of the issue (paragraph style)
        """
        super().__init__()
        self.title_input = TextInput(label="Title", max_length=100)
        self.description_input = TextInput(label="Describe your issue", style=TextStyle.paragraph)
        self.add_item(self.title_input)
        self.add_item(self.description_input)

    async def on_submit(self, interaction: discord.Interaction) -> None:
        """
        Handle the submission of the support request form.
        
        This method:
        1. Validates the existence of required channels and roles
        2. Creates a private thread in the support channel for the user and moderators to discuss the issue
        3. Sends an embed with the support request details to both the thread and reports channel
        4. Notifies staff members and provides them with a button to join the thread
        
        Args:
            interaction: Discord interaction object
            
        Note:
            The method will send error messages if required channels or roles are not found.
            Staff members are notified with a soft ping (hidden mention) in the reports channel.
        """
        reports_channel = discord.utils.get(interaction.guild.channels, name="reports")
        support_channel = discord.utils.get(interaction.guild.channels, name="support")
        thread_author = interaction.user
        mod_role = discord.utils.get(interaction.guild.roles, name="Moderator")
        if not mod_role:
            await interaction.response.send_message("Error: Could not find the moderator role. Please contact an administrator.", ephemeral=True)
            return
        
        # check if the channels and roles exist
        if not reports_channel:
            await interaction.response.send_message(
                "Error: Could not find the reports channel. Please contact an administrator.",
                ephemeral=True
            )
            return
        if not support_channel:
            await interaction.response.send_message("Error: Could not find the support channel. Please contact an administrator.", ephemeral=True)
            return
        
        # truncate description in case it's too long
        description = self.description_input.value
        if len(description) > 1000:
            description = description[:1000] + "..."
            
        # create embed for reports channel
        reports_embed = discord.Embed(
            title=f"⭐ New Request: {self.title_input.value}",
            description=f"Description:{self.description_input.value}",
            color=discord.Color.blue()
        )
        reports_embed.add_field(name="Opened by", value=f"{thread_author.mention}\n", inline=True)

        
        # create the thread with the next available name and add the initialuser to the thread
        thread_name = get_next_thread_name(support_channel)
        thread = await support_channel.create_thread(
            name=thread_name,
            type=discord.ChannelType.private_thread,
            reason=f"Support request from {thread_author}",
            auto_archive_duration=1440
        )
        await thread.add_user(thread_author)
        
        # send initial message as embed in thread with inquiry details
        thread_embed = discord.Embed(
            title=f"Request: {self.title_input.value}",
            description=f"Description: {self.description_input.value}\n\n✅ Thank you for your request, we will be with you shortly!",
            color=discord.Color.green(),
        )
        await thread.send(embed=thread_embed, view=CloseThreadButton(thread, thread_author))
        
        # send the embed to the reports channel
        await interaction.response.send_message(
            f"Support thread created: {thread.mention}",
            ephemeral=True
        )
                
        # soft ping staff and send the embed to the reports channel
        await reports_channel.send(
            content=f"||{mod_role.mention}||",
            embed=reports_embed,
            view=JoinThreadButton(thread),
            allowed_mentions=discord.AllowedMentions(roles=[])
        )
        
        await interaction.followup.send(
            "Thank you! Your issue has been submitted and a thread has been created for you.",
            ephemeral=True,
        )
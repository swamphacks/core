import discord
from discord.ui import Modal, TextInput, View
from discord import TextStyle
from components.support_thread_buttons import SupportThreadButtons, CloseThreadButton
from utils.get_next_thread_name import get_next_thread_name
from utils.get_next_support_vc_name import get_next_support_vc_name
from components.support_vc_buttons import SupportVCButtons
from utils.mentor_functions import get_available_mentors
from components.support_vc_buttons import SupportVCButtons
from components.mentor_ping_state import pinged_mentors
import random

class ThreadSupportModal(Modal, title="Support Inquiry"):
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
            
        # get available mentors
        available_mentors = get_available_mentors(mod_role)
        while selected_mentor not in pinged_mentors:
            # if all have been pinged reset the set
            if not available_mentors:
                pinged_mentors.clear()
                available_mentors = get_available_mentors(mod_role)
            selected_mentor = random.choice(available_mentors) if available_mentors else None
        
        
        
        # not_yet_pinged = [mentor for mentor in available_mentors]
        
        # # if all have been pinged reset the set
        # if not not_yet_pinged:
        #     pinged_mentors.clear()
        #     not_yet_pinged = available_mentors
        # print(pinged_mentors)
        
        # if not_yet_pinged:
        #     selected_mentor = random.choice(not_yet_pinged)
        #     pinged_mentors.add(selected_mentor.id)
        #     action_text = f"{selected_mentor.mention} Please join the thread to assist the hacker."
        # else:
        #     selected_mentor = None
        #     action_text = "No mentors available at this time."

        
        # create the thread with the next available name and add the initialuser to the thread
        thread_name = get_next_thread_name(support_channel)
        thread = await support_channel.create_thread(
            name=thread_name,
            type=discord.ChannelType.private_thread,
            reason=f"Support request from {thread_author}",
            auto_archive_duration=1440
        )
        await thread.add_user(thread_author)
        
        
        close_button = CloseThreadButton(thread, self.description_input)
        close_button_view = View()
        close_button_view.add_item(close_button)
        
        # send initial message as embed in thread with inquiry details
        thread_embed = discord.Embed(
            title=f"Request: {self.title_input.value}",
            description=f"Description: {self.description_input.value}\n\n✅ Thank you for your request, we will be with you shortly!",
            color=discord.Color.green(),
        )
        await thread.send(embed=thread_embed)
        
        # create embed for reports channel
        reports_embed = discord.Embed(
            title=f"⭐ New Thread Request: {self.title_input.value}",
            description=f"Issue: {self.description_input.value}\n\nActions: {action_text}",
            color=discord.Color.blue(),
        )
        reports_embed.add_field(name="Opened by", value=f"{thread_author.mention}\n", inline=True)
                
        # soft ping staff and send the embed to the reports channel
        await reports_channel.send(
            content=f"||{mod_role.mention}||",
            embed=reports_embed,
            view=SupportThreadButtons(thread, self.description_input),
            allowed_mentions=discord.AllowedMentions(roles=True)
        )

        await interaction.response.send_message(
            f"Thank you! Your issue has been submitted and a thread has been created for you. Please check {thread.mention} for updates.",
            ephemeral=True,
        )
        
class VCSupportModal(Modal, title="VC Support Inquiry"):
    """
    Support modal for creating a support vc and embeds.
    
    This modal provides a form interface for users to submit support requests.
    It creates a private vc in the Support-VCs category and notifies staff members
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
        # add a dropdown to select vc or thread support
        # self.add_item(discord.ui.Select(
        #     placeholder="Select support type",
        #     options=[
        #         discord.SelectOption(label="Voice Channel", value="vc"),
        #         discord.SelectOption(label="Thread", value="thread")
        #     ]
        # ))
        # FIX! NOT WORKING YET
        self.add_item(self.title_input)
        self.add_item(self.description_input)

    async def on_submit(self, interaction: discord.Interaction) -> None:
        """
        Handle the submission of the support request form for voice channels.
        
        This method:
        1. Validates the existence of required channels and roles
        2. Creates a private voice channel in the Support-VCs category for the user and moderators to discuss the issue
        3. Sends an embed with the support request details to both the voice channel and reports channel
        4. Notifies staff members and provides them with a button to join the voice channel
        
        Args:
            interaction: Discord interaction object
            
        Note:
            The method will send error messages if required channels or roles are not found.
            Staff members are notified with a soft ping (hidden mention) in the reports channel.
        """
        reports_channel = discord.utils.get(interaction.guild.channels, name="reports")
        mod_role = discord.utils.get(interaction.guild.roles, name="Moderator")
        category = discord.utils.get(interaction.guild.categories, name="Support-VCs")
        vc_author = interaction.user
        
        if not mod_role:
            await interaction.response.send_message("Error: Could not find the moderator role. Please contact an administrator.", ephemeral=True)
            return

        if not reports_channel:
            await interaction.response.send_message(
                "Error: Could not find the reports channel. Please contact an administrator.",
                ephemeral=True
            )
            return
        if category is None:
            await interaction.response.send_message("Category 'Support-VCs' not found.", ephemeral=True)
            return
        
        description = self.description_input.value
        if len(description) > 1000:
            description = description[:1000] + "..."
        
        
        # give permissions to the moderator role and the user who clicked the button
        overwrites = {
            interaction.guild.default_role: discord.PermissionOverwrite(view_channel=False, connect=False),
            mod_role: discord.PermissionOverwrite(view_channel=True, connect=True),
            vc_author: discord.PermissionOverwrite(view_channel=True, connect=True),
        }
        
        vc_name = get_next_support_vc_name(category)
        
        # Create the voice channel and get the channel object
        voice_channel = await interaction.guild.create_voice_channel(
            name=vc_name,
            category=category,
            user_limit=4,
            reason="Support VC created for mentor and inquirer",
            overwrites=overwrites
        )

        # Try to send a message in the voice channel's chat (if available)
        text_channel = interaction.guild.get_channel(voice_channel.id)
        if text_channel:
            await text_channel.send(f"{vc_author.mention} Your support voice channel is ready! A mentor will join you soon.")
        else:
            print("Voice channel does not have an associated text channel.")
            return

        # ping the user who created the thread
        await interaction.response.send_message(
            f"Voice channel created: {voice_channel.mention}",
            ephemeral=True
        )

        # # get available mentors
        available_mentors = get_available_mentors(mod_role)
        
        not_yet_pinged = [mentor for mentor in available_mentors]
        
        # if all have been pinged reset the set
        if not not_yet_pinged:
            pinged_mentors.clear()
        
        if not_yet_pinged:
            selected_mentor = random.choice(not_yet_pinged)
            pinged_mentors.add(selected_mentor)
        else:
            selected_mentor = None

        if selected_mentor:
            action_text = f"{selected_mentor.mention} Please join the thread to assist the hacker."
        else:
            action_text = "No mentors available at this time."

        # create embed for reports channel
        reports_embed = discord.Embed(
            title=f"⭐ New Voice Channel Request",
            description=f"Issue: {self.description_input.value}\n\nActions: {action_text}",
            timestamp=discord.utils.utcnow(),
            color=discord.Color.blue()
        )
        reports_embed.add_field(name="Opened by", value=f"{vc_author.mention}\n", inline=True)

        await reports_channel.send(
            content=f"||{mod_role.mention}||",
            embed=reports_embed,
            view=SupportVCButtons(voice_channel, self.description_input),
            allowed_mentions=discord.AllowedMentions(roles=True)
        )
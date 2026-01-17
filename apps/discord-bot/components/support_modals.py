import discord
from discord.ui import Modal, TextInput, View
from discord import TextStyle
from components.support_thread_buttons import SupportThreadButtons, CloseThreadButton
from utils.get_next_thread_name import get_next_thread_name
from utils.get_next_support_vc_name import get_next_support_vc_name
from components.support_vc_buttons import SupportVCButtons
from utils.mentor_functions import get_available_mentors
from components.support_vc_buttons import SupportVCButtons
from components.mentor_ping_state import last_pinged_mentor_index

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
        self.description_input = TextInput(label="Describe your issue", style=TextStyle.paragraph, max_length=1000)
        self.add_item(self.title_input)
        self.add_item(self.description_input)

    async def on_submit(self, interaction: discord.Interaction) -> None:
        """
        Handle the submission of the support request form.
        
        This method:
        1. Validates the existence of required channels and roles
        2. Creates a private thread in the support channel for the user and mentors to discuss the issue
        3. Sends an embed with the support request details to both the thread and reports channel
        4. Notifies staff members and provides them with a button to join the thread
        
        Args:
            interaction: Discord interaction object
            
        Note:
            The method will send error messages if required channels or roles are not found.
            Staff members are notified with a soft ping (hidden mention) in the reports channel.
        """
        from utils.roles_config import RoleNames, get_role_id
        
        global last_pinged_mentor_index
        reports_channel = discord.utils.get(interaction.guild.channels, name="reports")
        support_channel = discord.utils.get(interaction.guild.channels, name="support")
        thread_author = interaction.user
        mod_role_name = RoleNames.MODERATOR
        mod_role_id = get_role_id(mod_role_name)
        
        if mod_role_id:
            mod_role = interaction.guild.get_role(int(mod_role_id))
        else:
            mod_role = discord.utils.get(interaction.guild.roles, name=mod_role_name)
        
        if not mod_role:
            await interaction.response.send_message(f"Error: Could not find the **{mod_role_name}** role. Please create it before using this command.", ephemeral=True)
            return
        
        # check if the channels exist
        if not reports_channel:
            await interaction.response.send_message(
                "Error: Could not find the **reports** channel. Please create it before using this command.",
                ephemeral=True
            )
            return
        if not support_channel:
            await interaction.response.send_message("Error: Could not find the **support** channel. Please create it before using this command.", ephemeral=True)
            return
        
        # truncate description in case it's too long
        description = self.description_input.value
        shortened_description = ""
        if len(description) > 200:
            shortened_description = description[:200] + "..."
            
        # get available mentors from all acceptable roles
        available_mentors = get_available_mentors(interaction.guild)
        if not available_mentors:
            await interaction.response.send_message(
                "Error: No available mentors at this time. Please try again later.",
                ephemeral=True
            )
            return
        
        # implement round-robin pinging of mentors
        if last_pinged_mentor_index >= len(available_mentors):
            last_pinged_mentor_index = 0
        selected_mentor = available_mentors[last_pinged_mentor_index]
        last_pinged_mentor_index = (last_pinged_mentor_index + 1) % len(available_mentors)
        action_text = f"{selected_mentor.mention} Please join the thread to assist the user."
        
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
            description=f"Description: {description}\n\n✅ Thank you for your request, we will be with you shortly!",
            color=discord.Color.green(),
        )
        await thread.send(content=f"{thread_author.mention}")
        await thread.send(embed=thread_embed)
        
        # create embed for reports channel
        reports_embed = discord.Embed(
            title=f"💬 New Thread Request: {self.title_input.value}",
            description=f"Issue: {shortened_description}\n\nActions: {action_text}",
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
        It creates a private vc in the --- SwampHacks XI (Support-VCs) --- category and notifies staff members
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
        self.description_input = TextInput(label="Describe your issue", style=TextStyle.paragraph, max_length=1000)
        self.add_item(self.title_input)
        self.add_item(self.description_input)

    async def on_submit(self, interaction: discord.Interaction) -> None:
        """
        Handle the submission of the support request form for voice channels.
        
        This method:
        1. Validates the existence of required channels and roles
        2. Creates a private voice channel in the --- SwampHacks XI (Support-VCs) --- category for the user and mentors to discuss the issue
        3. Sends an embed with the support request details to both the voice channel and reports channel
        4. Notifies staff members and provides them with a button to join the voice channel
        
        Args:
            interaction: Discord interaction object
            
        Note:
            The method will send error messages if required channels or roles are not found.
            Staff members are notified with a soft ping (hidden mention) in the reports channel.
        """
        from utils.roles_config import RoleNames, get_role_id
        
        global last_pinged_mentor_index
        reports_channel = discord.utils.get(interaction.guild.channels, name="reports")
        mod_role_name = RoleNames.MODERATOR
        mod_role_id = get_role_id(mod_role_name)
        
        if mod_role_id:
            mod_role = interaction.guild.get_role(int(mod_role_id))
        else:
            mod_role = discord.utils.get(interaction.guild.roles, name=mod_role_name)
        
        category = discord.utils.get(interaction.guild.categories, name="--- SwampHacks XI (Support-VCs) ---")
        vc_author = interaction.user
        
        if not mod_role:
            await interaction.response.send_message(f"Error: Could not find the **{mod_role_name}** role. Please create it before using this command.", ephemeral=True)
            return

        if not reports_channel:
            await interaction.response.send_message(
                "Error: Could not find the **reports** channel. Please create it before using this command.",
                ephemeral=True
            )
            return
        if category is None:
            await interaction.response.send_message("Category **--- SwampHacks XI (Support-VCs) ---** not found.", ephemeral=True)
            return
        
         # truncate description in case it's too long
        description = self.description_input.value
        shortened_description = ""
        if len(description) > 200:
            shortened_description = description[:200] + "..."
        
        
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
        
        text_channel_embed = discord.Embed(
            title=f"Request: {self.title_input.value}",
            description=f"Description: {description}\n\n✅ Thank you for your request, we will be with you shortly!",
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

        # get available mentors from all acceptable roles
        available_mentors = get_available_mentors(interaction.guild)
        if not available_mentors:
            await interaction.response.send_message(
                "Error: No available mentors at this time. Please try again later.",
                ephemeral=True
            )
            return
        
        # ping the user who created the thread
        await interaction.response.send_message(
            f"Voice channel created: {voice_channel.mention}",
            ephemeral=True
        )
        
        if last_pinged_mentor_index >= len(available_mentors):
            last_pinged_mentor_index = 0
        selected_mentor = available_mentors[last_pinged_mentor_index]
        last_pinged_mentor_index = (last_pinged_mentor_index + 1) % len(available_mentors)
        action_text = f"{selected_mentor.mention} Please join the vc to assist the user."
        
        # create embed for reports channel
        reports_embed = discord.Embed(
            title=f"🎤 New Voice Channel Request: {self.title_input.value}",
            description=f"Issue: {shortened_description}\n\nActions: {action_text}",
            timestamp=discord.utils.utcnow(),
            color=discord.Color.purple()
        )
        reports_embed.add_field(name="Opened by", value=f"{vc_author.mention}\n", inline=True)

        await reports_channel.send(
            content=f"||{mod_role.mention}||",
            embed=reports_embed,
            view=SupportVCButtons(voice_channel, self.description_input),
            allowed_mentions=discord.AllowedMentions(roles=True)
        )
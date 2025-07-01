from discord.ui import View, Button
import discord
from discord import Interaction
from components.support_modals import ThreadSupportModal, VCSupportModal

class TicketView(View):
    def __init__(self):
        super().__init__(timeout=None)
        self.add_item(OpenThreadButton())
        self.add_item(OpenVoiceChatButton())
        self.add_item(ContactSupportButton())

class OpenThreadButton(Button):
    """
    Button to open the support modal and create a support thread
    """
    def __init__(self) -> None:
        """
        Initialize the OpenThreadButton

        Args:
            timeout: Timeout for the button in seconds
        """
        super().__init__(label="Chat in threads", style=discord.ButtonStyle.primary, emoji="💬")

    async def callback(self, interaction: Interaction):
        try:
            await interaction.response.send_modal(ThreadSupportModal())
        except Exception as e:
            await interaction.response.send_message(
                "Sorry, there was an error opening the support modal. Please try again later.",
                ephemeral=True
            )
            # Log the error for debugging
            print(f"Error in open_threads: {e}")

class OpenVoiceChatButton(Button):
    """
    Button to open private voice chat for mentor and inquirer
    """
    def __init__(self) -> None:
        """
        Initialize the OpenVoiceChatButton

        """
        super().__init__(
            label="Chat in VC",
            style=discord.ButtonStyle.primary,
            emoji="🎤"
        )

    async def callback(self, interaction: Interaction):
        try:
            await interaction.response.send_modal(VCSupportModal())
        except Exception as e:
            await interaction.response.send_message(
                "Sorry, there was an error opening the support modal. Please try again later.",
                ephemeral=True
            )
            # Log the error for debugging
            print(f"Error in open_threads: {e}")

class ContactSupportButton(Button):
    """
    Button to open private voice chat for mentor and inquirer
    """
    def __init__(self) -> None:
        """
        Initialize the ContactSupportButton
        """
        super().__init__(
            label="Contact Support",
            style=discord.ButtonStyle.primary,
            emoji="📞"
        )

    async def callback(self, interaction: Interaction):
        try:
            await interaction.response.send_modal(VCSupportModal())
        except Exception as e:
            await interaction.response.send_message(
                "Sorry, there was an error opening the support modal. Please try again later.",
                ephemeral=True
            )
            # Log the error for debugging
            print(f"Error in contact_support: {e}")
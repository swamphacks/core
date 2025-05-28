from discord.ui import View, Button
import discord
from components.support_modal import SupportModal
from discord import Interaction

class OpenSupportButton(View):
    """
    Button to open the support modal and create a support thread
    """
    def __init__(self) -> None:
        """
        Initialize the OpenSupportButton
        
        Args:
            timeout: Timeout for the button in seconds
        """
        super().__init__(timeout=300)  # 5 minute timeout

    @discord.ui.button(label="Contact Support", style=discord.ButtonStyle.primary, emoji="📞")
    async def open_ticket(self, interaction: Interaction, button: Button):
        try:
            await interaction.response.send_modal(SupportModal())
        except Exception as e:
            await interaction.response.send_message(
                "Sorry, there was an error opening the support modal. Please try again later.",
                ephemeral=True
            )
            # Log the error for debugging
            print(f"Error in open_ticket: {e}")
            
    async def on_timeout(self):
        # Disable the button when the view times out
        for item in self.children:
            item.disabled = True
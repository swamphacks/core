from discord.ui import View, Button
from discord import ButtonStyle, Interaction
import discord
from discord.errors import NotFound
from utils.checks import is_mod_slash

class JoinThreadButton(View):
    """
    Button to join a thread and notify the user that they have been added to the thread
    """
    def __init__(self, thread: discord.Thread) -> None:
        """
        Initialize the JoinThreadButton
        
        Args:
            thread: Thread object indicating the thread to join
        """
        super().__init__(timeout=None)
        self.thread = thread

    @is_mod_slash()
    @discord.ui.button(label="Join thread", style=ButtonStyle.primary, custom_id="join_support_thread", emoji="➕")
    async def join_thread(self, interaction: Interaction, button: Button):
        """
        Join the thread and notify the user that they have been added to the thread, restrict command to @moderators.
        
        Args:
            interaction: Discord interaction object
            button: Discord button object
        """
        try:
            await self.thread.add_user(interaction.user)
            await interaction.response.send_message(f"You've been added to the thread: {self.thread.mention}", ephemeral=True)
        except NotFound:
            await interaction.response.send_message(
                "This thread was closed and can no longer be joined.",
                ephemeral=True
            )
        except Exception as e:
            await interaction.response.send_message(f"Failed to add you to the thread. Error: {e}", ephemeral=True)
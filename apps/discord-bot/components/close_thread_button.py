from discord.ui import View, Button
import discord
from discord import Interaction

class CloseThreadButton(View):
    """
    Button to close a thread and notify the user that the thread has been closed
    """
    def __init__(self, thread: discord.Thread, thread_author: discord.Member) -> None:
        """
        Initialize the CloseThreadButton
        
        Args:
            thread: The thread to close
            thread_author: The author/creator of the thread
        """
        super().__init__()
        self.thread = thread
        self.thread_author = thread_author

    @discord.ui.button(label="Close Thread", style=discord.ButtonStyle.danger, emoji="🔒")
    async def close_thread(self, interaction: Interaction, button: Button) -> None:
        """
        Close and delete the thread and send embed to reports channel with support request details.

        Args:
            interaction: Discord interaction object
            button: Discord button object
        """
        
        # get thread id/number
        try:
            thread_number = self.thread.name.split("-")[1]
        except IndexError:
            await interaction.response.send_message("❌ Unable to parse thread number from thread name.", ephemeral=True)
            return
        
        # create embed to send to reports channel
        reports_channel = discord.utils.get(interaction.guild.channels, name="reports")
        if not reports_channel:
            await interaction.response.send_message("❌ Reports channel not found.", ephemeral=True)
            return
            
        closed_embed = discord.Embed(
            title=f"Support Request **#{thread_number}** Closed",
            color=discord.Color.red()
        )
        closed_embed.add_field(name="Thread Opened by", value=f"{self.thread_author.mention}", inline=True)
        closed_embed.add_field(name="Thread Closed by", value=f"{interaction.user.mention}", inline=True)
        closed_embed.add_field(name="Thread Opened at", value=f"{self.thread.created_at.strftime('%Y-%m-%d %H:%M:%S')}", inline=True)

        try:
            await reports_channel.send(embed=closed_embed)
            await self.thread.delete()
            await interaction.response.send_message("Thread has been closed and deleted.", ephemeral=True)
        except discord.Forbidden:
            await interaction.response.send_message("❌ Missing permissions to close thread or send to reports.", ephemeral=True)
        except Exception as e:
            await interaction.response.send_message(f"❌ An error occurred: {str(e)}", ephemeral=True)
from discord.ext import commands
from discord.ui import View, Button
import discord


class TicketView(View):
    """
    A view for creating tickets
    
    This view includes a button for creating a ticket
    """
    def __init__(self):
        super().__init__(timeout=None)

    @discord.ui.button(label="Open a Ticket", style=discord.ButtonStyle.primary)
    async def create_ticket(self, interaction: discord.Interaction):
        await interaction.response.send_message("Ticket created", ephemeral=True)
        

class Tickets(commands.Cog):
    """
    A cog for creating and managing tickets
    
    This cog includes commands for:
    - Creating tickets
    - Managing tickets
    - Closing tickets
    - Viewing ticket logs
    """
    def __init__(self, bot: commands.Bot) -> None:
        """Initialize the Tickets cog
        
        Args:
            bot: The bot instance
        """
        self.bot = bot
    
    @commands.command()
    async def ticketpanel(self, ctx):
        embed = discord.Embed(
            title="Open Ticket",
            description="Click the button to open a ticket",
            color=discord.Color.green(),
        )
        embed.set_footer(text="Powered by SwampHacksXI")
        await ctx.send(embed=embed, view=TicketView())
        

async def setup(bot: commands.Bot) -> None:
    await bot.add_cog(Tickets(bot))


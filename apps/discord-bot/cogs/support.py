from discord.ext import commands
from discord import app_commands, Interaction, Embed, Colour, ButtonStyle
from discord.ui import View, Button, button
from typing import Literal
from components.support_modal import SupportModal


class SupportView(View):
    """
    A view for creating tickets
    
    This view includes a button for creating a ticket
    """
    def __init__(self):
        super().__init__(timeout=None)

    @button(label="Contact Support", style=ButtonStyle.primary)
    async def create_report(self, interaction: Interaction, button: Button):
        # await interaction.response.send_message("Support ticket created", ephemeral=True)
        await interaction.response.send_modal(SupportModal())

class Support(commands.Cog):
    """
    A cog for creating and managing tickets
    
    This cog includes commands for:
    - Creating panels
    - Managing tickets
    - Closing tickets
    - Viewing ticket logs
    """
    def __init__(self, bot: commands.Bot) -> None:
        """Initialize the Support cog
        
        Args:
            bot: The bot instance
        """
        self.bot = bot
    
    @app_commands.command(name="create_panel", description="Create a support panel")
    @app_commands.describe(
        title="The title of the support panel",
        description="The description of the support panel",
        color="Choose the panel's color"
    )
    async def supportpanel(
        self,
        interaction: Interaction,
        title: str,
        description: str,
        color: Literal["red", "blue", "green", "purple", "orange"]
    ):
        color_map = {
            "red": Colour.red(),
            "blue": Colour.blue(),
            "green": Colour.green(),
            "purple": Colour.purple(),
            "orange": Colour.orange()
        }

        embed = Embed(
            title=title,
            description=description,
            color=color_map[color]
        )
        embed.set_footer(text="Powered by SwampHacksXI")
        await interaction.response.send_message(embed=embed, view=SupportView())

        

async def setup(bot: commands.Bot) -> None:
    await bot.add_cog(Support(bot))


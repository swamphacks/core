from discord.ext import commands
from discord import app_commands, Interaction, Embed, Colour
from typing import Literal
from components.open_support_button import OpenSupportButton
from utils.checks import is_mod_slash

class Support(commands.Cog):
    """
    A cog for creating and managing support requests
    
    This cog includes commands for:
    - Creating panels
    - Managing support requests and opening threads
    - Closing support requests and closing threads
    - Viewing support request logs
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
    @is_mod_slash()
    async def supportpanel(
        self,
        interaction: Interaction,
        title: str,
        description: str,
        color: Literal["red", "blue", "green", "purple", "orange"]
    ) -> None:
        """
        Create a support panel mainly used in the #support channel and restricted to @moderators.
        
        Args:
            interaction: The interaction object
            title: The title of the support panel
            description: The description of the support panel
            color: The color of the support panel
        """
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
        await interaction.response.defer(ephemeral=True)
        await interaction.delete_original_response()
        await interaction.channel.send(embed=embed, view=OpenSupportButton())

        

async def setup(bot: commands.Bot) -> None:
    await bot.add_cog(Support(bot))


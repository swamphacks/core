from discord.ext import commands
from discord import app_commands
import discord
from typing import Literal
from utils.checks import is_mod_slash


class General(commands.Cog):
    """A cog containing general utility commands for the server
    
    This cog includes commands for:
    - Basic server interactions
    - Role management
    - Fun commands
    """
    def __init__(self, bot: commands.Bot) -> None:
        """Initialize the General cog
        
        Args:
            bot: Discord bot instance
        """
        self.bot: commands.Bot = bot
    
    @commands.command()
    async def test(self, ctx: commands.Context) -> None:
        """Send a test message
        
        Args:
            ctx: The command context
        """
        await ctx.send("Testing")

    @app_commands.command(name="delete", description="Delete X amount of messages based on the number you provide")
    @app_commands.describe(
        amount="The amount of messages to delete"
    )
    @is_mod_slash()
    async def delete(
        self,
        interaction: discord.Interaction,
        amount: int
    ) -> None:
        """Delete X amount of messages based on the number you provide
        
        Args:
            interaction: The interaction that triggered this command
            amount: The amount of messages to delete
        """
        await interaction.response.defer(ephemeral=True)
        deleted = await interaction.channel.purge(limit=amount)
        await interaction.followup.send(
            f"Deleted {len(deleted)} messages.",
            ephemeral=True
        )
    
    @app_commands.command(name="delete_all_threads", description="Delete all threads in a specified channel")
    @is_mod_slash()
    async def delete_all_threads(self, interaction: discord.Interaction, channel: discord.TextChannel) -> None:
        """Delete all threads in a specified channel
        
        Args:
            interaction: The interaction that triggered this command
        """
        guild = interaction.guild
        if channel not in guild.text_channels:
            await interaction.response.send_message("Error: Could not find the specified channel.", ephemeral=True)
            return
        
        for thread in channel.threads:
            await thread.delete()
        
        await interaction.response.send_message("All threads have been deleted.", ephemeral=True)
    
    @app_commands.command(
        name="role",
        description="Assign or remove a role from yourself"
    )
    @app_commands.describe(
        action="Whether to assign or remove the role",
        role="The role to assign or remove"
    )
    @is_mod_slash()
    async def manage_role(
        self,
        interaction: discord.Interaction,
        action: Literal["assign", "remove"],
        role: discord.Role,
        member: discord.Member
    ) -> None:
        """Manage roles for the user who triggered the command
        
        Args:
            interaction: The interaction that triggered this command
            action: Whether to assign or remove the role
            role: The role to assign or remove
            
        Note:
            This command will:
            1. Check if the user already has/doesn't have the role
            2. Assign or remove the role if conditions are met
            3. Send appropriate feedback messages
        """
        # fetch the member to give the role to
        member = await interaction.guild.fetch_member(member.id)

        has_role = role in member.roles

        if action == "assign":
            if has_role:
                await interaction.response.send_message(
                    f"{member.mention} already has the **{role.name}** role.",
                    ephemeral=True
                )
                return
                
            try:
                await member.add_roles(role)
                await interaction.response.send_message(
                    f"Assigned **{role.name}** to {member.mention}.",
                    ephemeral=True
                )
                # Send a followup message that will be deleted after 5 seconds
                await interaction.followup.send(
                    f"{interaction.user.mention} assigned **{role.name}** role to {member.mention}.",
                    delete_after=5
                )
            except discord.Forbidden:
                await interaction.response.send_message(
                    "I don't have permission to assign roles!",
                    ephemeral=True
                )
        elif action == "remove":
            if not has_role:
                await interaction.response.send_message(
                    f"{member.mention} does not have the **{role.name}** role.",
                    ephemeral=True
                )
                return
                
            try:
                await member.remove_roles(role)
                await interaction.response.send_message(
                    f"Removed **{role.name}** from {member.mention}.",
                    ephemeral=True
                )
                # Send a followup message that will be deleted after 5 seconds
                await interaction.followup.send(
                    f"{interaction.user.mention} removed **{role.name}** role from {member.mention}.",
                    delete_after=5
                )
            except discord.Forbidden:
                await interaction.response.send_message(
                    "I don't have permission to remove roles!",
                    ephemeral=True
                )
        else:
            # This will not be reached but just wanted to show add and remove for commands
            await interaction.response.send_message(
                "Invalid action. Please use 'assign' or 'remove'.",
                ephemeral=True
            )


async def setup(bot: commands.Bot) -> None:
    """Add the General cog to the bot
    
    Args:
        bot: Discord bot instance
    """
    await bot.add_cog(General(bot))
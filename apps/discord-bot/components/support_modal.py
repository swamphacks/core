import discord
from discord.ui import Modal, TextInput
from discord import TextStyle
from components.buttons import ThreadOpenButton

class SupportModal(Modal, title="Support Inquiry"):
    def __init__(self):
        super().__init__()
        self.title_input = TextInput(label="Title", max_length=100)
        self.description_input = TextInput(label="Describe your issue", style=TextStyle.paragraph)
        self.add_item(self.title_input)
        self.add_item(self.description_input)


    async def on_submit(self, interaction: discord.Interaction):
        reports_channel = discord.utils.get(interaction.guild.channels, name="reports")
        
        if not reports_channel:
            await interaction.response.send_message(
                "Error: Could not find the reports channel. Please contact an administrator.",
                ephemeral=True
            )
            return

        # create embed for reports channel
        embed = discord.Embed(
            title=f"New Support Request: {self.title_input.value}",
            description=self.description_input.value,
            color=discord.Color.blue()
        )
        embed.set_author(name=f"From: {interaction.user}", icon_url=interaction.user.display_avatar.url)
        
                
        view = ThreadOpenButton()
        await reports_channel.send(embed=embed, view=view)
        await interaction.response.send_message("Thank you! Your issue has been submitted. A staff member may create a thread for you.", ephemeral=True)
        
        
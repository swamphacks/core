import discord
from discord.ui import View, Button
from discord import ButtonStyle, Interaction


class ThreadOpenButton(View):
    def __init__(self):
        super().__init__(timeout=None)

    @discord.ui.button(label="Open Support Thread", style=ButtonStyle.primary, custom_id="open_support_thread")
    async def open_thread(self, interaction: Interaction, button: Button):
        # replace this with actual thread-creation logic
        await interaction.response.send_message("Opening support thread...", ephemeral=True)


class LinkButtonView(View):
    def __init__(self, label: str, url: str):
        super().__init__()
        self.add_item(
            Button(
                label=label,
                style=ButtonStyle.link,
                url=url
            )
        )

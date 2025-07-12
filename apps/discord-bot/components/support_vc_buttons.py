from discord.ui import View, Button
from discord import ButtonStyle, Interaction
import discord
from discord.errors import NotFound
from utils.checks import is_mod_slash
from utils.mentor_functions import set_busy_mentor, set_available_mentor
from components.ticket_state import claimed_tickets

class SupportVCButtons(View):
    def __init__(self, voice_channel: discord.VoiceChannel, description_input: discord.ui.TextInput) -> None:
        super().__init__(timeout=None)
        self.voice_channel = voice_channel
        self.description_input = description_input
        self.add_item(ClaimTicketButton(voice_channel, description_input))
        self.add_item(CloseTicketButton(voice_channel, description_input))


class CloseTicketButton(Button):
    def __init__(self, voice_channel: discord.VoiceChannel, description_input: discord.ui.TextInput):
        super().__init__(label="Close Ticket", style=ButtonStyle.primary, custom_id="close_ticket", emoji="❌")
        self.voice_channel = voice_channel
        self.description_input = description_input

    async def callback(self, interaction: Interaction):
        claimed_tickets.pop(self.voice_channel.id, None)
        # print(claimed_tickets)
        try:
            await self.voice_channel.delete()
            await interaction.response.send_message(f"Voice channel: {self.voice_channel.mention} has been deleted.", ephemeral=True)
            await set_available_mentor(interaction.user, True)
            await set_busy_mentor(interaction.user, False)
            
            # edit original message to disable claim button
            message = interaction.message
            if not message:
                await interaction.response.send_message(
                    "Message not found.",
                    ephemeral=True
                )
                return
            new_view = SupportVCButtons(self.voice_channel, self.description_input)
            # disable all buttons in the view
            for item in new_view.children:
                item.disabled = True
            # copy the original embed and update its description
            embed = message.embeds[0] if message.embeds else None
            if embed:
                new_embed = embed.copy()
                new_embed.description = f"Issue: {self.description_input.value}\n\nActions: {interaction.user.mention} closed {self.voice_channel.name}."
                new_embed.color = discord.Color.red()
                await message.edit(embed=new_embed, view=new_view)
            else:
                await message.edit(view=new_view)
        except NotFound:
            await interaction.response.send_message(
                "This voice channel no longer exists.",
                ephemeral=True
            )
        except Exception as e:
            await interaction.response.send_message(f"Failed to delete the voice channel. Error: {e}", ephemeral=True)

class ClaimTicketButton(Button):
    def __init__(self, voice_channel: discord.VoiceChannel, description_input: discord.ui.TextInput):
        super().__init__(label="Claim Ticket", style=ButtonStyle.primary, custom_id="claim_ticket", emoji="📥")
        self.voice_channel = voice_channel
        self.description_input = description_input

    async def callback(self, interaction: Interaction):
        try:
            # Check if the ticket is already claimed
            if claimed_tickets.get(self.voice_channel.id):
                await interaction.response.send_message(
                    "This ticket has already been claimed by another mentor.",
                    ephemeral=True
                )
                return
            
            # check if mentor already has an active ticket
            if interaction.user.id in claimed_tickets.values():
                await interaction.response.send_message(
                    "You already have an active support thread or VC. Please close it before claiming a new one.",
                    ephemeral=True
                )
                return
            
            # Mark as claimed
            claimed_tickets[self.voice_channel.id] = interaction.user.id
            # print(claimed_tickets)
            
            
            await self.voice_channel.set_permissions(
                interaction.user,
                connect=True,
                view_channel=True
            )
            await interaction.response.send_message(
                f"Click here to join the voice channel: {self.voice_channel.mention}",
                ephemeral=True
            )
            await set_available_mentor(interaction.user, False)
            await set_busy_mentor(interaction.user, True)
            
            
            # Edit the original message to disable the claim button
            message = interaction.message
            new_view = SupportVCButtons(self.voice_channel, self.description_input)
            for item in new_view.children:
                if isinstance(item, ClaimTicketButton):
                    item.disabled = True

            # Copy the original embed and update its description
            embed = message.embeds[0] if message.embeds else None
            if embed:
                new_embed = embed.copy()
                new_embed.description = f"Issue: {self.description_input.value}\n\nActions: {interaction.user.mention} claimed the ticket. Please join the vc to assist the member."
                await message.edit(embed=new_embed, view=new_view)
            else:
                await message.edit(view=new_view)
            
        except NotFound:
            await interaction.response.send_message(
                "This voice channel no longer exists.",
                ephemeral=True
            )
        except Exception as e:
            await interaction.response.send_message(f"Failed to notify you about the voice channel. Error: {e}", ephemeral=True)
from discord.ui import View, Button, Select
import discord
from discord import Interaction
from components.support_modals import ThreadSupportModal, VCSupportModal

class TicketView(View):
    def __init__(self, mentors=None):
        super().__init__(timeout=None)
        self.add_item(TicketSelect())
        # self.add_item(MentorSelect(mentors))
        
        
class TicketSelect(Select):
    def __init__(self):
        options = [
            discord.SelectOption(label="Chat in threads", value="thread", emoji="💬", description="Open a support thread"),
            discord.SelectOption(label="Chat in VC", value="vc", emoji="🎤", description="Open a private voice chat"),
        ]
        super().__init__(
            placeholder="Choose your support type...",
            min_values=1,
            max_values=1,
            options=options
        )

    async def callback(self, interaction: Interaction):
        if self.values[0] == "thread":
            try:
                await interaction.response.send_modal(ThreadSupportModal())
                
                # reset selection when clicked
                if interaction.message:
                    from components.ticket_view import TicketView
                    await interaction.message.edit(view=TicketView())
            except Exception as e:
                await interaction.response.send_message(
                    "Sorry, there was an error opening the support modal. Please try again later.",
                    ephemeral=True
                )
                print(f"Error in open_threads: {e}")
        elif self.values[0] == "vc":
            try:
                await interaction.response.send_modal(VCSupportModal())
                
                # reset selection when clicked
                if interaction.message:
                    from components.ticket_view import TicketView
                    await interaction.message.edit(view=TicketView())
            except Exception as e:
                await interaction.response.send_message(
                    "Sorry, there was an error opening the support modal. Please try again later.",
                    ephemeral=True
                )
                print(f"Error in open_vc: {e}")

# This is for selecting mentor dropdown not sure if we will use it yet.
# class MentorSelect(Select):
#     """ Select a specific mentor from all mentors.
#     """
#     def __init__(self, mentors):
#         self.mentors = mentors
#         options = []
#         for mentor in mentors:
#             options.append(
#                 discord.SelectOption(
#                     label=mentor.name,
#                     value=str(mentor.id),
#                     emoji="👨‍🏫" if mentor else None
#                 ))
#         super().__init__(
#             placeholder="(Optional) Choose a mentor...",
#             min_values=1,
#             max_values=1,
#             options=options
#         )
    
#     async def callback(self, interaction: Interaction):
#         if not self.values:
#             await interaction.response.send_message(
#                 "Please select a mentor to proceed.",
#                 ephemeral=True
#             )
#             return
        
#         mentor_id = int(self.values[0])
#         mentor = discord.utils.get(self.mentors, id=mentor_id)
        
#         if mentor:
#             await interaction.response.send_message(
#                 f"You have selected {mentor.name} as your mentor.",
#                 ephemeral=True
#             )
#         else:
#             await interaction.response.send_message(
#                 "Selected mentor is not available. Please try again.",
#                 ephemeral=True
#             ) 
    

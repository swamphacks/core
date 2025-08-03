import discord

def get_available_mentors(mod_role: discord.Role) -> list[discord.Member]:
    available_mentors = []
    for mentor in mod_role.members:
        if "Available Mentor" in [role.name for role in mentor.roles]:
            available_mentors.append(mentor)
    return available_mentors
    
async def set_all_mentors_available(mod_role: discord.Role) -> None:
    available_role = discord.utils.get(mentor.guild.roles, name="Available Mentor")
    if not available_role:
        return
    for mentor in mod_role.members:
        # print(f"Checking mentor: {mentor.name}")
        if "Available Mentor" not in [role.name for role in mentor.roles]:
            
            await mentor.add_roles(available_role)
        else:
            print(f"{mentor.name} is already an available mentor")
            return


async def set_busy_mentor(mentor: discord.Member, busy: bool) -> None:
    busy_role = discord.utils.get(mentor.guild.roles, name="Busy Mentor")
    if not busy_role:
        return
    if busy == True:
        await mentor.add_roles(busy_role)
    else:
        if busy_role in mentor.roles:
            await mentor.remove_roles(busy_role)
        else:
            print(f"No busy mentor role to remove from {mentor.name}")
            return
        

async def set_available_mentor(mentor: discord.Member, available: bool) -> None:
    available_role = discord.utils.get(mentor.guild.roles, name="Available Mentor")
    if not available_role:
        print(f"No available mentor role found for {mentor.name}")
        return
    if available == True:
        await mentor.add_roles(available_role)
    else:
        
        if available_role in mentor.roles:
            await mentor.remove_roles(available_role)
        else:
            print(f"No available mentor role to remove from {mentor.name}")
            return
        
from discord.ext import commands

class General(commands.Cog):
    def __init__(self, bot):
        self.bot = bot
    
    @commands.command()
    async def pickle(self, ctx):
        await ctx.send("PICKLEBALL!")


async def setup(bot):
    await bot.add_cog(General(bot))
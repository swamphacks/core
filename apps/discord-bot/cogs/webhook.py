from discord.ext import commands
import discord
import json
from aiohttp import web
from typing import Optional
import os
import logging


class Webhooks(commands.Cog):
    """A cog that handles webhooks from the API to assign Discord roles"""
    
    def __init__(self, bot: commands.Bot) -> None:
        """Initialize the Webhooks cog
        
        Args:
            bot: Discord bot instance
        """
        self.bot: commands.Bot = bot
        self.webhook_port: int = int(os.getenv('WEBHOOK_PORT', '8081'))
        self.webhook_app: Optional[web.Application] = None
        self.webhook_runner: Optional[web.AppRunner] = None
        self.logger = logging.getLogger(__name__)
        
    async def setup_webhook_server(self) -> None:
        """Set up the webhook HTTP server"""
        self.webhook_app = web.Application()
        self.webhook_app.router.add_post('/webhook', self.handle_webhook)
        
        aiohttp_access_logger = logging.getLogger('aiohttp.access')
        aiohttp_access_logger.setLevel(logging.WARNING)
        
        self.webhook_runner = web.AppRunner(self.webhook_app)
        await self.webhook_runner.setup()
        
        site = web.TCPSite(self.webhook_runner, '0.0.0.0', self.webhook_port)
        await site.start()
        print(f"Webhook server started on port {self.webhook_port}")
    
    async def handle_webhook(self, request: web.Request) -> web.Response:
        """Handle incoming webhook requests
        
        Args:
            request: The incoming HTTP request
            
        Returns:
            HTTP response
        """
        try:
            data = await request.json()
            action = data.get('action')
            user_id = data.get('user_id')
            role_name = data.get('role_name')
            guild_id = data.get('guild_id')
            
            if action == 'assign_role' and user_id and role_name and guild_id:
                # get the guild by id
                guild = self.bot.get_guild(int(guild_id))
                if not guild:
                    return web.Response(text=f'Guild {guild_id} not found', status=404)
                
                # get the member by discord user id
                try:
                    member = await guild.fetch_member(int(user_id))
                except discord.NotFound:
                    return web.Response(text=f'User {user_id} not found in server', status=404)
                except discord.HTTPException as e:
                    return web.Response(text=f'Error fetching member: {str(e)}', status=500)
                
                # get the role by name
                role = discord.utils.get(guild.roles, name=role_name)
                if not role:
                    return web.Response(text=f'Role "{role_name}" not found', status=404)
                
                # check if user already has the role
                if role in member.roles:
                    return web.Response(text='User already has the role', status=200)
                
                # assign the role
                try:
                    await member.add_roles(role, reason='Role assignment via webhook')
                    return web.Response(text='Role assigned successfully', status=200)
                except discord.Forbidden:
                    return web.Response(text='Bot lacks permission to assign roles', status=403)
                except discord.HTTPException as e:
                    return web.Response(text=f'Error assigning role: {str(e)}', status=500)
            else:
                return web.Response(text='Invalid webhook payload', status=400)
                
        except json.JSONDecodeError:
            return web.Response(text='Invalid JSON', status=400)
        except Exception as e:
            print(f"Error handling webhook: {e}")
            return web.Response(text=f'Internal server error: {str(e)}', status=500)
    
    @commands.Cog.listener()
    async def on_ready(self) -> None:
        """Start webhook server when bot is ready"""
        await self.setup_webhook_server()


async def setup(bot: commands.Bot) -> None:
    """Setup function for the cog
    
    Args:
        bot: The bot instance
    """
    await bot.add_cog(Webhooks(bot))

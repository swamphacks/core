from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Optional
import asyncio
import os
from dotenv import load_dotenv
from pathlib import Path

# Attempt to load a repo-level infra env file first (core/infra/.env.dev.discord)
# This helps when the working directory is `core/apps/discord-bot` but env files live in `core/infra/`.
try:
    repo_root = Path(__file__).resolve().parents[2]
    infra_env = repo_root / "infra" / ".env.dev.discord"
    if infra_env.exists():
        load_dotenv(infra_env)
        print(f"Loaded env from {infra_env}")
    else:
        # fallback to default behavior (load .env from current working directory)
        load_dotenv()
except Exception:
    # fallback
    load_dotenv()

app = FastAPI()


class CheckResponse(BaseModel):
    in_guild: bool


# This module expects a `bot` attribute to be attached at runtime.
# main.py will attach the running discord bot instance as `api_bot`.
bot = None


def attach_bot(discord_bot):
    """Attach the running discord bot instance so endpoints can use it."""
    global bot
    bot = discord_bot


@app.get("/check_user_in_guild", response_model=CheckResponse)
async def check_user_in_guild(guild_id: Optional[int] = None, user_id: int = None):
    """Check whether a user with the given Discord user ID exists in the guild.

    Query params:
    - guild_id: integer guild id to search in (optional, falls back to GUILD_ID env)
    - user_id: Discord user id (snowflake) to check membership for (required)
    """
    # If guild_id not supplied, try env var GUILD_ID
    if guild_id is None:
        env_gid = os.getenv("GUILD_ID")
        # print(env_gid)
        if env_gid:
            try:
                guild_id = int(env_gid)
            except Exception:
                raise HTTPException(status_code=400, detail="Invalid GUILD_ID in environment")
        else:
            raise HTTPException(status_code=400, detail="guild_id query param or GUILD_ID env var is required")

    if bot is None:
        raise HTTPException(status_code=503, detail="Bot not attached")

    if not bot.is_ready():
        # wait a bit for ready state
        try:
            await asyncio.wait_for(bot.wait_until_ready(), timeout=5)
        except Exception:
            raise HTTPException(status_code=503, detail="Bot not ready")

    guild = bot.get_guild(guild_id)
    if guild is None:
        # try fetching from API as fallback
        try:
            guild = await bot.fetch_guild(guild_id)
        except Exception:
            raise HTTPException(status_code=404, detail="Guild not found")

    # FastAPI enforces user_id is provided (required query param)

    try:
        member = guild.get_member(user_id)
        if member is not None:
            return CheckResponse(in_guild=True)
        # try API fetch
        try:
            member = await guild.fetch_member(user_id)
            if member is not None:
                return CheckResponse(in_guild=True)
        except Exception:
            pass
        return CheckResponse(in_guild=False)
    except Exception:
        return CheckResponse(in_guild=False)

# Installation

The Discord bot is a Python service managed with [uv](https://docs.astral.sh/uv/).

## Prerequisites

- Python 3.12+
- uv: `pip install uv` or see the [uv docs](https://docs.astral.sh/uv/getting-started/installation/)
- A Discord bot token ([Discord Developer Portal](https://discord.com/developers/applications))
- A Google Gemini API key ([Google AI Studio](https://aistudio.google.com/))

## Environment

Copy the example file:

```bash
cp apps/discord-bot/.env.example apps/discord-bot/.env
```

| Variable | Description |
|---|---|
| `DISCORD_TOKEN` | Bot token from the Discord Developer Portal |
| `API_KEY` | API key for authenticating requests to the SwampHacks API |
| `GEMINI_API_KEY` | Google Gemini API key for AI chatbot features |
| `API_URL` | SwampHacks API base URL (`http://localhost:8080` for local dev) |
| `SESSION_COOKIE` | Session cookie value for API authentication |
| `WEBHOOK_URL` | Discord webhook URL for outbound notifications |
| `WEBHOOK_PORT` | Port for the internal webhook listener |
| `EVENT_ID` | The current hackathon event ID |

## Running

```bash
cd apps/discord-bot
uv sync
uv run main.py
```

The bot connects to Discord and starts an HTTP server (FastAPI) on `WEBHOOK_PORT` for incoming webhooks.

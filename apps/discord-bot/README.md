# Discord Bot

## Prerequisites

Before you begin, ensure you have the following installed:

- **uv** by Astral.sh  
  Install from the official page: [https://github.com/astral-sh/uv](https://github.com/astral-sh/uv)
- **Python 3.12 or higher**  
  Download from: [https://www.python.org/downloads/](https://www.python.org/downloads/)

- **Git**  
  Download from: [https://git-scm.com/downloads](https://git-scm.com/downloads)

---

## Setup & Run the Bot

### 1. Configure Environment Variables

Create a `.env` file in the `discord-bot` directory with your credentials (or follow `.env.example`):

```env
DISCORD_TOKEN=YOUR_DISCORD_BOT_TOKEN_HERE
GEMINI_API_KEY=YOUR_GEMINI_API_KEY_HERE
```

---

### 2. Create virtual environment (venv) and run the bot

Make sure you are inside the `core/apps/discord-bot` directory, then run:

```bash
# This commands creates a venv automatically and runs main.py
uv run main.py
```

You should see a confirmation message like:

```
Bot: HackrBot#3118 is ready to go.
```

---

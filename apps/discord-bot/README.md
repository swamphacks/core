# Discord Bot

## Prerequisites

Before you begin, make sure you have:

* **Python 3.8+** installed and on your PATH
* **Git** (optional if you clone the repo)

*All commands assume you are in the `core/apps/discord-bot/` directory.*

## 1. Create a Virtual Environment

```bash
python -m venv .venv
```

Make sure to select the path to the interpreter usually located in .venv/bin/pythonX

## 2. Activate the Environment

### macOS / Linux

```bash
source .venv/bin/activate
```

### Windows (PowerShell)

```powershell
.\.venv\Scripts\Activate.ps1
```

*(Your prompt should now be prefixed with `(.venv)`.)*

## 3. Install Dependencies

```bash
# Upgrade pip
python -m pip install --upgrade pip

# Install runtime dependencies
pip install -r requirements.txt
```

## 4. Configure Environment Variables

Create a file named `.env` in this directory with the following content:

```
DISCORD_TOKEN=YOUR_DISCORD_BOT_TOKEN_HERE
```

## 5. Run the Bot

With the venv active, start your bot:

```bash
python main.py
```

You should see something like:

```
Bot: SwampHackr is ready to go!
```

---

> If you ever switch to a fresh terminal session, re-activate the venv before running any commands:
>
> ```bash
> source .venv/bin/activate  # (macOS/Linux)
> # or
> .\.venv\Scripts\Activate.ps1  # (Windows)
> ```


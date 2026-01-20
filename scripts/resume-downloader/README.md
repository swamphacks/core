# Resume Downloader

Simple tool used for gathering all resumes from Cloudflare R2 and saving the file as the user's name.

## Setup

1. Clone the repository if you haven't
2. Install dependencies through uv (install uv first)
Run the following:
```
  uv sync
```

3. Create a .env file. Copy the environment variables from .env.example

4. Get a csv file with a structure matching users.csv.example

## Usage

### Download all resumes corresponding to csv

uv run main.py 


# Bulk Email CLI

This is a simple command-line tool to manage SES email templates and send bulk emails.

## Setup

1. Clone the repository and navigate to the project folder.
2. Install dependencies through uv (install uv first)
Run the following:
```
  uv sync
```

3. Create a .env file in the root directory with your AWS credentials:

AWS_ACCESS_KEY=YOUR_ACCESS_KEY
AWS_ACCESS_KEY_SECRET=YOUR_ACCESS_KEY_SECRET

## Templates

- All HTML templates should live in the templates/ folder.
- Templates can be uploaded or updated from this folder.

## Usage

### List all templates in SES

uv run main.py list

### Upload or update a template

uv run main.py upload TEMPLATE_NAME template_file.html --subject "Email Subject" --text "Plain text fallback"

> Note: Place all HTML files inside the `templates/` folder. Then simply specify the HTML file. The CLI tool parses for templates inside the templates folder.

### Send bulk emails

uv run main.py send TEMPLATE_NAME --csv contacts.csv

- contacts.csv should have an email column with recipient addresses.

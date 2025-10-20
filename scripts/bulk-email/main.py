import csv
import os
import boto3
import argparse
from typing import List
from dotenv import load_dotenv
from botocore.exceptions import ClientError
from math import ceil

load_dotenv()

AWS_ACCESS_KEY = os.getenv("AWS_ACCESS_KEY")
AWS_ACCESS_KEY_SECRET = os.getenv("AWS_ACCESS_KEY_SECRET")
REGION = "us-east-1"


class Contact:
    def __init__(self, email: str):
        self.email = email

    def get_email(self):
        return self.email


def parse_csv(csv_path="contacts.csv") -> List[Contact]:
    contacts: List[Contact] = []
    with open(csv_path, newline="") as f:
        reader = csv.DictReader(f)
        for row in reader:
            if "email" not in row or not row["email"]:
                raise ValueError(f"Missing email field in row: {row}")
            contacts.append(Contact(row["email"]))
    return contacts


def get_ses_client():
    return boto3.client(
        "ses",
        region_name=REGION,
        aws_access_key_id=AWS_ACCESS_KEY,
        aws_secret_access_key=AWS_ACCESS_KEY_SECRET
    )


def list_templates():
    ses = get_ses_client()
    response = ses.list_templates()
    templates = response.get("TemplatesMetadata", [])
    if not templates:
        print("No templates found.")
        return
    print("SES Templates:")
    for t in templates:
        print(f"- {t['Name']} (Created: {t['CreatedTimestamp']})")


def upload_template(template_name: str, html_file_name: str, subject: str = "No Subject", text: str = "Your email client does not support HTML"):
    # Construct path relative to the templates/ folder
    html_file_path = os.path.join("templates", html_file_name)

    if not os.path.exists(html_file_path):
        print(f"File not found: {html_file_path}")
        return

    with open(html_file_path, "r", encoding="utf-8") as f:
        html_content = f.read()

    ses = get_ses_client()

    try:
        ses.get_template(TemplateName=template_name)
        print(f"Template '{template_name}' exists. Updating it...")
        ses.update_template(
            Template={
                "TemplateName": template_name,
                "SubjectPart": subject,
                "TextPart": text,
                "HtmlPart": html_content
            }
        )
        print(f"Template '{template_name}' updated successfully.")
    except ClientError as e:
        if e.response['Error']['Code'] == "TemplateDoesNotExist":
            ses.create_template(
                Template={
                    "TemplateName": template_name,
                    "SubjectPart": subject,
                    "TextPart": text,
                    "HtmlPart": html_content
                }
            )
            print(f"Template '{template_name}' created successfully.")
        else:
            raise


def chunk_list(lst, chunk_size: int):
    """Split a list into chunks of size chunk_size."""
    for i in range(0, len(lst), chunk_size):
        yield lst[i:i + chunk_size]


def send_emails(template_name: str, csv_path="contacts.csv"):
    contacts = parse_csv(csv_path)
    if not contacts:
        print("No contacts found in CSV.")
        return

    ses = get_ses_client()

    # SES doesn't support sending more that 50 emails in one API call
    chunk_size = 50

    for chunk_index, chunk in enumerate(chunk_list(contacts, chunk_size), start=1):
        destinations = [{"Destination": {"ToAddresses": [contact.get_email()]}} for contact in chunk]

        try:
            ses.send_bulk_templated_email(
                Source="noreply@swamphacks.com",
                Template=template_name,
                Destinations=destinations,
                DefaultTemplateData="{}"
            )
            print(f"Chunk {chunk_index}: Emails sent successfully!")
        except ClientError as e:
            print(f"Chunk {chunk_index}: Failed to send emails: {e}")


def main():
    parser = argparse.ArgumentParser(description="SES Bulk Email CLI")
    subparsers = parser.add_subparsers(dest="command")

    # List templates
    subparsers.add_parser("list", help="List SES templates")

    # Upload template
    upload_parser = subparsers.add_parser("upload", help="Upload or update SES template from HTML file")
    upload_parser.add_argument("name", help="Template name")
    upload_parser.add_argument("html_path", help="Path to HTML file")
    upload_parser.add_argument("--subject", help="Email subject line", default="No Subject")
    upload_parser.add_argument("--text", help="Plain text fallback", default="Your email client does not support HTML")

    # Send emails
    send_parser = subparsers.add_parser("send", help="Send bulk emails using a template")
    send_parser.add_argument("template_name", help="SES template name to use")
    send_parser.add_argument("--csv", help="Path to contacts CSV", default="contacts.csv")

    args = parser.parse_args()

    if args.command == "list":
        list_templates()
    elif args.command == "upload":
        upload_template(args.name, args.html_path, args.subject, args.text)
    elif args.command == "send":
        send_emails(args.template_name, args.csv)
    else:
        parser.print_help()


if __name__ == "__main__":
    main()

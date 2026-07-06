#!/usr/bin/env python3
"""
Generate fake user data and insert it into the `users` table.

Usage:
    export DATABASE_URL="postgresql://user:password@host:port/dbname"
    python seed_fake_users.py --count 100
"""

import argparse
import os
import random
import sys
import uuid
from datetime import timedelta

import psycopg2
from psycopg2.extras import execute_values
from faker import Faker
from dotenv import load_dotenv

load_dotenv()

fake = Faker()

# Must match: create type user_role as enum ('admin', 'staff', 'attendee', 'applicant', 'visitor');
# All generated users get this role.
ROLE = "applicant"


def random_past_datetime(start_days_ago=365, end_days_ago=0):
    """Return a random aware datetime between `start_days_ago` and `end_days_ago` days in the past."""
    from datetime import datetime, timezone

    delta_days = random.uniform(end_days_ago, start_days_ago)
    return datetime.now(timezone.utc) - timedelta(days=delta_days)


def generate_user(existing_emails: set) -> dict:
    """Generate a single fake user row as a dict matching the users table columns."""
    name = fake.name()

    # Make a reasonably unique, unique-per-row email
    email = fake.unique.email()
    while email in existing_emails:
        email = fake.unique.email()
    existing_emails.add(email)

    email_verified = random.random() < 0.85  # most emails verified
    onboarded = random.random() < 0.7  # most users onboarded
    image = fake.image_url() if random.random() < 0.5 else None

    created_at = random_past_datetime(start_days_ago=400)
    # updated_at is always >= created_at
    updated_at = created_at + timedelta(days=random.uniform(0, 30))

    # preferred_email: sometimes set, sometimes same domain different alias, sometimes None
    preferred_email = (
        email
        if random.random() < 0.3
        else (fake.email() if random.random() < 0.3 else None)
    )

    email_consent = random.random() < 0.6

    # checked_in_at: only some users have checked in, and only after created_at
    checked_in_at = None
    if random.random() < 0.4:
        checked_in_at = created_at + timedelta(hours=random.uniform(1, 24 * 60))

    # rfid: only some users (e.g. staff/attendees who picked up a badge) have one
    rfid = (
        fake.bothify(text="RFID-########").upper() if random.random() < 0.35 else None
    )

    role = ROLE

    # role_assigned_at: since every user here is explicitly assigned 'applicant'
    # (not the default 'visitor'), always set this.
    role_assigned_at = created_at + timedelta(hours=random.uniform(0, 24 * 10))

    return {
        "id": str(uuid.uuid4()),
        "name": name,
        "email": email,
        "email_verified": email_verified,
        "onboarded": onboarded,
        "image": image,
        "created_at": created_at,
        "updated_at": updated_at,
        "preferred_email": preferred_email,
        "email_consent": email_consent,
        "checked_in_at": checked_in_at,
        "rfid": rfid,
        "role_assigned_at": role_assigned_at,
        "role": role,
        "is_fake": True
    }


COLUMNS = [
    "id",
    "name",
    "email",
    "email_verified",
    "onboarded",
    "image",
    "created_at",
    "updated_at",
    "preferred_email",
    "email_consent",
    "checked_in_at",
    "rfid",
    "role_assigned_at",
    "role",
    "is_fake"
]


def insert_users(conn, users: list[dict]):
    """Bulk insert a list of user dicts into the users table."""
    rows = [tuple(u[col] for col in COLUMNS) for u in users]

    # %s::user_role casts the role column explicitly so psycopg2 doesn't need
    # to know about the custom enum type ahead of time.
    template = (
        "(%(id)s, %(name)s, %(email)s, %(email_verified)s, %(onboarded)s, %(image)s, "
        "%(created_at)s, %(updated_at)s, %(preferred_email)s, %(email_consent)s, "
        "%(checked_in_at)s, %(rfid)s, %(role_assigned_at)s, %(role)s::user_role, %(is_fake)s)"
    )

    query = f"INSERT INTO users ({', '.join(COLUMNS)}) VALUES %s"

    with conn.cursor() as cur:
        execute_values(cur, query, users, template=template, page_size=500)
    conn.commit()


def main():
    parser = argparse.ArgumentParser(description="Seed the users table with fake data.")
    parser.add_argument(
        "--count", type=int, default=50, help="Number of fake users to generate."
    )
    parser.add_argument(
        "--database-url",
        default=os.getenv("DATABASE_URL", ""),
        help="Postgres connection string. Defaults to the DATABASE_URL env var.",
    )
    args = parser.parse_args()

    if not args.database_url:
        print(
            "Error: no database URL provided. Set DATABASE_URL or pass --database-url.",
            file=sys.stderr,
        )
        sys.exit(1)

    print(f"Generating {args.count} fake users...")
    existing_emails: set[str] = set()
    users = [generate_user(existing_emails) for _ in range(args.count)]

    print("Connecting to database...")
    conn = psycopg2.connect(args.database_url)
    try:
        insert_users(conn, users)
        print(f"Inserted {len(users)} fake users into `users`.")
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()


if __name__ == "__main__":
    main()

"""
Seed script for the `users` table.

Usage:
    export DATABASE_URL="postgresql://user:password@host:port/dbname"
    python seed_users.py

Requires:
    pip install psycopg2-binary faker
"""

import os
import random
import psycopg2
from psycopg2.extras import execute_values
from faker import Faker
from dotenv import load_dotenv

load_dotenv()
fake = Faker()

# NOTE: adjust to match your actual `role` enum values in Postgres
ROLES = ["visitor", "admin"]

NUM_ROWS = 50


def generate_row():
    name = fake.name()
    email_verified = random.choice([True, False])
    onboarded = random.choice([True, False])
    email_consent = random.choice([True, False])
    role = random.choice(ROLES)
    is_fake = True  # mark seeded rows as fake data

    return {
        "name": name,
        "email": fake.unique.email(),
        "email_verified": email_verified,
        "onboarded": onboarded,
        "image": fake.image_url() if random.random() < 0.5 else None,
        "preferred_email": fake.email() if random.random() < 0.5 else None,
        "email_consent": email_consent,
        "checked_in_at": fake.date_time_this_year(tzinfo=None) if random.random() < 0.4 else None,
        "rfid": fake.bothify(text="RFID-########") if random.random() < 0.3 else None,
        "role_assigned_at": fake.date_time_this_year(tzinfo=None) if random.random() < 0.5 else None,
        "role": role,
        "has_seen_new_application_status": random.choice([True, False, None]),
        "is_fake": is_fake,
    }


def main():
    database_url = os.getenv("DATABASE_URL")
    if not database_url:
        raise SystemExit("Set the DATABASE_URL environment variable before running this script.")

    rows = [generate_row() for _ in range(NUM_ROWS)]

    columns = [
        "name",
        "email",
        "email_verified",
        "onboarded",
        "image",
        "preferred_email",
        "email_consent",
        "checked_in_at",
        "rfid",
        "role_assigned_at",
        "role",
        "has_seen_new_application_status",
        "is_fake",
    ]

    values = [tuple(row[col] for col in columns) for row in rows]

    insert_sql = f"""
        INSERT INTO users ({", ".join(columns)})
        VALUES %s
    """

    conn = psycopg2.connect(database_url)
    try:
        with conn:
            with conn.cursor() as cur:
                execute_values(cur, insert_sql, values)
        print(f"Inserted {len(values)} rows into users.")
    finally:
        conn.close()


if __name__ == "__main__":
    main()
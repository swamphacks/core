"""
Seed script for the `hackathons` table.

Usage:
    export DATABASE_URL="postgresql://user:password@host:port/dbname"
    python seed_hackathons.py

Requires:
    pip install psycopg2-binary faker

NOTE: `id` is a `text` primary key with no default, so this script
generates readable slug-style ids (e.g. "hackfest-2024-07"). Run this
BEFORE seed_redeemables.py, since redeemables.hackathon_id has a foreign
key to hackathons.id.

Also note: the DB enforces a unique constraint (`only_one_hackathon_active`)
that allows at most one row in the whole table with is_active = true. Since
this script can't know whether an active hackathon already exists before it
runs, every seeded row is set to is_active = False. Flip one to True by hand
afterward if you need an "active" hackathon for testing.
"""

import os
import random
import psycopg2
from psycopg2.extras import execute_values
from faker import Faker
from datetime import timedelta
from dotenv import load_dotenv

load_dotenv()

fake = Faker()

NUM_ROWS = 50

NAME_PREFIXES = [
    "HackFest", "CodeStorm", "InnovateX", "ByteBuild", "HackWave",
    "TechSprint", "DevJam", "HackNorth", "CodeCraft", "BuildAthon",
]


def make_id(name, index):
    slug = name.lower().replace(" ", "-")
    return f"{slug}-{index}"


def generate_row(index):
    name = f"{random.choice(NAME_PREFIXES)} {fake.year()}"
    hackathon_id = make_id(name, index)

    application_open = fake.date_time_between(start_date="-6M", end_date="-2M")
    application_close = application_open + timedelta(days=random.randint(14, 45))
    rsvp_deadline = application_close + timedelta(days=random.randint(3, 10))
    decision_release = application_close + timedelta(days=random.randint(1, 7))
    start_time = rsvp_deadline + timedelta(days=random.randint(3, 14))
    end_time = start_time + timedelta(days=random.choice([1, 2, 3]))

    accept_early = random.choice([True, False])
    if accept_early:
        early_application_open = application_open - timedelta(days=random.randint(10, 30))
        early_application_close = application_open - timedelta(days=1)
    else:
        early_application_open = None
        early_application_close = None

    return {
        "id": hackathon_id,
        "name": name,
        "description": fake.paragraph(nb_sentences=4),
        "location": fake.city(),
        "location_url": fake.url() if random.random() < 0.6 else None,
        "max_attendees": random.choice([100, 150, 200, 300, 500]),
        "application_open": application_open,
        "application_close": application_close,
        "rsvp_deadline": rsvp_deadline,
        "decision_release": decision_release,
        "start_time": start_time,
        "end_time": end_time,
        "is_active": False,
        "banner": fake.image_url() if random.random() < 0.5 else None,
        "application_review_started": random.choice([True, False]),
        "accept_early_applications": accept_early,
        "early_application_open": early_application_open,
        "early_application_close": early_application_close,
    }


def main():
    database_url = os.environ.get("DATABASE_URL")
    if not database_url:
        raise SystemExit("Set the DATABASE_URL environment variable before running this script.")

    rows = [generate_row(i) for i in range(NUM_ROWS)]

    columns = [
        "id",
        "name",
        "description",
        "location",
        "location_url",
        "max_attendees",
        "application_open",
        "application_close",
        "rsvp_deadline",
        "decision_release",
        "start_time",
        "end_time",
        "is_active",
        "banner",
        "application_review_started",
        "accept_early_applications",
        "early_application_open",
        "early_application_close",
    ]

    values = [tuple(row[col] for col in columns) for row in rows]

    insert_sql = f"""
        INSERT INTO hackathons ({", ".join(columns)})
        VALUES %s
    """

    conn = psycopg2.connect(database_url)
    try:
        with conn:
            with conn.cursor() as cur:
                execute_values(cur, insert_sql, values)
        print(f"Inserted {len(values)} rows into hackathons.")
    finally:
        conn.close()


if __name__ == "__main__":
    main()
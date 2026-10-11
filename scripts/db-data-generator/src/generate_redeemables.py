"""
Seed script for the `redeemables` table.

Usage:
    export DATABASE_URL="postgresql://user:password@host:port/dbname"
    python seed_redeemables.py

Requires:
    pip install psycopg2-binary faker

NOTE: `hackathon_id` has a foreign key to `hackathons.id`. This script
fetches existing hackathon ids from the database rather than making them
up, so run seed_hackathons.py FIRST.
"""

import os
import random
import psycopg2
from psycopg2.extras import execute_values
from faker import Faker
from dotenv import load_dotenv

load_dotenv()
fake = Faker()

NUM_ROWS = 50

REDEEMABLE_TYPES = ["meal", "tshirt"]

REDEEMABLE_NAMES = {
    "tshirt": ["T-Shirt", "Hoodie"],
    "meal": [
        "Energy Drink", "Coffee Voucher", "Pizza Slice", "Snack Pack",
        "Water Bottle",
    ],
}


def fetch_hackathon_ids(cur):
    cur.execute("SELECT id FROM hackathons")
    ids = [row[0] for row in cur.fetchall()]
    if not ids:
        raise SystemExit(
            "No rows found in `hackathons`. Run seed_hackathons.py first, "
            "or make sure the table is populated."
        )
    return ids


def generate_row(hackathon_ids):
    amount = random.randint(10, 500)
    max_user_amount = random.randint(1, min(5, amount))
    redeemable_type = random.choice(REDEEMABLE_TYPES)
    name = random.choice(REDEEMABLE_NAMES[redeemable_type])

    return {
        "name": name,
        "amount": amount,
        "max_user_amount": max_user_amount,
        "type": redeemable_type,
        "hackathon_id": random.choice(hackathon_ids),
    }

def main():
    database_url = os.getenv("DATABASE_URL")
    if not database_url:
        raise SystemExit("Set the DATABASE_URL environment variable before running this script.")

    columns = [
        "name",
        "amount",
        "max_user_amount",
        "type",
        "hackathon_id",
    ]

    insert_sql = f"""
        INSERT INTO redeemables ({", ".join(columns)})
        VALUES %s
    """

    conn = psycopg2.connect(database_url)
    try:
        with conn.cursor() as cur:
            hackathon_ids = fetch_hackathon_ids(cur)

        rows = [generate_row(hackathon_ids) for _ in range(NUM_ROWS)]
        values = [tuple(row[col] for col in columns) for row in rows]

        with conn:
            with conn.cursor() as cur:
                execute_values(cur, insert_sql, values)
        print(f"Inserted {len(values)} rows into redeemables.")
    finally:
        conn.close()


if __name__ == "__main__":
    main()
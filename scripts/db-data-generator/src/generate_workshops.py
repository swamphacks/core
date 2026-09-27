"""
Seed script for the `workshops` table.

Usage:
    export DATABASE_URL="postgresql://user:password@host:port/dbname"
    python seed_workshops.py

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

NUM_ROWS = 50

WORKSHOP_TYPES = ["workshop", "social"]
WORKSHOP_TOPICS = [
    "Intro to Machine Learning",
    "Building REST APIs with Flask",
    "React for Beginners",
    "Intro to Cybersecurity",
    "Data Visualization with Python",
    "Getting Started with Docker",
    "Intro to Blockchain",
    "UX Design Fundamentals",
    "Git & GitHub Basics",
    "Intro to Cloud Computing (AWS)",
    "Building Chatbots with LLMs",
    "SQL for Data Analysis",
    "Intro to Computer Vision",
    "Hardware Hacking 101",
    "Public Speaking for Engineers",
    "Intro to Game Development",
    "Web Scraping with Python",
    "Intro to Kubernetes",
    "Mobile App Development with Flutter",
    "Intro to Cryptography",
]

LOCATIONS = [
    "Room 101", "Room 202", "Main Auditorium", "Workshop Hall A",
    "Workshop Hall B", "Innovation Lab", "Conference Room 3", "Atrium Stage",
]


def generate_row():
    start_time = fake.date_time_between(start_date="-1d", end_date="+3d")
    duration_minutes = random.choice([30, 45, 60, 90, 120])
    end_time = start_time + __import__("datetime").timedelta(minutes=duration_minutes)

    return {
        "title": random.choice(WORKSHOP_TOPICS),
        "description": fake.paragraph(nb_sentences=3),
        "start_time": start_time,
        "end_time": end_time,
        "num_attendees": random.randint(0, 150),
        "location": random.choice(LOCATIONS),
        "presenter": fake.name(),
        "type": random.choice(WORKSHOP_TYPES),
    }


def main():
    database_url = os.getenv("DATABASE_URL")
    if not database_url:
        raise SystemExit("Set the DATABASE_URL environment variable before running this script.")

    rows = [generate_row() for _ in range(NUM_ROWS)]

    columns = [
        "title",
        "description",
        "start_time",
        "end_time",
        "num_attendees",
        "location",
        "presenter",
        "type",
    ]

    values = [tuple(row[col] for col in columns) for row in rows]

    insert_sql = f"""
        INSERT INTO workshops ({", ".join(columns)})
        VALUES %s
    """

    conn = psycopg2.connect(database_url)
    try:
        with conn:
            with conn.cursor() as cur:
                execute_values(cur, insert_sql, values)
        print(f"Inserted {len(values)} rows into workshops.")
    finally:
        conn.close()


if __name__ == "__main__":
    main()
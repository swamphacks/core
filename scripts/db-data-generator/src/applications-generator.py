#!/usr/bin/env python3
"""
Generate fake application data and insert it into the `applications` table,
one row per existing `applicant` user that doesn't already have an application.

Usage:
    export DATABASE_URL="postgresql://user:password@host:port/dbname"
    python seed_fake_applications.py --hackathon-id xii
"""

import argparse
import copy
import json
import os
import random
import sys
from datetime import timedelta, datetime, timezone
from dotenv import load_dotenv

load_dotenv()

import psycopg2
from psycopg2.extras import execute_values, Json

# The fixed application payload used for every generated application.
FAKE_APPLICATION = {
    "age": ">21",
    "diet": "celiac-disease",
    "race": "asian-pacific-islander",
    "year": "second_year",
    "level": "undergraduate",
    "phone": "1234567891",
    "essay1": "Example Response",
    "essay2": "Example Response",
    "essay3": "Example Response",
    "gender": "man",
    # "github": "https://www.github.com/my-profile",
    "majors": "Computer Science",
    "minors": "Entrepreneurship",
    "school": "University of Florida",
    "country": "US",
    "lastName": "Doe",
    "linkedin": "https://www.linkedin.com/in/my-profile",
    "pronouns": "he/him",
    "referral": "website,class_shoutout",
    "firstName": "John",
    "shirtSize": "L",
    "experience": "one",
    "race-other": "",
    "year-other": "",
    "level-other": "",
    "orientation": "heterosexual",
    "gender-other": "",
    "agreeToConduct": "agree",
    "graduationYear": "2026",
    "pictureConsent": "agree",
    "preferredEmail": "testEmail@gmail.com",
    "ufHackathonExp": "no",
    "universityEmail": "testEmail@school.edu",
    "ageCertification": True,
    "agreeToMLHEmails": "",
    "projectExperience": "independent_project",
    "infoShareAuthorization": "agree",
    "inpersonAcknowledgement": "agree",
}

STATUS = "submitted"


def fetch_applicant_user_ids(conn, hackathon_id: str, role: str, only_missing: bool):
    """
    Return a list of user ids with the given role who don't already have an
    application for this hackathon_id (if only_missing is True), or all
    matching users (if only_missing is False).
    """
    with conn.cursor() as cur:
        if only_missing:
            cur.execute(
                """
                SELECT u.id
                FROM users u
                LEFT JOIN applications a
                    ON a.user_id = u.id AND a.hackathon_id = %s
                WHERE u.role = %s::user_role
                  AND a.user_id IS NULL
                """,
                (hackathon_id, role),
            )
        else:
            cur.execute(
                "SELECT id FROM users WHERE role = %s::user_role",
                (role,),
            )
        return [row[0] for row in cur.fetchall()]


def generate_application(user_id: str, hackathon_id: str) -> dict:
    """Generate a single fake application row as a dict matching the applications table columns."""
    created_at = datetime.now(timezone.utc) - timedelta(days=random.uniform(1, 60))
    saved_at = created_at + timedelta(minutes=random.uniform(5, 240))
    submitted_at = saved_at + timedelta(minutes=random.uniform(1, 120))
    updated_at = submitted_at

    application_payload = copy.deepcopy(FAKE_APPLICATION)

    return {
        "user_id": user_id,
        "status": STATUS,
        "application": Json(application_payload),
        "created_at": created_at,
        "saved_at": saved_at,
        "updated_at": updated_at,
        "submitted_at": submitted_at,
        "hackathon_id": hackathon_id,
        "is_fake": True
    }


COLUMNS = [
    "user_id",
    "status",
    "application",
    "created_at",
    "saved_at",
    "updated_at",
    "submitted_at",
    "hackathon_id",
    "is_fake"
]


def insert_applications(conn, applications: list[dict]):
    """Bulk insert a list of application dicts into the applications table."""
    # %(status)s::application_status casts explicitly so psycopg2 doesn't need
    # to know about the custom enum type ahead of time.
    template = (
        "(%(user_id)s, %(status)s::application_status, %(application)s, %(created_at)s, "
        "%(saved_at)s, %(updated_at)s, %(submitted_at)s, %(hackathon_id)s, %(is_fake)s)"
    )

    query = f"INSERT INTO applications ({', '.join(COLUMNS)}) VALUES %s"

    with conn.cursor() as cur:
        execute_values(cur, query, applications, template=template, page_size=500)
    conn.commit()


def main():
    parser = argparse.ArgumentParser(
        description="Seed the applications table with fake data."
    )
    parser.add_argument(
        "--hackathon-id",
        default="xii",
        help="hackathon_id to attach to each generated application (default: xii).",
    )
    parser.add_argument(
        "--role",
        default="applicant",
        help="Only generate applications for users with this role (default: applicant).",
    )
    parser.add_argument(
        "--all",
        action="store_true",
        help="Generate applications for all matching users, even ones that already have an "
        "application for this hackathon_id (will fail on the primary key unless you "
        "expect duplicates to error out). By default, only users without an existing "
        "application for this hackathon_id are used.",
    )
    parser.add_argument(
        "--database-url",
        default=os.environ.get("DATABASE_URL"),
        help="Postgres connection string. Defaults to the DATABASE_URL env var.",
    )
    args = parser.parse_args()

    if not args.database_url:
        print(
            "Error: no database URL provided. Set DATABASE_URL or pass --database-url.",
            file=sys.stderr,
        )
        sys.exit(1)

    print("Connecting to database...")
    conn = psycopg2.connect(args.database_url)
    try:
        print(f"Fetching users with role '{args.role}'...")
        user_ids = fetch_applicant_user_ids(
            conn, args.hackathon_id, args.role, only_missing=not args.all
        )

        if not user_ids:
            print(
                f"No eligible users found (role='{args.role}', hackathon_id='{args.hackathon_id}'). "
                "Either there are no such users, or they all already have an application. "
                "Nothing to insert."
            )
            return

        print(
            f"Generating {len(user_ids)} fake applications for hackathon_id='{args.hackathon_id}'..."
        )
        applications = [
            generate_application(uid, args.hackathon_id) for uid in user_ids
        ]

        insert_applications(conn, applications)
        print(f"Inserted {len(applications)} fake applications into `applications`.")
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()


if __name__ == "__main__":
    main()

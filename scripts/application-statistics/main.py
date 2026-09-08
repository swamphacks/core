"""
Generate aggregate statistics from a SwampHacks applications CSV export.

The current SwampHacks export stores each applicant's form answers as JSON
inside the `application` column.

Essay "insights" are heuristic keyword/theme matches. They are useful for
spotting broad trends. An applicant can match more than one theme.

Usage:
    python main.py applications.csv

Optional JSON output:
    python main.py applications.csv --json application_summary.json
"""

from __future__ import annotations

import argparse
import csv
import json
import re
import statistics
from collections import Counter
from pathlib import Path


UF_SCHOOLS = {
    "University of Florida",
    "University of Florida Online",
}

GRADUATION_YEAR_MIN = 2023
GRADUATION_YEAR_MAX = 2050


LABELS = {
    "age": {
        "<18": "Under 18",
        "18-22": "18-22",
        ">22": "Over 22",
    },
    "gender": {
        "man": "Man",
        "woman": "Woman",
        "non-binary": "Non-binary",
        "self-describe": "Prefer to self-describe",
        "no-answer": "Prefer not to answer",
    },
    "pronouns": {
        "she/her": "She/her",
        "he/him": "He/him",
        "they/them": "They/them",
        "she/they": "She/they",
        "he/they": "He/they",
        "not-represented": "Not represented here",
        "no-answer": "Prefer not to answer",
    },
    "race": {
        "native-american-alaska-native": "Native American or Alaska Native",
        "asian-pacific-islander": "Asian or Pacific Islander",
        "black-african-american": "Black or African American",
        "hispanic-latino": "Hispanic or Latino",
        "white": "White",
        "middle-eastern": "Middle Eastern",
        "multiracial": "Mixed race",
        "other": "Other",
        "no-answer": "Prefer not to answer",
    },
    "orientation": {
        "heterosexual": "Heterosexual / straight",
        "homosexual": "Gay / lesbian",
        "bisexual": "Bisexual",
        "not-represented": "Other / not represented here",
        "no-answer": "Prefer not to answer",
    },
    "level": {
        "undergrad_two_year": "2-year / community college or similar",
        "undergrad_three_plus_year": "Undergraduate university (3+ year)",
        "graduate": "Graduate university",
        "postdoc": "Post doctorate",
        "bootcamp": "Code school / bootcamp",
        "vocational": "Vocational / trade / apprenticeship",
        "not_student": "Not currently a student",
        "prefer_no_answer": "Prefer not to answer",
        "other": "Other",
    },
    "year": {
        "first_year": "1st year",
        "second_year": "2nd year",
        "third_year": "3rd year",
        "fourth_year": "4th year",
        "graduate": "Graduate student",
        "other": "Other",
    },
    "experience": {
        "first_time": "SwampHacks would be my first hackathon",
        "one": "1 previous hackathon",
        "two": "2 previous hackathons",
        "three": "3 previous hackathons",
        "four_or_more": "4+ previous hackathons",
    },
    "ufHackathonExp": {
        "yes": "Yes",
        "no": "No",
    },
    "projectExperience": {
        "no_experience": "No - first programming experience",
        "course_experience": "No independent project, but course experience",
        "independent_project": "Yes",
    },
    "shirtSize": {
        "S": "Small",
        "M": "Medium",
        "L": "Large",
        "XL": "X-Large",
        "XXL": "XX-Large",
    },
    "diet": {
        "vegetarian": "Vegetarian",
        "vegan": "Vegan",
        "celiac-disease": "Celiac disease",
        "allergies": "Allergies",
        "kosher": "Kosher",
        "halal": "Halal",
        "other": "Other",
    },
    "referral": {
        "instagram": "Instagram",
        "discord": "Discord",
        "linkedin": "LinkedIn",
        "word_of_mouth": "Word of mouth",
        "website": "Website",
        "class_shoutout": "Class shoutout",
        "other_fl_hackathon": "Other Florida hackathon",
        "other": "Other",
    },
}


# Essay theme definitions
#
# Each category is counted at most once per applicant even if the applicant
# uses several matching words. Categories may overlap.

ESSAY1_MOTIVATIONS = {
    "Build / create something": [
        r"\bbuild(?:ing)?\b",
        r"\bcreat(?:e|ing)\b",
        r"\bdevelop(?:ing)? (?:a|an|something|project|app|application|website|tool|solution)\b",
        r"\bturn (?:an|a|my|our) idea\b",
    ],
    "Learn / improve technical skills": [
        r"\blearn(?:ing|ed)?\b",
        r"\bimprov(?:e|ing|ement)\b",
        r"\bexpand (?:my|our) (?:skills|knowledge)\b",
        r"\bnew (?:skills|technologies|tools|frameworks)\b",
    ],
    "Collaborate / work with others": [
        r"\bcollaborat\w*\b",
        r"\bteamwork\b",
        r"\bwork(?:ing)? with (?:a|my|other)? ?team\b",
        r"\bwork(?:ing)? (?:with|alongside) (?:others|people|students|peers)\b",
    ],
    "Meet people / network": [
        r"\bmeet (?:new )?(?:people|students|developers|engineers|peers|friends)\b",
        r"\bnetwork(?:ing)?\b",
        r"\bconnect(?:ing)? with (?:new )?(?:people|students|peers|others)\b",
        r"\bmake (?:new )?friends\b",
    ],
    "Solve meaningful / real-world problems": [
        r"\breal[- ]world\b",
        r"\bsolve (?:real|meaningful|actual)\b",
        r"\bmake (?:a|an)? ?(?:positive )?impact\b",
        r"\bsocial impact\b",
        r"\bmeaningful (?:project|solution|impact)\b",
    ],
    "Challenge themselves / try something new": [
        r"\bchallenge (?:myself|me|ourselves)\b",
        r"\boutside (?:of )?my comfort zone\b",
        r"\bpush (?:myself|my skills|ourselves)\b",
        r"\btry something new\b",
        r"\bstep outside\b",
    ],
    "Contribute / help other hackers": [
        r"\bcontribut\w*\b",
        r"\bhelp (?:other )?(?:hackers|participants|students|teammates)\b",
        r"\bshare (?:my|our) (?:knowledge|skills|experience|ideas)\b",
    ],
}


TECHNICAL_AREAS = {
    "AI / Machine Learning": [
        r"\bartificial intelligence\b",
        r"\bmachine learning\b",
        r"\bdeep learning\b",
        r"\bgenerative ai\b",
        r"\bcomputer vision\b",
        r"\bnatural language processing\b",
        r"\bnlp\b",
        r"\bllm[s]?\b",
        r"\bai\b",
        r"\bml\b",
    ],
    "Web / Full-stack Development": [
        r"\bweb development\b",
        r"\bweb app(?:lication)?s?\b",
        r"\bwebsite[s]?\b",
        r"\bfull[- ]?stack\b",
        r"\bfront[- ]?end\b",
        r"\bback[- ]?end\b",
        r"\breact(?:\.js)?\b",
        r"\bnext\.?js\b",
        r"\bnode\.?js\b",
    ],
    "Mobile App Development": [
        r"\bmobile app(?:lication)?s?\b",
        r"\bmobile development\b",
        r"\bios app(?:lication)?s?\b",
        r"\bandroid app(?:lication)?s?\b",
        r"\bflutter\b",
        r"\breact native\b",
        r"\bswift\b",
        r"\bkotlin\b",
    ],
    "Hardware / IoT / Embedded": [
        r"\bhardware\b",
        r"\biot\b",
        r"\binternet of things\b",
        r"\bembedded\b",
        r"\barduino\b",
        r"\braspberry pi\b",
        r"\bmicrocontroller[s]?\b",
        r"\b3d print\w*\b",
    ],
    "Data / Data Science": [
        r"\bdata science\b",
        r"\bdata analy(?:sis|tics)\b",
        r"\bdata visualization\b",
        r"\bdatabase[s]?\b",
        r"\bsql\b",
    ],
    "Cybersecurity": [
        r"\bcybersecurity\b",
        r"\bcyber security\b",
        r"\bctf\b",
        r"\bmalware\b",
        r"\bphishing\b",
        r"\bthreat detection\b",
    ],
    "Game Development": [
        r"\bgame development\b",
        r"\bvideo game[s]?\b",
        r"\bunity\b",
        r"\bunreal engine\b",
    ],
    "Cloud / DevOps": [
        r"\bcloud computing\b",
        r"\baws\b",
        r"\bazure\b",
        r"\bgcp\b",
        r"\bdocker\b",
        r"\bkubernetes\b",
        r"\bdevops\b",
    ],
    "Robotics": [
        r"\brobotics\b",
        r"\brobotic\b",
        r"\brobot[s]?\b",
    ],
    "AR / VR": [
        r"\baugmented reality\b",
        r"\bvirtual reality\b",
        r"\bar/vr\b",
    ],
}


ESSAY2_PROJECT_DOMAINS = {
    "Entertainment / games / media": [
        r"\bgame development\b",
        r"\bvideo game[s]?\b",
        r"\bgame\b",
        r"\bmusic\b",
        r"\bmovie[s]?\b",
        r"\bmedia\b",
    ],
    "Finance / budgeting": [
        r"\bfinanc(?:e|ial)\b",
        r"\bbudget(?:ing)?\b",
        r"\bexpense[s]?\b",
        r"\bbank(?:ing)?\b",
        r"\binvest(?:ing|ment)\b",
        r"\bpayment[s]?\b",
        r"\bmoney management\b",
    ],
    "Education / learning": [
        r"\beducat(?:ion|ional)\b",
        r"\btutor(?:ing)?\b",
        r"\bteacher[s]?\b",
        r"\bclassroom\b",
        r"\bstudy (?:tool|app|platform)\b",
        r"\blearning (?:app|platform|tool)\b",
        r"\blanguage[- ]learning\b",
    ],
    "Health / healthcare": [
        r"\bhealthcare\b",
        r"\bmedical\b",
        r"\bmedicine\b",
        r"\bpatient[s]?\b",
        r"\bmedication\b",
        r"\bmental health\b",
        r"\bwellness\b",
        r"\bhealth app\b",
    ],
    "Accessibility / assistive tech": [
        r"\baccessib(?:ility|le)\b",
        r"\bdisabilit(?:y|ies)\b",
        r"\bassistive\b",
        r"\bvisually impaired\b",
        r"\bblind\b",
        r"\bdeaf\b",
    ],
    "Transportation / mobility": [
        r"\btransportation\b",
        r"\btransit\b",
        r"\btraffic\b",
        r"\bnavigation\b",
        r"\broute optimization\b",
        r"\bmobility\b",
    ],
    "Environment / sustainability": [
        r"\benvironmental\b",
        r"\bsustainab(?:ility|le)\b",
        r"\bclimate\b",
        r"\brecycl(?:e|ing)\b",
        r"\bcarbon\b",
        r"\brenewable energy\b",
    ],
    "Cybersecurity / safety": [
        r"\bcybersecurity\b",
        r"\bcyber security\b",
        r"\bthreat detection\b",
        r"\bsecurity event\b",
        r"\bphishing\b",
        r"\bmalware\b",
    ],
}


ESSAY2_COLLABORATION = {
    "Communication / coordination": [
        r"\bcommunicat\w*\b",
        r"\bcoordinat\w*\b",
        r"\bcheck[- ]?in[s]?\b",
        r"\bmeeting[s]?\b",
        r"\bdiscuss(?:ed|ing|ion)?\b",
    ],
    "Divided tasks / assigned roles": [
        r"\bsplit\b",
        r"\bdivid(?:e|ed|ing)\b",
        r"\bassign(?:ed|ing)? (?:tasks|roles|responsibilities)\b",
        r"\bresponsib(?:le|ility|ilities)\b",
        r"\bdelegat\w*\b",
    ],
    "Feedback / iteration / testing": [
        r"\bfeedback\b",
        r"\biterat\w*\b",
        r"\btest(?:ed|ing)?\b",
        r"\brefin\w*\b",
    ],
    "Frontend / backend split": [
        r"\bfront[- ]?end\b",
        r"\bback[- ]?end\b",
    ],
    "Leadership / project management": [
        r"\bled\b",
        r"\blead(?:er|ership|ing)?\b",
        r"\bproject manag\w*\b",
        r"\bmanaged\b",
        r"\bdelegat\w*\b",
    ],
    "Git / GitHub / version control": [
        r"\bgit\b",
        r"\bgithub\b",
        r"\bversion control\b",
        r"\bpull request[s]?\b",
        r"\bbranch(?:es|ing)\b",
    ],
}


TECHNOLOGY_MENTIONS = {
    "Python": [r"\bpython\b"],
    "React": [r"\breact(?:\.js)?\b"],
    "JavaScript": [r"\bjavascript\b"],
    "TypeScript": [r"\btypescript\b"],
    "Java": [r"\bjava\b"],
    "C / C++": [r"\bc\+\+\b", r"\bcpp\b"],
    "Go": [r"\bgolang\b", r"\bwritten in go\b", r"\bbuilt in go\b", r"\busing go\b"],
    "Swift": [r"\bswift\b"],
    "Kotlin": [r"\bkotlin\b"],
    "Flutter": [r"\bflutter\b"],
    "Node.js": [r"\bnode(?:\.js|js)\b"],
    "FastAPI": [r"\bfastapi\b"],
    "Flask": [r"\bflask\b"],
    "Django": [r"\bdjango\b"],
    "Angular": [r"\bangular\b"],
    "Firebase": [r"\bfirebase\b"],
    "MongoDB": [r"\bmongodb\b", r"\bmongo db\b"],
    "SQL": [r"\bsql\b", r"\bpostgres(?:ql)?\b", r"\bmysql\b", r"\bsqlite\b"],
    "AWS": [r"\baws\b", r"\bamazon web services\b"],
    "Docker": [r"\bdocker\b"],
    "Git / GitHub": [r"\bgit\b", r"\bgithub\b"],
    "TensorFlow": [r"\btensorflow\b"],
    "PyTorch": [r"\bpytorch\b"],
    "OpenAI / GPT": [r"\bopenai\b", r"\bgpt(?:-\d+(?:\.\d+)?)?\b"],
    "Gemini": [r"\bgemini\b"],
    "Arduino": [r"\barduino\b"],
    "Raspberry Pi": [r"\braspberry pi\b"],
    "Unity": [r"\bunity\b"],
    "Figma": [r"\bfigma\b"],
}


ESSAY3_INTERESTS = {
    "Sports / fitness / outdoors": [
        r"\bsports?\b",
        r"\brunn(?:ing|er)?\b",
        r"\bgym\b",
        r"\bwork(?:ing)? out\b",
        r"\bhik(?:e|ing)\b",
        r"\btrek(?:king)?\b",
        r"\bsoccer\b",
        r"\bfootball\b",
        r"\bbasketball\b",
        r"\btennis\b",
        r"\bvolleyball\b",
        r"\bswimm(?:ing)?\b",
        r"\bcycling\b",
        r"\bbiking\b",
        r"\boutdoors?\b",
    ],
    "Gaming": [
        r"\bvideo game[s]?\b",
        r"\bgam(?:e|es|ing|er)\b",
    ],
    "Music": [
        r"\bmusic\b",
        r"\bsing(?:ing)?\b",
        r"\bguitar\b",
        r"\bpiano\b",
        r"\bdrums?\b",
        r"\bconcert[s]?\b",
    ],
    "Art / design / photography": [
        r"\bart\b",
        r"\bdrawing\b",
        r"\bpaint(?:ing)?\b",
        r"\bdesign\b",
        r"\bphotograph(?:y|er|ing)\b",
    ],
    "Reading / writing": [
        r"\bread(?:ing)?\b",
        r"\bbook[s]?\b",
        r"\bwriting\b",
        r"\bwriter\b",
    ],
    "Cooking / food": [
        r"\bcook(?:ing)?\b",
        r"\bbak(?:e|ing)\b",
        r"\bfood\b",
        r"\brestaurant[s]?\b",
    ],
    "Movies / TV": [
        r"\bmovie[s]?\b",
        r"\bfilm[s]?\b",
        r"\btv\b",
        r"\btelevision\b",
        r"\bshows?\b",
    ],
    "Travel": [
        r"\btravel(?:ing|led)?\b",
        r"\btrip[s]?\b",
    ],
    "Pets / animals": [
        r"\bpet[s]?\b",
        r"\bcat[s]?\b",
        r"\bdog[s]?\b",
        r"\banimal[s]?\b",
    ],
    "Anime / manga": [
        r"\banime\b",
        r"\bmanga\b",
    ],
}


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Generate aggregate SwampHacks application statistics."
    )
    parser.add_argument("csv_file", help="Path to the application CSV export.")
    parser.add_argument(
        "--json",
        dest="json_file",
        help="Optionally save the aggregate results to a JSON file.",
    )
    parser.add_argument(
        "--top",
        type=int,
        default=10,
        help="Number of schools, majors, minors, and technologies to show (default: 10).",
    )
    return parser.parse_args()


def clean(value) -> str:
    if value is None:
        return ""
    return str(value).strip()


def load_applications(csv_path: str) -> list[dict]:
    applications = []

    with open(csv_path, newline="", encoding="utf-8-sig") as file:
        reader = csv.DictReader(file)

        if "application" not in (reader.fieldnames or []):
            raise ValueError("CSV must contain an 'application' column.")

        for row_number, row in enumerate(reader, start=2):
            if clean(row.get("is_fake")).lower() == "true":
                continue

            if "submitted_at" in row and not clean(row.get("submitted_at")):
                continue

            try:
                application = json.loads(row["application"])
            except (json.JSONDecodeError, TypeError) as exc:
                print(
                    f"Warning: skipping row {row_number}; "
                    f"invalid application JSON ({exc})"
                )
                continue

            applications.append(application)

    return applications


def percent(count: int, total: int) -> float:
    return round((count / total) * 100, 1) if total else 0.0


def count_field(applications: list[dict], field: str) -> Counter:
    return Counter(clean(app.get(field)) or "Blank" for app in applications)


def count_multiselect(applications: list[dict], field: str) -> Counter:
    counts = Counter()

    for app in applications:
        raw = clean(app.get(field))
        if not raw:
            continue

        for item in raw.split(","):
            item = item.strip()
            if item:
                counts[item] += 1

    return counts


def count_minors(applications: list[dict]) -> Counter:
    """
    Minor(s) is open-ended text.

    Do not force responses into predefined categories because doing so could
    incorrectly merge different programs. Only obvious "no minor" answers are
    ignored.
    """
    ignored = {"", "n/a", "na", "none", "no", "not applicable"}
    counts = Counter()

    for app in applications:
        value = clean(app.get("minors"))
        if value.casefold() not in ignored:
            counts[value] += 1

    return counts


def display_label(field: str, value: str) -> str:
    # Country is open-ended. Whatever appears in the export is accepted.
    if field == "country":
        return value

    return LABELS.get(field, {}).get(value, value)


def counter_rows(
    counts: Counter,
    total: int,
    field: str = "",
    top: int | None = None,
) -> list[dict]:
    return [
        {
            "value": value,
            "label": display_label(field, value),
            "count": count,
            "percent": percent(count, total),
        }
        for value, count in counts.most_common(top)
    ]


def graduation_year_rows(
    applications: list[dict],
    total: int,
) -> tuple[list[dict], list[dict]]:
    """
    Graduation Year is open-ended from 2023 through 2050.

    The script derives the years from the CSV and sorts them chronologically.
    Unexpected values are reported separately.
    """
    valid = Counter()
    unexpected = Counter()

    for app in applications:
        raw = clean(app.get("graduationYear"))

        try:
            year = int(raw)
        except ValueError:
            unexpected[raw or "Blank"] += 1
            continue

        if GRADUATION_YEAR_MIN <= year <= GRADUATION_YEAR_MAX:
            valid[year] += 1
        else:
            unexpected[raw] += 1

    valid_rows = [
        {
            "value": str(year),
            "label": str(year),
            "count": valid[year],
            "percent": percent(valid[year], total),
        }
        for year in sorted(valid)
    ]

    return valid_rows, counter_rows(unexpected, total)


def word_count(text: str) -> int:
    return len(re.findall(r"\b[\w'-]+\b", text))


def essay_word_stats(
    applications: list[dict],
    field: str,
    minimum_words: int | None,
    maximum_words: int | None,
) -> dict:
    lengths = [
        word_count(clean(app.get(field)))
        for app in applications
        if clean(app.get(field))
    ]

    result = {
        "responses": len(lengths),
        "response_percent": percent(len(lengths), len(applications)),
        "average_words": round(statistics.mean(lengths), 1) if lengths else 0,
        "median_words": round(statistics.median(lengths), 1) if lengths else 0,
        "shortest_words": min(lengths) if lengths else 0,
        "longest_words": max(lengths) if lengths else 0,
    }

    if minimum_words is not None:
        result["below_minimum"] = sum(length < minimum_words for length in lengths)

    if maximum_words is not None:
        result["above_maximum"] = sum(length > maximum_words for length in lengths)

    return result


def theme_rows(
    applications: list[dict],
    field: str,
    themes: dict[str, list[str]],
    top: int | None = None,
) -> list[dict]:
    """
    Count how many applicants' responses match each broad text theme.

    Each applicant is counted at most once per theme, even if the response
    contains several matching keywords.
    """
    responses = [
        clean(app.get(field))
        for app in applications
        if clean(app.get(field))
    ]
    response_total = len(responses)

    counts = Counter()

    for theme, patterns in themes.items():
        compiled = [re.compile(pattern, re.IGNORECASE) for pattern in patterns]

        counts[theme] = sum(
            1
            for text in responses
            if any(pattern.search(text) for pattern in compiled)
        )

    return [
        {
            "label": theme,
            "count": count,
            "percent": percent(count, response_total),
        }
        for theme, count in counts.most_common(top)
        if count > 0
    ]


def build_summary(applications: list[dict], top: int) -> dict:
    total = len(applications)

    schools = count_field(applications, "school")
    years = count_field(applications, "year")

    uf_count = sum(schools.get(school, 0) for school in UF_SCHOOLS)
    non_uf_count = total - uf_count

    first_year_count = years.get("first_year", 0)
    first_year_uf = sum(
        1
        for app in applications
        if clean(app.get("year")) == "first_year"
        and clean(app.get("school")) in UF_SCHOOLS
    )

    graduation_years, unexpected_graduation_years = graduation_year_rows(
        applications, total
    )

    dietary = count_multiselect(applications, "diet")
    with_dietary_restrictions = sum(
        1 for app in applications if clean(app.get("diet"))
    )

    return {
        "overview": {
            "total_applications": total,
            "uf_students": {
                "count": uf_count,
                "percent": percent(uf_count, total),
            },
            "non_uf_students": {
                "count": non_uf_count,
                "percent": percent(non_uf_count, total),
            },
            "unique_schools": len(
                [school for school in schools if school != "Blank"]
            ),
            "first_year_students": {
                "count": first_year_count,
                "percent": percent(first_year_count, total),
                "uf": first_year_uf,
                "non_uf": first_year_count - first_year_uf,
            },
        },

        "personal_information": {
            "age": counter_rows(
                count_field(applications, "age"), total, "age"
            ),
            "country": counter_rows(
                count_field(applications, "country"), total, "country"
            ),
            "gender": counter_rows(
                count_field(applications, "gender"), total, "gender"
            ),
            "pronouns": counter_rows(
                count_field(applications, "pronouns"), total, "pronouns"
            ),
            "race_ethnicity": counter_rows(
                count_field(applications, "race"), total, "race"
            ),
            "sexual_orientation": counter_rows(
                count_field(applications, "orientation"),
                total,
                "orientation",
            ),
        },

        "education": {
            "school": {
                "uf": {
                    "count": uf_count,
                    "percent": percent(uf_count, total),
                },
                "non_uf": {
                    "count": non_uf_count,
                    "percent": percent(non_uf_count, total),
                },
                "unique_schools": len(
                    [school for school in schools if school != "Blank"]
                ),
                "top_schools": counter_rows(
                    schools, total, top=top
                ),
            },
            "level_of_study": counter_rows(
                count_field(applications, "level"), total, "level"
            ),
            "year_in_college": counter_rows(
                years, total, "year"
            ),
            "graduation_year": graduation_years,
            "unexpected_graduation_year_values": unexpected_graduation_years,
            "top_majors": counter_rows(
                count_multiselect(applications, "majors"),
                total,
                top=top,
            ),
            "top_minors": counter_rows(
                count_minors(applications),
                total,
                top=top,
            ),
        },

        "experience_and_preferences": {
            "hackathon_experience": counter_rows(
                count_field(applications, "experience"),
                total,
                "experience",
            ),
            "previous_uf_hackathon": counter_rows(
                count_field(applications, "ufHackathonExp"),
                total,
                "ufHackathonExp",
            ),
            "programming_project_experience": counter_rows(
                count_field(applications, "projectExperience"),
                total,
                "projectExperience",
            ),
            "shirt_size": counter_rows(
                count_field(applications, "shirtSize"),
                total,
                "shirtSize",
            ),
            "dietary_restrictions": {
                "applicants_with_restrictions": with_dietary_restrictions,
                "percent_with_restrictions": percent(
                    with_dietary_restrictions, total
                ),
                "selections": counter_rows(
                    dietary, total, "diet"
                ),
            },
        },

        "get_to_know_you": {
            "referral_sources": counter_rows(
                count_multiselect(applications, "referral"),
                total,
                "referral",
            ),

            "essay_1": {
                "word_stats": essay_word_stats(
                    applications, "essay1", 100, 250
                ),
                "motivation_themes": theme_rows(
                    applications, "essay1", ESSAY1_MOTIVATIONS
                ),
                "technical_areas_mentioned": theme_rows(
                    applications, "essay1", TECHNICAL_AREAS
                ),
            },

            "essay_2": {
                "word_stats": essay_word_stats(
                    applications, "essay2", 50, 250
                ),
                "project_domains": theme_rows(
                    applications, "essay2", ESSAY2_PROJECT_DOMAINS
                ),
                "technical_areas": theme_rows(
                    applications, "essay2", TECHNICAL_AREAS
                ),
                "collaboration_patterns": theme_rows(
                    applications, "essay2", ESSAY2_COLLABORATION
                ),
                "technologies_mentioned": theme_rows(
                    applications,
                    "essay2",
                    TECHNOLOGY_MENTIONS,
                    top=top,
                ),
            },

            "essay_3_optional": {
                "word_stats": essay_word_stats(
                    applications, "essay3", None, 250
                ),
                "interests_mentioned": theme_rows(
                    applications, "essay3", ESSAY3_INTERESTS
                ),
            },
        },
    }


def heading(title: str) -> None:
    print()
    print("=" * 78)
    print(title)
    print("=" * 78)


def print_rows(rows: list[dict]) -> None:
    for row in rows:
        print(
            f"{row['label']:<52}"
            f"{row['count']:>5}  "
            f"({row['percent']:>5.1f}%)"
        )


def print_theme_rows(rows: list[dict]) -> None:
    for row in rows:
        print(
            f"{row['label']:<52}"
            f"{row['count']:>5}  "
            f"({row['percent']:>5.1f}% of responses)"
        )


def print_essay_word_stats(title: str, stats: dict) -> None:
    print(f"\n{title}")
    print(
        f"  Responses:            "
        f"{stats['responses']} ({stats['response_percent']:.1f}%)"
    )
    print(f"  Average length:       {stats['average_words']} words")
    print(f"  Median length:        {stats['median_words']} words")
    print(
        f"  Observed range:       "
        f"{stats['shortest_words']}-{stats['longest_words']} words"
    )

    if "below_minimum" in stats:
        print(f"  Below minimum:        {stats['below_minimum']}")

    if "above_maximum" in stats:
        print(f"  Above maximum:        {stats['above_maximum']}")


def print_summary(summary: dict) -> None:
    overview = summary["overview"]

    heading("SWAMPHACKS APPLICATION OVERVIEW")
    print(f"Total applications:              {overview['total_applications']}")
    print(
        f"UF students:                     "
        f"{overview['uf_students']['count']} "
        f"({overview['uf_students']['percent']:.1f}%)"
    )
    print(
        f"Non-UF students:                 "
        f"{overview['non_uf_students']['count']} "
        f"({overview['non_uf_students']['percent']:.1f}%)"
    )
    print(f"Unique schools represented:      {overview['unique_schools']}")

    first_year = overview["first_year_students"]
    print(
        f"1st-year students:               "
        f"{first_year['count']} ({first_year['percent']:.1f}%)"
    )
    print(f"  UF 1st-years:                  {first_year['uf']}")
    print(f"  Non-UF 1st-years:              {first_year['non_uf']}")

    personal = summary["personal_information"]

    heading("PERSONAL INFORMATION")

    print("\nAge")
    print_rows(personal["age"])

    print("\nCountry")
    print_rows(personal["country"])

    print("\nGender")
    print_rows(personal["gender"])

    print("\nPronouns")
    print_rows(personal["pronouns"])

    print("\nRace / Ethnicity")
    print_rows(personal["race_ethnicity"])

    print("\nSexual Orientation")
    print_rows(personal["sexual_orientation"])

    education = summary["education"]

    heading("EDUCATION")

    print("\nTop Schools")
    print_rows(education["school"]["top_schools"])

    print("\nLevel of Study")
    print_rows(education["level_of_study"])

    print("\nYear in College")
    print_rows(education["year_in_college"])

    print(
        f"\nGraduation Year "
        f"({GRADUATION_YEAR_MIN}-{GRADUATION_YEAR_MAX})"
    )
    print_rows(education["graduation_year"])

    if education["unexpected_graduation_year_values"]:
        print("\nUnexpected Graduation Year Values")
        print_rows(education["unexpected_graduation_year_values"])

    print("\nTop Majors")
    print("Multi-select responses can appear in more than one category.")
    print_rows(education["top_majors"])

    print("\nTop Minors")
    print("Minor(s) is open text, so spelling variants remain separate.")
    print_rows(education["top_minors"])

    experience = summary["experience_and_preferences"]

    heading("EXPERIENCE & PREFERENCES")

    print("\nHackathons Participated In")
    print_rows(experience["hackathon_experience"])

    print("\nAttended SwampHacks / Another UF Hackathon Before")
    print_rows(experience["previous_uf_hackathon"])

    print("\nCreated an Independent Programming Project Before")
    print_rows(experience["programming_project_experience"])

    print("\nT-Shirt Size")
    print_rows(experience["shirt_size"])

    diet = experience["dietary_restrictions"]
    print(
        f"\nApplicants with dietary restrictions: "
        f"{diet['applicants_with_restrictions']} "
        f"({diet['percent_with_restrictions']:.1f}%)"
    )
    print("Multi-select responses may appear in more than one category.")
    print_rows(diet["selections"])

    get_to_know = summary["get_to_know_you"]

    heading("GET TO KNOW YOU")

    print("\nHow Applicants Learned About SwampHacks")
    print("This is multi-select, so percentages do not add to 100%.")
    print_rows(get_to_know["referral_sources"])

    print(
        "\nEssay insights are keyword-based aggregate trends. "
        "One response may match multiple themes."
    )

    essay1 = get_to_know["essay_1"]
    print_essay_word_stats(
        "Essay 1 - Why SwampHacks? (100-250 words)",
        essay1["word_stats"],
    )

    print("\nEssay 1 - Common Motivations")
    print_theme_rows(essay1["motivation_themes"])

    print("\nEssay 1 - Technical Areas Applicants Mention Wanting to Explore")
    print_theme_rows(essay1["technical_areas_mentioned"])

    essay2 = get_to_know["essay_2"]
    print_essay_word_stats(
        "Essay 2 - Project and collaboration (50-250 words)",
        essay2["word_stats"],
    )

    print("\nEssay 2 - Project Domains Mentioned")
    print_theme_rows(essay2["project_domains"])

    print("\nEssay 2 - Technical Areas Used in Past Projects")
    print_theme_rows(essay2["technical_areas"])

    print("\nEssay 2 - Collaboration Patterns Mentioned")
    print_theme_rows(essay2["collaboration_patterns"])

    print("\nEssay 2 - Technologies / Tools Mentioned")
    print_theme_rows(essay2["technologies_mentioned"])

    essay3 = get_to_know["essay_3_optional"]
    print_essay_word_stats(
        "Essay 3 - Optional get-to-know-you response (max 250 words)",
        essay3["word_stats"],
    )

    print("\nEssay 3 - Interests / Hobbies Mentioned")
    print_theme_rows(essay3["interests_mentioned"])


def main() -> None:
    args = parse_args()

    applications = load_applications(args.csv_file)

    if not applications:
        raise SystemExit("No submitted applications were found.")

    summary = build_summary(applications, args.top)
    print_summary(summary)

    if args.json_file:
        output_path = Path(args.json_file)
        output_path.write_text(
            json.dumps(summary, indent=2, ensure_ascii=False),
            encoding="utf-8",
        )
        print(f"\nSaved aggregate JSON summary to: {output_path}")


if __name__ == "__main__":
    main()

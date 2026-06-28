import argparse
import csv
import json


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Convert a CSV file into a JSON array."
    )
    parser.add_argument(
        "csv_file",
        help="Path to the input CSV file.",
    )
    parser.add_argument(
        "json_file",
        nargs="?",
        default=None,
        help="Path to the output JSON file. If omitted, prints to stdout.",
    )
    return parser.parse_args()


def read_csv(csv_path: str) -> list:
    with open(csv_path, newline="", encoding="utf-8") as csvfile:
        reader = csv.reader(csvfile)
        rows = [row for row in reader if row]

    if not rows:
        return []

    if all(len(row) == 1 for row in rows):
        return [row[0] for row in rows]

    if len(rows) > 1 and all(len(r) == len(rows[0]) for r in rows[1:]):
        header, *data_rows = rows
        return [dict(zip(header, row)) for row in data_rows]

    return rows


def write_json(json_path: str, data: list) -> None:
    data = list(dict.fromkeys(data))
    text = json.dumps(data, ensure_ascii=False, indent=2)
    if json_path:
        with open(json_path, "w", encoding="utf-8") as jsonfile:
            jsonfile.write(text)
    else:
        print(text)


def main() -> None:
    args = parse_args()
    data = read_csv(args.csv_file)
    write_json(args.json_file, data)


if __name__ == "__main__":
    main()

import boto3
import os
import csv
from dotenv import load_dotenv

load_dotenv()

ACCOUNT_ID: str = os.getenv("ACCOUNT_ID", "")
ACCESS_KEY_ID: str = os.getenv("ACCESS_KEY_ID", "")
ACCESS_KEY: str = os.getenv("ACCESS_KEY", "")
BUCKET_NAME: str = os.getenv("BUCKET_NAME", "")
DOWNLOAD_PATH: str= os.getenv("DOWNLOAD_PATH", "")
CSV_PATH: str = os.getenv("CSV_PATH", "")
EVENT_ID: str = os.getenv("EVENT_ID", "")

def download_file_from_r2(client, bucket, key, download_path):
    try:
        client.download_file(bucket, key, download_path)
    except Exception as e:
        print(f"Error downloading '{bucket}/{key}': {e}")
        return False

def main():
    try:
        s3 = boto3.client(
            service_name="s3",
            endpoint_url=f"https://{ACCOUNT_ID}.r2.cloudflarestorage.com",
            aws_access_key_id=ACCESS_KEY_ID,
            aws_secret_access_key=ACCESS_KEY,
            # aws_account_id=ACCOUNT_ID,
            region_name="auto",
        )
        print("R2 client initialized successfully.")

        os.makedirs(DOWNLOAD_PATH, exist_ok=True)

        with open(CSV_PATH, mode='r', newline='', encoding='utf-8') as csvfile:
            reader = csv.DictReader(csvfile)
            folder = EVENT_ID
            for row in reader:
                id = row['id'].strip()
                name = row['name'].strip()
                key = f"{folder}/{id}"
                filename= f"{DOWNLOAD_PATH}/{name}.pdf"
                download_file_from_r2(s3, BUCKET_NAME, key, filename)
                print(f"Saved {BUCKET_NAME}/{key}' to {filename}")

    except Exception as e:
        print(f"Failed to initialize R2 client or process files: {e}")

if __name__ == "__main__":
    main()


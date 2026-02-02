# Resume Downloader

Simple tool used for gathering all resumes from Cloudflare R2 and saving the
file as the user's name.

## Setup

1. Clone the repository if you haven't
2. Install dependencies through uv (install uv first)
Run the following:
```
  uv sync
```

3. Create a .env file. Copy the environment variables from .env.example

4. Get a csv file with a structure matching users.csv.example. If using the
database as of 2026/2/1, then the following query should be of help:

```sql
SELECT
  au.id,
  (app.application->>'firstName') || '_' || (app.application->>'lastName') AS name,
  er.role,
  app.status
FROM auth.users AS au
JOIN applications AS app ON app.user_id = au.id AND app.status = 'accepted'
JOIN event_roles AS er ON er.user_id = au.id AND er.role = 'attendee'
GROUP BY au.id, (app.application->>'firstName') || '_' || (app.application->>'lastName'), er.role, app.status;
```

This returns user.ids, semi-formatted name, and other parameters to visually
verify that these users have the correct role and application status. For the
name, please check the data to make sure that there are no spaces in any name
(which is easy to do in vim). Also look for any line with no name (search for
',,'), which are likely applicants who were admitted the day of the event and
had not submitted an application normally, as well as a resume. Remove these
people and try to add their resume to the pile manually if possible.

## Usage

### Download all resumes corresponding to csv

uv run main.py 


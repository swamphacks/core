-- Queries used for statistics, mainly used by overview dashboards etc

-- name: GetApplicationGenderSplit :one
SELECT
  COUNT(*) FILTER (WHERE application->>'gender' = 'man') AS male,
  COUNT(*) FILTER (WHERE application->>'gender' = 'woman') AS female,
  COUNT(*) FILTER (WHERE application->>'gender' = 'non-binary') AS non_binary,
  COUNT(*) FILTER (WHERE application->>'gender' = '') AS other
FROM applications
WHERE status <> 'started' AND status IS NOT NULL
  AND event_id = $1;

-- name: GetApplicationAgeSplit :one
SELECT
    COUNT(*) FILTER (WHERE (application->>'age')::int < 18) AS underage,
    COUNT(*) FILTER (WHERE (application->>'age')::int = 18) AS age_18,
    COUNT(*) FILTER (WHERE (application->>'age')::int = 19) AS age_19,
    COUNT(*) FILTER (WHERE (application->>'age')::int = 20) AS age_20,
    COUNT(*) FILTER (WHERE (application->>'age')::int = 21) AS age_21,
    COUNT(*) FILTER (WHERE (application->>'age')::int = 22) AS age_22,
    COUNT(*) FILTER (WHERE (application->>'age')::int >= 23) AS age_23_plus
FROM applications
WHERE status <> 'started' AND status IS NOT NULL
  AND event_id = $1;

-- name: GetApplicationRaceSplit :many
SELECT
    CASE 
        WHEN application->>'race' IS NOT NULL AND application->>'race' <> '' THEN application->>'race'
        WHEN application->>'race-other' IS NOT NULL AND application->>'race-other' <> '' THEN application->>'race-other'
        ELSE 'prefer_not_to_say'
    END AS race_group,
    COUNT(*) AS count
FROM applications
WHERE status <> 'started' AND status IS NOT NULL
  AND event_id = $1
GROUP BY 
    CASE 
        WHEN application->>'race' IS NOT NULL AND application->>'race' <> '' THEN application->>'race'
        WHEN application->>'race-other' IS NOT NULL AND application->>'race-other' <> '' THEN application->>'race-other'
        ELSE 'prefer_not_to_say'
    END
ORDER BY count DESC;

-- name: GetApplicationSchoolSplit :many
SELECT
    (application->>'school')::text AS school,
    COUNT(*) AS count
FROM applications
WHERE status <> 'started' AND status IS NOT NULL
  AND event_id = $1
GROUP BY (application->>'school')::text
ORDER BY count DESC;

-- name: GetApplicationMajorSplit :many
SELECT
    trim(major) AS major,
    COUNT(*) AS count
FROM applications,
LATERAL unnest(string_to_array(application->>'majors', ',')) AS major
WHERE status <> 'started' AND status IS NOT NULL
  AND event_id = $1
GROUP BY trim(major)
ORDER BY count DESC;

-- name: GetApplicationStatusSplit :one
SELECT
    COUNT(*) FILTER (WHERE status = 'started')       AS started,
    COUNT(*) FILTER (WHERE status = 'submitted')     AS submitted,
    COUNT(*) FILTER (WHERE status = 'under_review')  AS under_review,
    COUNT(*) FILTER (WHERE status = 'accepted')      AS accepted,
    COUNT(*) FILTER (WHERE status = 'rejected')      AS rejected,
    COUNT(*) FILTER (WHERE status = 'waitlisted')    AS waitlisted,
    COUNT(*) FILTER (WHERE status = 'withdrawn')     AS withdrawn
FROM applications
WHERE event_id = $1;

-- name: GetSubmissionTimes :many
SELECT
  date_trunc('day', submitted_at AT TIME ZONE 'US/Eastern')::date AS day,
  COUNT(*) AS count
FROM applications
WHERE event_id = $1 AND submitted_at IS NOT NULL
GROUP BY day
ORDER By day;

-- +goose Up

-- A policy must be explicitly enabled for an event before dispatch begins.
CREATE TABLE waitlist_dispatch_policies (
    hackathon_id text PRIMARY KEY REFERENCES hackathons(id) ON DELETE CASCADE,
    enabled boolean NOT NULL DEFAULT false,
    invitations_open_at timestamptz NOT NULL,
    online_join_closes_at timestamptz NOT NULL,
    in_person_opens_at timestamptz NOT NULL,
    invitations_close_at timestamptz NOT NULL,
    CHECK (invitations_open_at < invitations_close_at),
    CHECK (online_join_closes_at <= in_person_opens_at),
    CHECK (in_person_opens_at < invitations_close_at)
);

-- Online join order remains intact when an organizer records physical arrival.
ALTER TABLE waitlist
    ADD COLUMN in_person_joined_at timestamptz,
    ADD COLUMN in_person_recorded_by uuid REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX waitlist_in_person_order
    ON waitlist (hackathon_id, in_person_joined_at, user_id)
    WHERE in_person_joined_at IS NOT NULL;

-- Store invitations durably alongside the admission transaction.
-- Redis or email outages must not lose an invitation.
CREATE TABLE waitlist_invitation_outbox (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id uuid NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    offered_at timestamptz NOT NULL,
    confirmation_deadline timestamptz NOT NULL,
    recipient text NOT NULL CHECK (length(btrim(recipient)) > 0),
    first_name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    sent_at timestamptz,
    cancelled_at timestamptz,
    attempts integer NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    next_attempt_at timestamptz NOT NULL DEFAULT now(),
    last_error text,
    UNIQUE (application_id, offered_at),
    CHECK (confirmation_deadline > offered_at)
);

CREATE INDEX waitlist_invitation_outbox_pending
    ON waitlist_invitation_outbox (next_attempt_at, created_at, id)
    WHERE sent_at IS NULL AND cancelled_at IS NULL;

-- +goose Down

DROP TABLE waitlist_invitation_outbox;
DROP INDEX waitlist_in_person_order;
ALTER TABLE waitlist
    DROP COLUMN in_person_recorded_by,
    DROP COLUMN in_person_joined_at;
DROP TABLE waitlist_dispatch_policies;

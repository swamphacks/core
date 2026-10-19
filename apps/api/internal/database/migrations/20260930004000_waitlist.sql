-- +goose Up
CREATE TABLE waitlist (
    hackathon_id text NOT NULL REFERENCES hackathons(id),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (hackathon_id, user_id),
    FOREIGN KEY (user_id, hackathon_id)
        REFERENCES applications(user_id, hackathon_id) ON DELETE CASCADE
);

CREATE INDEX waitlist_join_order
    ON waitlist (hackathon_id, created_at, user_id);

-- +goose Down
DROP TABLE waitlist;

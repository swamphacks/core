-- +goose Up
-- +goose StatementBegin
CREATE TABLE campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    title TEXT CONSTRAINT title_charlim CHECK (char_length(title) <= 200),
    description TEXT CONSTRAINT description_charlim CHECK (char_length(description) <= 1000),
    recipient_addresses TEXT[],
    recipient_roles TEXT[],
    emails UUID[] REFERENCES campaign_emails(id),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE campaign_emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_addresses TEXT[],
    recipient_roles TEXT[] NOT NULL,
    send_from TEXT NOT NULL,
    subject TEXT CONSTRAINT subject_charlim CHECK (char_length(subject) < 78), -- From SMTP spec.
    body TEXT, -- Limit determined by mail clients.
    template UUID NOT NULL REFERENCES campaign_email_template(id),
    send_on TIMESTAMPTZ NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE campaign_email_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT CONSTRAINT title_charlim CHECK (char_length(title) <= 200),
    html TEXT, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Each row contains stats associated with specific a role and email, and all links used to modify its stats.
CREATE TABLE campaign_email_statistics (
    email_id UUID NOT NULL REFERENCES campaign_emails(id) ON DELETE CASCADE,
    role event_role_type NOT NULL,
    hit_count INT DEFAULT 0,
    unsubscribed_count INT DEFAULT 0,
    link UUID[] NOT NULL REFERENCES campaign_link(id),

    PRIMARY KEY (email_id, role)
);

CREATE TABLE campaign_links (
    -- While the redirect is done by an api endpoint, both the role type and redirect_to could be passed as query params.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role event_role_type NOT NULL,
    redirect_to TEXT NOT NULL,
    api_endpoint TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS campaigns
DROP TABLE IF EXISTS campaign_emails
DROP TABLE IF EXISTS campaign_email_templates
DROP TABLE IF EXISTS campaign_email_statistics
DROP TABLE IF EXISTS campaigns_links

-- +goose StatementEnd
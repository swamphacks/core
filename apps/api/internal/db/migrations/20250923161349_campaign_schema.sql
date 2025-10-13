-- +goose Up
-- +goose StatementBegin

CREATE TABLE campaign_email_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT CONSTRAINT title_charlim CHECK (char_length(title) <= 200),
    html TEXT, 

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE campaign_email_links (
    -- May be able to simplify the primary key being used here
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role event_role_type NOT NULL,
    hit_count INT DEFAULT 0,
    unsubscribed_count INT DEFAULT 0,

    -- Example endpoint: GET localhost/<endpoint>?redirect=<redirect_to>&role=<role>
    api_endpoint TEXT NOT NULL,
    redirect_to TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE campaign_emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_addresses TEXT[],
    recipient_roles TEXT[] NOT NULL,
    send_from TEXT NOT NULL,
    subject TEXT CONSTRAINT subject_charlim CHECK (char_length(subject) < 78), -- From SMTP spec.
    body TEXT, -- Char limit determined by mail clients.

    -- Metadata
    template UUID NOT NULL REFERENCES campaign_email_templates(id),
    send_on TIMESTAMPTZ NOT NULL,
    links UUID[],
    created_by UUID NOT NULL REFERENCES auth.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    title TEXT CONSTRAINT title_charlim CHECK (char_length(title) <= 200) NOT NULL,
    description TEXT CONSTRAINT description_charlim CHECK (char_length(description) <= 1000),
    recipient_roles event_role_type[] NOT NULL,
    recipient_addresses TEXT[],
    emails UUID[],

    created_by UUID NOT NULL REFERENCES auth.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_campaign_email_templates
BEFORE UPDATE ON campaign_email_templates
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER set_updated_at_campaign_email_links
BEFORE UPDATE ON campaign_email_links
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER set_updated_at_campaign_emails
BEFORE UPDATE ON campaign_emails
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER set_updated_at_campaigns
BEFORE UPDATE ON campaigns
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_updated_at_campaigns ON campaigns;
DROP TRIGGER IF EXISTS set_updated_at_campaign_emails ON campaign_emails;
DROP TRIGGER IF EXISTS set_updated_at_campaign_email_links ON campaign_email_links;
DROP TRIGGER IF EXISTS set_updated_at_campaign_email_templates ON campaign_email_templates;

DROP TABLE IF EXISTS campaigns;
DROP TABLE IF EXISTS campaign_emails;
DROP TABLE IF EXISTS campaign_email_links;
DROP TABLE IF EXISTS campaign_email_templates;

-- +goose StatementEnd
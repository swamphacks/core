-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE redeemables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID REFERENCES events(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount INT NOT NULL CHECK (amount >= 0),
    max_user_amount INT NOT NULL CHECK (max_user_amount >= 1),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_updated_at_redeemables
BEFORE UPDATE ON redeemables
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TABLE user_redemptions (
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE NOT NULL,
    redeemable_id UUID REFERENCES redeemables(id) ON DELETE CASCADE NOT NULL,
    amount INT NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, redeemable_id)
);

CREATE TRIGGER set_updated_at_user_redemptions
BEFORE UPDATE ON user_redemptions
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS user_redemptions;
DROP TRIGGER IF EXISTS set_updated_at_redeemables ON redeemables;
DROP TABLE IF EXISTS redeemables;
DROP TRIGGER IF EXISTS set_updated_at_user_redemptions ON user_redemptions;
-- +goose StatementEnd

ALTER TABLE users
  ADD COLUMN is_active          BOOLEAN     NOT NULL DEFAULT false,
  ADD COLUMN activation_token   TEXT,
  ADD COLUMN activation_expires TIMESTAMPTZ,
  ADD COLUMN last_email_sent_at TIMESTAMPTZ,
  ADD COLUMN email_send_count   INT         NOT NULL DEFAULT 0;
CREATE TABLE IF NOT EXISTS subscriptions (
  id uuid PRIMARY KEY,
  service_name text NOT NULL,
  price integer NOT NULL,
  user_id uuid NOT NULL,
  start_date text NOT NULL,
  end_date text NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
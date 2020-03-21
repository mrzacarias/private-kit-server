CREATE TABLE IF NOT EXISTS had_contact_with_infected(
  "uuid" UUID PRIMARY KEY,
  since_ts TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX had_contact_with_infected_since_ts_index ON had_contact_with_infected (since_ts);

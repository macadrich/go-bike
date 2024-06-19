CREATE TABLE IF NOT EXISTS coordinates (
    id SERIAL PRIMARY KEY,
    kiosk_id INTEGER,
    longitude NUMERIC(10,8) NOT NULL,
    latitude NUMERIC(11,8) NOT NULL
)
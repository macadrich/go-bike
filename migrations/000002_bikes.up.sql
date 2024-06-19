CREATE TABLE IF NOT EXISTS bikes (
    id SERIAL PRIMARY KEY,
    at TIMESTAMP NOT NULL,
    kiosk_id INTEGER,
    dock_number INTEGER,
    is_electric BOOL NOT NULL,
    is_available BOOL NOT NULL,
    battery INTEGER
)
CREATE TABLE IF NOT EXISTS geometry (
    id SERIAL PRIMARY KEY,
    kiosk_id INTEGER REFERENCES stations(id),
    type VARCHAR(50) NOT NULL,
    coordinates NUMERIC[]
)
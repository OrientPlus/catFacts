-- schema.sql
CREATE TABLE facts
(
    id SERIAL PRIMARY KEY,
    message TEXT,
    length INT,
    time_point TIMESTAMPTZ
);


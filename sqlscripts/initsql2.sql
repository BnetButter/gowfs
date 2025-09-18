-- ==========================================================
-- Core metadata table
-- ==========================================================
CREATE TABLE layer_metadata (
    id SERIAL PRIMARY KEY,
    layer_name TEXT UNIQUE,
    layer_title TEXT,
    user_id INTEGER
);

-- ==========================================================
-- User 0: already has test_layer_1 and test_layer_2
-- ==========================================================
CREATE TABLE test_layer_1 (
    fid SERIAL PRIMARY KEY,
    geom GEOMETRY(POINT, 4326),
    name TEXT,
    address TEXT
);

INSERT INTO test_layer_1 (geom, name, address) VALUES
(ST_GeomFromText('POINT(12.123 42.789)', 4326), 'Kevin', '4510 Laclede Ave'),
(ST_GeomFromText('POINT(14.128 12.238)', 4326), 'Lawrance', '2145 Tenbrink'),
(ST_GeomFromText('POINT(28.479 94.177)', 4326), 'Peter', '23 Olde Hamlet Dr.');

CREATE TABLE test_layer_2 (
    fid SERIAL PRIMARY KEY,
    geom GEOMETRY(POINT, 4326),
    username TEXT,
    filename TEXT
);

INSERT INTO test_layer_2 (geom, username, filename) VALUES
(ST_GeomFromText('POINT(19.123 62.789)', 4326), 'Kevin', 'main.c'),
(ST_GeomFromText('POINT(17.128 92.238)', 4326), 'Lawrance', 'foobar.txt'),
(ST_GeomFromText('POINT(21.479 34.177)', 4326), 'Peter', 'helloworld.cpp');

INSERT INTO layer_metadata (layer_name, layer_title, user_id) VALUES
('test_layer_1', 'Who Lives Where', 0),
('test_layer_2', 'File Owners', 0);

-- ==========================================================
-- User 1 Layers
-- ==========================================================
CREATE TABLE restaurants_sf (
    fid SERIAL PRIMARY KEY,
    geom GEOMETRY(POINT, 4326),
    name TEXT,
    cuisine TEXT
);

INSERT INTO restaurants_sf (geom, name, cuisine) VALUES
(ST_GeomFromText('POINT(-122.4194 37.7749)', 4326), 'Golden Gate Grill', 'American'),
(ST_GeomFromText('POINT(-122.414 37.781)', 4326), 'Dragon Wok', 'Chinese'),
(ST_GeomFromText('POINT(-122.407 37.783)', 4326), 'Pasta Amore', 'Italian');

CREATE TABLE parks_ny (
    fid SERIAL PRIMARY KEY,
    geom GEOMETRY(POINT, 4326),
    park_name TEXT,
    size_acres INT
);

INSERT INTO parks_ny (geom, park_name, size_acres) VALUES
(ST_GeomFromText('POINT(-73.9654 40.7829)', 4326), 'Central Park', 843),
(ST_GeomFromText('POINT(-73.968 40.660)', 4326), 'Prospect Park', 526),
(ST_GeomFromText('POINT(-73.971 40.676)', 4326), 'Brooklyn Botanic Garden', 52);

CREATE TABLE wifi_chicago (
    fid SERIAL PRIMARY KEY,
    geom GEOMETRY(POINT, 4326),
    hotspot_name TEXT,
    provider TEXT
);

INSERT INTO wifi_chicago (geom, hotspot_name, provider) VALUES
(ST_GeomFromText('POINT(-87.6298 41.8781)', 4326), 'Union Station WiFi', 'AT&T'),
(ST_GeomFromText('POINT(-87.6244 41.8827)', 4326), 'Millennium Park WiFi', 'Comcast'),
(ST_GeomFromText('POINT(-87.623 41.885)', 4326), 'Chicago Public Library WiFi', 'City of Chicago');

INSERT INTO layer_metadata (layer_name, layer_title, user_id) VALUES
('restaurants_sf', 'San Francisco Restaurants', 1),
('parks_ny', 'New York Parks', 1),
('wifi_chicago', 'Chicago Public WiFi Spots', 1);

-- ==========================================================
-- User 2 Layers
-- ==========================================================
CREATE TABLE crime_la (
    fid SERIAL PRIMARY KEY,
    geom GEOMETRY(POINT, 4326),
    crime_type TEXT,
    date DATE
);

INSERT INTO crime_la (geom, crime_type, date) VALUES
(ST_GeomFromText('POINT(-118.2437 34.0522)', 4326), 'Robbery', '2024-07-15'),
(ST_GeomFromText('POINT(-118.25 34.05)', 4326), 'Burglary', '2024-08-02'),
(ST_GeomFromText('POINT(-118.27 34.04)', 4326), 'Vandalism', '2024-09-01');

CREATE TABLE museums_dc (
    fid SERIAL PRIMARY KEY,
    geom GEOMETRY(POINT, 4326),
    museum_name TEXT,
    admission BOOLEAN
);

INSERT INTO museums_dc (geom, museum_name, admission) VALUES
(ST_GeomFromText('POINT(-77.026 38.891)', 4326), 'Smithsonian Air & Space', TRUE),
(ST_GeomFromText('POINT(-77.023 38.891)', 4326), 'National Gallery of Art', TRUE),
(ST_GeomFromText('POINT(-77.019 38.888)', 4326), 'US Holocaust Memorial', FALSE);

CREATE TABLE schools_boston (
    fid SERIAL PRIMARY KEY,
    geom GEOMETRY(POINT, 4326),
    school_name TEXT,
    level TEXT
);

INSERT INTO schools_boston (geom, school_name, level) VALUES
(ST_GeomFromText('POINT(-71.0589 42.3601)', 4326), 'Boston Latin School', 'High School'),
(ST_GeomFromText('POINT(-71.09 42.34)', 4326), 'James F. Condon Elementary', 'Elementary'),
(ST_GeomFromText('POINT(-71.07 42.36)', 4326), 'Boston Arts Academy', 'High School');

INSERT INTO layer_metadata (layer_name, layer_title, user_id) VALUES
('crime_la', 'Los Angeles Crime Reports', 2),
('museums_dc', 'Washington DC Museums', 2),
('schools_boston', 'Boston Public Schools', 2);

-- ==========================================================
-- User 3 Layers
-- ==========================================================
CREATE TABLE hospitals_texas (
    fid SERIAL PRIMARY KEY,
    geom GEOMETRY(POINT, 4326),
    hospital_name TEXT,
    capacity INT
);

INSERT INTO hospitals_texas (geom, hospital_name, capacity) VALUES
(ST_GeomFromText('POINT(-95.3698 29.7604)', 4326), 'Houston Medical Center', 1500),
(ST_GeomFromText('POINT(-96.797 32.7767)', 4326), 'Dallas General', 1200),
(ST_GeomFromText('POINT(-97.7431 30.2672)', 4326), 'Austin Health', 800);

CREATE TABLE libraries_portland (
    fid SERIAL PRIMARY KEY,
    geom GEOMETRY(POINT, 4326),
    library_name TEXT,
    open_hours TEXT
);

INSERT INTO libraries_portland (geom, library_name, open_hours) VALUES
(ST_GeomFromText('POINT(-122.6765 45.5231)', 4326), 'Portland Central Library', '9-5'),
(ST_GeomFromText('POINT(-122.68 45.52)', 4326), 'Belmont Library', '10-6'),
(ST_GeomFromText('POINT(-122.66 45.53)', 4326), 'Hollywood Library', '11-7');

CREATE TABLE breweries_denver (
    fid SERIAL PRIMARY KEY,
    geom GEOMETRY(POINT, 4326),
    brewery_name TEXT,
    founded_year INT
);

INSERT INTO breweries_denver (geom, brewery_name, founded_year) VALUES
(ST_GeomFromText('POINT(-104.9903 39.7392)', 4326), 'Mile High Brewing', 1995),
(ST_GeomFromText('POINT(-105.00 39.74)', 4326), 'Rocky Mountain Ales', 2005),
(ST_GeomFromText('POINT(-104.98 39.74)', 4326), 'Denver Beer Works', 2010);

CREATE TABLE bike_trails_seattle (
    fid SERIAL PRIMARY KEY,
    geom GEOMETRY(POINT, 4326),
    trail_name TEXT,
    length_miles DECIMAL
);

INSERT INTO bike_trails_seattle (geom, trail_name, length_miles) VALUES
(ST_GeomFromText('POINT(-122.3321 47.6062)', 4326), 'Burke-Gilman Trail', 18.8),
(ST_GeomFromText('POINT(-122.34 47.61)', 4326), 'Elliott Bay Trail', 5.0),
(ST_GeomFromText('POINT(-122.31 47.63)', 4326), 'Chief Sealth Trail', 4.5);

INSERT INTO layer_metadata (layer_name, layer_title, user_id) VALUES
('hospitals_texas', 'Texas Hospitals', 3),
('libraries_portland', 'Portland Libraries', 3),
('breweries_denver', 'Denver Breweries', 3),
('bike_trails_seattle', 'Seattle Bike Trails', 3);
CREATE TABLE layer_metadata (id SERIAL PRIMARY KEY, layer_name TEXT UNIQUE, layer_title TEXT);
INSERT INTO layer_metadata (layer_name, layer_title) VALUES ('test_layer_1', 'Who Lives Where');
INSERT INTO layer_metadata (layer_name, layer_title) VALUES ('test_layer_2', 'File Owners');

CREATE TABLE test_layer_1 (fid SERIAL PRIMARY KEY, geom GEOMETRY(POINT, 4326), name TEXT, address TEXT);
CREATE TABLE test_layer_2 (fid SERIAL PRIMARY KEY, geom GEOMETRY(POINT, 4326), username TEXT, filename TEXT);

INSERT INTO test_layer_1 (geom, name, address) VALUES (ST_GeomFromText('POINT(12.123 42.789)', 4326), 'Kevin', '4510 Laclede Ave');
INSERT INTO test_layer_1 (geom, name, address) VALUES (ST_GeomFromText('POINT(14.128 12.238)', 4326), 'Lawrance', '2145 Tenbrink');
INSERT INTO test_layer_1 (geom, name, address) VALUES (ST_GeomFromText('POINT(28.479 94.177)', 4326), 'Peter', '23 Olde Hamlet Dr.');

INSERT INTO test_layer_2 (geom, username, filename) VALUES (ST_GeomFromText('POINT(19.123 62.789)', 4326), 'Kevin', 'main.c');
INSERT INTO test_layer_2 (geom, username, filename) VALUES (ST_GeomFromText('POINT(17.128 92.238)', 4326), 'Lawrance', 'foobar.txt');
INSERT INTO test_layer_2 (geom, username, filename) VALUES (ST_GeomFromText('POINT(21.479 34.177)', 4326), 'Peter', 'helloworld.cpp');


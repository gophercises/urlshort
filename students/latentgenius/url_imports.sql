CREATE TABLE IF NOT EXISTS urlmaps (shortpath VARCHAR(30) PRIMARY KEY, url VARCHAR(256) NOT NULL);
INSERT INTO urlmaps(shortpath, url) VALUES (
"/urlshort-godoc", "https://godoc.org/github.com/gophercises/urlshort");
INSERT INTO urlmaps(shortpath, url) VALUES (
"/yaml-godoc", "https://godoc.org/gopkg.in/yaml.v2");
INSERT INTO urlmaps(shortpath, url) VALUES (
"/urlshort", "https://github.com/gophercises/urlshort");
INSERT INTO urlmaps(shortpath, url) VALUES (
"/urlshort-final", "https://github.com/gophercises/urlshort/tree/final");

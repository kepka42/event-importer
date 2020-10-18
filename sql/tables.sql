CREATE TABLE locations (
    id          INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    radius      INT NOT NULL,
    coordinates POINT NOT NULL,
    start_from  TIMESTAMP NOT NULL,
    title       VARCHAR(255) NOT NULL
);

CREATE TABLE points (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    location_id  INT UNSIGNED NOT NULL,
    photo        VARCHAR(255),
    gender       TINYINT(1) DEFAULT 0,
    age          TINYINT UNSIGNED,
    has_children TINYINT(1) NOT NULL,
    coordinates  POINT NOT NULL,
    is_tourist   TINYINT(1),
    FOREIGN KEY (location_id) REFERENCES locations(id)
);
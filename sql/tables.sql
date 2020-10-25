CREATE TABLE locations (
    id          INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    radius      INT NOT NULL,
    lat         DOUBLE(10, 7) NOT NULL,
    lng         DOUBLE(10, 7) NOT NULL,
    start_from  TIMESTAMP NULL DEFAULT NULL,
    title       VARCHAR(255) NOT NULL
);

CREATE TABLE points (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    location_id  INT UNSIGNED NOT NULL,
    photo        VARCHAR(255),
    gender       CHAR(6),
    age          TINYINT UNSIGNED,
    has_children TINYINT(1) NOT NULL,
    lat          DOUBLE(10, 7) NOT NULL,
    lng          DOUBLE(10, 7) NOT NULL,
    is_tourist   TINYINT(1),
    vk_user_id   INT,
    user_city    VARCHAR(255),
    user_city_id INT,
    created_at   TIMESTAMP,
    FOREIGN KEY (location_id) REFERENCES locations(id)
);
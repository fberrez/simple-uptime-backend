CREATE TYPE keyword_condition AS ENUM ('contains','not_contains');
CREATE TYPE route_type AS ENUM ('http','keyword','ping');

CREATE TABLE account (
    id SERIAL PRIMARY KEY,
    email VARCHAR(256) NOT NULL,
    password VARCHAR(2048) NOT NULL,
    last_connection timestamp
);

CREATE TABLE authentication(
    id SERIAL PRIMARY KEY,
    username VARCHAR(1024) NOT NULL,
    password VARCHAR(1024) NOT NULL
);

CREATE TABLE header (
    id SERIAL PRIMARY KEY,
    name VARCHAR(1024) NOT NULL,
    value VARCHAR(1024) NOT NULL
);

CREATE TABLE http (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(2048),
    ip INET,
    interval SMALLINT NOT NULL,
    authentication_id INT NOT NULL,
    status_code SMALLINT NOT NULL,
    FOREIGN KEY(authentication_id) REFERENCES authentication(id)
);

CREATE TABLE keyword (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(2048),
    ip INET,
    keyword VARCHAR(255) NOT NULL,
    condition keyword_condition NOT NULL,
    interval SMALLINT NOT NULL,
    authentication_id INT NOT NULL
);

CREATE TABLE ping (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(2048),
    ip INET,
    port INTEGER,
    interval SMALLINT NOT NULL
);

CREATE TABLE route (
    id SERIAL PRIMARY KEY,
    type route_type NOT NULL,
    http_id INT,
    keyword_id INT,
    ping_id INT,
    FOREIGN KEY(http_id) REFERENCES http(id),
    FOREIGN KEY(keyword_id) REFERENCES http(id),
    FOREIGN KEY(ping_id) REFERENCES ping(id)
);

CREATE TABLE route_header(
    route_id INT NOT NULL,
    header_id INT NOT NULL,
    PRIMARY KEY(route_id, header_id),
    FOREIGN KEY(route_id) REFERENCES route(id),
    FOREIGN KEY(header_id) REFERENCES header(id)
);

CREATE TABLE account_route(
    account_id INT NOT NULL,
    route_id INT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY(account_id, route_id),
    FOREIGN KEY(account_id) REFERENCES account(id),
    FOREIGN KEY(route_id) REFERENCES route(id)
);

CREATE TABLE result(
    id SERIAL PRIMARY KEY,
    route_id INT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status_code SMALLINT NOT NULL,
    is_success BOOLEAN NOT NULL,
    FOREIGN KEY(route_id) REFERENCES route(id)
);

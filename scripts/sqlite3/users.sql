-- users definition

CREATE TABLE users (
	name TEXT NOT NULL,
	passwd TEXT NOT NULL,
	last_login_time NUMERIC,
	CONSTRAINT NewTable_PK PRIMARY KEY (name)
);

INSERT INTO users (name, passwd,last_login_time) values ('admin', 'admin@123', '2024-04-29 12:00:00');
-- user
CREATE TABLE IF NOT EXISTS user (
    pk_user_id        INTEGER NOT NULL AUTO_INCREMENT,
    password          VARCHAR(255) NOT NULL,
    name              VARCHAR(255) NOT NULL UNIQUE,
    email             VARCHAR(255) NOT NULL UNIQUE,
    role              VARCHAR(10) NOT NULL,
    active            BOOLEAN,
    PRIMARY KEY (pk_user_id)
);

-- subscription sensor
CREATE TABLE IF NOT EXISTS subscription (
    fk_sensor_id      INTEGER NOT NULL,
    fk_user_id      INTEGER NOT NULL,
    show_entries    INTEGER NOT NULL DEFAULT 10,
    PRIMARY KEY (fk_sensor_id, fk_user_id),
    FOREIGN KEY (fk_sensor_id) REFERENCES sensor (pk_sensor_id) ON DELETE CASCADE,
    FOREIGN KEY (fk_user_id) REFERENCES user (pk_user_id) ON DELETE CASCADE
);

-- user
INSERT INTO user (password, name, email, role, active)
VALUES('$2a$14$BHkC8UDmVJ3YbUOjwZaEa.4T.kG54L2bRc1561R0067CG5MHok04S', 'Admin', 'admin@localhost', 'admin', 1); -- dudeli

INSERT INTO user (password, name, email, role, active)
VALUES('$2a$12$jlO46pJTt9xmwszAHaVi4OvqgyFVxko/lNCYwE2sLJtQ4mo97YQ9S', 'SimRacer', 'simracer@localhost', 'user', 1); -- iracing

-- sensor
INSERT INTO sensor (name, fk_sensor_type_id, description)
VALUES('Admin Basement', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Admins Basement Temperature');
INSERT INTO sensor (name, fk_sensor_type_id, description)
VALUES('Admin Living Room', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Admins Living Room Temperature');
INSERT INTO sensor (name, fk_sensor_type_id, description)
VALUES('Simracer Living Room', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Simracer Living Room Temperature');

-- subscription
INSERT INTO subscription (fk_sensor_id, fk_user_id)
VALUES(
	(select pk_sensor_id from sensor where name = 'Admin Basement'),
	(select pk_user_id from user where name = 'Admin')
);
INSERT INTO subscription (fk_sensor_id, fk_user_id)
VALUES(
	(select pk_sensor_id from sensor where name = 'Admin Living Room'),
	(select pk_user_id from user where name = 'Admin')
);
INSERT INTO subscription (fk_sensor_id, fk_user_id)
VALUES(
	(select pk_sensor_id from sensor where name = 'Simracer Living Room'),
	(select pk_user_id from user where name = 'SimRacer')
);





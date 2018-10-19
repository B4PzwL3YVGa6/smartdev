-- user
DROP TABLE user;

-- subscription
DROP TABLE subscription;

-- subscription
DELETE FROM subscription
WHERE fk_user_id = (select pk_user_id from user where name = 'SimRacer');

-- feed
DELETE FROM sensor
WHERE name = 'Admin Basement';
DELETE FROM sensor
WHERE name = 'Admin Living Room';
DELETE FROM sensor
WHERE name = 'Simracer Living Room';

-- user
DELETE FROM user
WHERE name = 'SimRacer';

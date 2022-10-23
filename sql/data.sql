drop table if exists users cascade;
drop table if exists transactions cascade;
drop table if exists details cascade;

CREATE TABLE users (
    ID integer not null UNIQUE CHECK (id >= 0) PRIMARY KEY,
    balance integer not null CHECK (balance >= 0),
    reserved integer not null CHECK (reserved >= 0)
);

CREATE TABLE details (
    ID SERIAL PRIMARY KEY,
    orderId integer not null CHECK (orderId >= 0),
    serviceId integer not null CHECK (serviceId >= 0),
    status bool not null
);

ALTER TABLE details
ADD CONSTRAINT details_unique UNIQUE (orderId, serviceId);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    fromId integer REFERENCES users (id) ON DELETE CASCADE,
    toId integer REFERENCES users (id) ON DELETE CASCADE,
    amount integer not null CHECK(amount > 0),
    date timestamp not null,
    type integer not null CHECK(type >= 0 and type < 5),
    detailsId integer REFERENCES details (id) ON DELETE CASCADE
);
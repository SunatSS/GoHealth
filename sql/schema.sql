-- table of customers
-- todo payment_sources system
CREATE TABLE customers 
(
    id       BIGSERIAL   PRIMARY KEY,
    name     TEXT        NOT NULL,
    phone    TEXT        NOT NULL UNIQUE,
    password TEXT        NOT NULL,
    address  TEXT,
    active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
    -- payment_sources TEXT[][]
);

--table of customers_tokens
CREATE TABLE customers_tokens
(
    token       TEXT        NOT NULL    UNIQUE,
    customer_id BIGINT      NOT NULL    REFERENCES customers,
    expires     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);

-- table of pharmacies
-- todo license system and chande default of active to false
-- todo working_hours system
CREATE TABLE pharmacies 
(
    id          BIGSERIAL   PRIMARY KEY,
    name        TEXT        NOT NULL,
    address     TEXT,
    admin_login TEXT        NOT NULL UNIQUE,
    password    TEXT        NOT NULL,
    contacts    TEXT[],
    active      BOOLEAN     NOT NULL DEFAULT true,
    created     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
    -- license     TEXT        NOT NULL UNIQUE,
    -- working_hours TIMESTAMP[]
);

-- table of medicines
-- todo add image
CREATE TABLE medicines 
(
    id              BIGSERIAL   PRIMARY KEY,
    name            TEXT        NOT NULL,
    manafacturer    TEXT        NOT NULL,
    description     TEXT        NOT NULL,
    compounds       TEXT[],
    pharm_id        BIGINT      NOT NULL REFERENCES pharmacies,
    price           INTEGER     NOT NULL CHECK ( price > 0 ),
    qty             INTEGER     NOT NULL DEFAULT 0,
    recipe_needed   BOOLEAN     NOT NULL DEFAULT FALSE,
    active          BOOLEAN     NOT NULL DEFAULT TRUE,
    created         TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--todo table of drivers
-- --table of drivers
-- CREATE TABLE drivers 
-- (
--     id          BIGSERIAL   PRIMARY KEY,
-- );


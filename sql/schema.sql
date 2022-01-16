-- table of customers
CREATE TABLE customers 
(
    id       BIGSERIAL   PRIMARY KEY,
    name     TEXT        NOT NULL,
    phone    TEXT        NOT NULL UNIQUE,
    password TEXT        NOT NULL,
    address  TEXT,
    is_admin BOOLEAN     NOT NULL DEFAULT FALSE,
    active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--table of customers_tokens
CREATE TABLE customers_tokens
(
    token       TEXT        NOT NULL    UNIQUE,
    customer_id BIGINT      NOT NULL    REFERENCES customers,
    expires     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);

-- table of medicines
CREATE TABLE medicines 
(
    id              BIGSERIAL   PRIMARY KEY,
    name            TEXT        NOT NULL,
    manafacturer    TEXT        NOT NULL,
    description     TEXT        NOT NULL,
    components      TEXT[],
    recipe_needed   BOOLEAN     NOT NULL DEFAULT FALSE,
    price           INTEGER     NOT NULL CHECK ( price > 0 ),
    qty             INTEGER     NOT NULL DEFAULT 0,
    pharmacy_name   TEXT        NOT NULL,
    active          BOOLEAN     NOT NULL DEFAULT TRUE,
    created         TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);
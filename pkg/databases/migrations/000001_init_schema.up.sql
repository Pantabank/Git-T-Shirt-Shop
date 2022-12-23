START TRANSACTION:

CREATE TYPE genders AS ENUM (men, women);

CREATE TYPE style_types AS ENUM ('plain_color', 'patterns', 'figures');

CREATE TYPE sizes AS ENUM ('xs', 's', 'm', 'l', 'xl');

CREATE TYPE statuses AS ENUM ('placed_order', 'paid', 'out_of_shipping', 'completed');

CREATE TABLE products (
    id            BIGSERIAL PRIMARY KEY,
    gender        genders,
    style_type    style_types,
    style_detail  varchar(255),
    size          sizes,
    price         int,
    enable        bool
);

CREATE TABLE orders (
    id                SERIAL PRIMARY KEY,
    shipping_address  JSONB
);

CREATE TABLE product_order (
    id            BIGSERIAL PRIMARY KEY,
    order_id      BIGINT,
    products      JSONB,
    qty           int,
    price         int,
    status        statuses,
    created_at    timestamp,
    enable        bool
);

ALTER TABLE product_order ADD FOREIGN KEY (order_id) REFERENCES orders ("id");

COMMIT;
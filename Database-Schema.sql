CREATE TYPE "genders" AS ENUM (
  'Men',
  'Women'
);

CREATE TYPE "style_types" AS ENUM (
  'plain_color',
  'patterns',
  'figures'
);

CREATE TYPE "sizes" AS ENUM (
  'XS',
  'S',
  'M',
  'L',
  'XL'
);

CREATE TYPE "statuses" AS ENUM (
  'placed_order',
  'paid',
  'out_of_shipping',
  'completed'
);

CREATE TABLE "products" (
  "id" serial PRIMARY KEY,
  "gender" genders,
  "style_type" style_types,
  "style_detail" varchar(255),
  "size" sizes,
  "price" int,
  "enable" bool
);

CREATE TABLE "orders" (
  "id" serial PRIMARY KEY,
  "shipping_address" JSONB
);

CREATE TABLE "product_order" (
  "id" serial PRIMARY KEY,
  "order_id" int,
  "products" JSONB,
  "qty" int,
  "price" int,
  "enable" bool,
  "status" statuses,
  "created_at" timestamp
);

ALTER TABLE "product_order" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

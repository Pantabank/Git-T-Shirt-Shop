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

CREATE TYPE "roles" AS ENUM (
  'users',
  'admin'
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
  "customer_id" int5555555,
  "shipping_address" JSONB,
  "qty" int,
  "price" int,
  "enable" bool,
  "status" statuses,
  "created_at" timestamp
);

CREATE TABLE "product_order" (
  "id" serial PRIMARY KEY,
  "order_id" int,
  "product_id" int,
  "products" JSONB
);

CREATE TABLE "users" (
  "id" serial PRIMARY KEY,
  "username" varchar(255),
  "password" varchar(255),
  "role" roles
);

ALTER TABLE "product_order" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "product_order" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("customer_id") REFERENCES "users" ("id");

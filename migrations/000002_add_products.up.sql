CREATE TABLE IF NOT EXISTS "products" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "url" VARCHAR NOT NULL,
    "seller_id" INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS "sellers" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "url" VARCHAR
);

ALTER TABLE "products" ADD FOREIGN KEY ("seller_id") REFERENCES "sellers" ("id") ON DELETE CASCADE ;

CREATE TABLE IF NOT EXISTS "image_product" (
    "id" SERIAL PRIMARY KEY,
    "image_url" VARCHAR NOT NULL,
    "product_id" INTEGER NOT NULL
);

ALTER TABLE "image_product" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "price_now" (
    "id" SERIAL PRIMARY KEY,
    "price" NUMERIC(12,2) NOT NULL,
    "currency" VARCHAR NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "product_id" INTEGER NOT NULL
);

ALTER TABLE "price_now" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE;
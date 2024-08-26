CREATE TABLE IF NOT EXISTS "group_product" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS "match_product_group" (
    "id" SERIAL PRIMARY KEY,
    "group_id" INTEGER NOT NULL,
    "product_id" INTEGER NOT NULL,
    UNIQUE (group_id, product_id)
);

ALTER TABLE "match_product_group" ADD FOREIGN KEY ("group_id") REFERENCES "group_product" ("id") ON DELETE CASCADE;
ALTER TABLE "match_product_group" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE;
CREATE TABLE IF NOT EXISTS "user_sub_product" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "group_product_id" INTEGER NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, group_product_id)
);

ALTER TABLE "user_sub_product" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "user_sub_product" ADD FOREIGN KEY ("group_product_id") REFERENCES "group_product" ("id") ON DELETE CASCADE;
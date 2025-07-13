CREATE TABLE "sub_categories" (
  "id" bigserial PRIMARY KEY,
  "category_id" bigint NOT NULL,
  "name" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),

   CONSTRAINT "sub_categories_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES "categories" ("id")
);

ALTER TABLE "products" ADD COLUMN "sub_category" varchar(255) NOT NULL DEFAULT '';
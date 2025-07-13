CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "phone_number" varchar(255) UNIQUE NOT NULL,
  "refresh_token" text NOT NULL,
  "password" varchar(255) NOT NULL,
  "is_admin" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "price" decimal(10,2) NOT NULL DEFAULT 0,
  "category" varchar(255) NOT NULL,
  "image_url" text[] NOT NULL DEFAULT '{}',
  "size" text[] NOT NULL DEFAULT '{}',
  "color" text[] NOT NULL DEFAULT '{}',
  "stock_quantity" bigint NOT NULL DEFAULT 0,
  "deleted_at" timestamptz NULL,
  "updated_by" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT "products_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users" ("id")
);

CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "user_name" varchar(255) NOT NULL,
  "user_phone_number" varchar(255) NOT NULL,
  "total_amount" decimal(10,2) NOT NULL,
  "status" varchar(255) NOT NULL,
  "shipping_address" text NOT NULL,
  "payment_status" bool NOT NULL,
  "deleted_at" timestamptz,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "order_items" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "size" varchar(255) NOT NULL,
  "color" varchar(255) NOT NULL,
  "quantity" int NOT NULL DEFAULT 1,
  "amount" decimal(10,2) NOT NULL,

  CONSTRAINT "order_item_order_id_fkey" FOREIGN KEY ("order_id") REFERENCES "orders" ("id"),
  CONSTRAINT "order_item_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "products" ("id")
);

CREATE TABLE "payments" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL,
  "transaction_id" varchar(255) NOT NULL,
  "amount" decimal(10,2) NOT NULL,
  "payment_method" varchar(255) NOT NULL,
  "payment_status" bool NOT NULL DEFAULT false,
  "paid_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT "payment_order_id_fkey" FOREIGN KEY ("order_id") REFERENCES "orders" ("id")
);

CREATE INDEX ON "products" ("price");

CREATE INDEX ON "products" ("deleted_at");

CREATE INDEX ON "orders" ("status");

CREATE INDEX ON "order_items" ("order_id");

CREATE INDEX ON "payments" ("order_id");

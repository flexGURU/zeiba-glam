ALTER TABLE "order_items" DROP CONSTRAINT "order_item_order_id_fkey";
ALTER TABLE "order_items" DROP CONSTRAINT "order_item_product_id_fkey";
ALTER TABLE "payments" DROP CONSTRAINT "payment_order_id_fkey";

DROP TABLE "payments";
DROP TABLE "order_items";
DROP TABLE "orders";
DROP TABLE "products";
DROP TABLE "users";
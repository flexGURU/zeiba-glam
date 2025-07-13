ALTER TABLE "products" DROP COLUMN "sub_category";
ALTER TABLE "sub_categories" DROP CONSTRAINT "sub_categories_category_id_fkey";

DROP TABLE "sub_categories";
-- Удаление индексов
DROP INDEX IF EXISTS idx_pma_attr_value;
DROP INDEX IF EXISTS idx_product_models_status;
DROP INDEX IF EXISTS idx_product_models_product_id;
DROP INDEX IF EXISTS idx_products_brand_id;
DROP INDEX IF EXISTS idx_categories_parent_id;

-- Удаление таблиц (в правильном порядке для избежания нарушений внешних ключей)
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS product_model_attributes;
DROP TABLE IF EXISTS product_models;
DROP TABLE IF EXISTS attribute_value_translations;
DROP TABLE IF EXISTS attribute_values;
DROP TABLE IF EXISTS attribute_translations;
DROP TABLE IF EXISTS attributes;
DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS product_translations;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS category_translations;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS brands;
drop table if exists model_badges;
drop table if exists badge_translations;
drop table if exists badges;
DROP TABLE IF EXISTS languages;

create table if not exists languages (
  code varchar(12) primary key,
  is_default boolean not null default false,
  is_active boolean not null default true
);


create table if not exists brands (
  id uuid primary key default gen_random_uuid(),
  title varchar(100) not null,
  slug varchar(255) not null unique,
  description text,
  logo_url text,
  is_active boolean not null default true,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp default null
);

create table if not exists categories (
  id uuid primary key default gen_random_uuid(),
  parent_id uuid references categories(id),
  slug varchar(255) not null unique,
  logo_url text,
  is_active boolean not null default true,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp default null
);

create index idx_categories_parent_id on categories(parent_id);


create table if not exists category_translations (
    id uuid primary key default gen_random_uuid(),
    category_id uuid not null references categories(id) on delete cascade,
    language_code varchar(12) not null references languages(code) on delete cascade,
    title varchar(255) not null,  -- переведённое название атрибута
    description text,
    unique(category_id, language_code)
);


create table if not exists products (
  id uuid primary key default gen_random_uuid(),
  brand_id uuid references brands(id),
  status varchar(64) not null check (status in ('draft', 'active', 'archived')),
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp default null
);

create index idx_products_brand_id on products(brand_id);


create table if not exists product_translations (
    id uuid primary key default gen_random_uuid(),
    product_id uuid not null references products(id) on delete cascade,
    language_code varchar(12) not null references languages(code) on delete cascade,
    title varchar(255) not null,  -- переведённое название атрибута
    description text,
    unique(product_id, language_code)
);


create table if not exists product_categories (
  product_id uuid not null references products(id) on delete cascade,
  category_id uuid not null references categories(id) on delete cascade,
  primary key (product_id, category_id)
);

create table if not exists attributes (
  id uuid primary key default gen_random_uuid(),
  slug varchar(255) not null unique,
  data_type   text not null check (data_type in ('string', 'number', 'boolean')),
  is_variant boolean not null default false,
  created_at timestamp not null
);

create table if not exists attribute_translations (
    id uuid primary key default gen_random_uuid(),
    attribute_id uuid not null references attributes(id) on delete cascade,
    language_code varchar(12) not null references languages(code) on delete cascade,
    title varchar(255) not null,  -- переведённое название атрибута
    unique(attribute_id, language_code)
);


create table if not exists attribute_values (
  id uuid primary key default gen_random_uuid(),
  attribute_id uuid not null references attributes(id) on delete cascade,
  value varchar(255) not null,
  meta varchar(255),
  meta_type varchar(255) check (meta_type in ('string', 'number', 'color')),
  created_at timestamp not null
);

create unique index unique_attr_value on attribute_values (attribute_id, value);

create table if not exists attribute_value_translations (
  id uuid primary key default gen_random_uuid(),
  attribute_value_id uuid not null references attribute_values(id) on delete cascade,
  language_code varchar(12) not null references languages(code) on delete cascade,
  value varchar(255) not null,  -- переведённое значение
  unique(attribute_value_id, language_code)
);


create table if not exists product_models (
  id uuid primary key default gen_random_uuid(),
  product_id uuid not null references products(id) on delete cascade,

  sku varchar(255) not null unique,
  slug text not null unique,

  status varchar(64) not null check (status in ('active', 'inactive')),

  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp default null

);

create index idx_product_models_product_id on product_models(product_id);

create index idx_product_models_status on product_models(status);

create table if not exists product_model_attributes (
    product_model_id   uuid not null references product_models(id) on delete cascade,
    attribute_id       uuid not null references attributes(id),
    attribute_value_id uuid not null references attribute_values(id),
    PRIMARY KEY (product_model_id, attribute_id)
);

create index idx_pma_attr_value on product_model_attributes(attribute_id, attribute_value_id);



/*
Запрос с выбором fallback локали

SELECT 
    COALESCE(t.name, d.name) AS attribute_name
FROM attributes a
LEFT JOIN attribute_translations t
    ON t.attribute_id = a.id AND t.language = 'fr'
LEFT JOIN attribute_translations d
    ON d.attribute_id = a.id AND d.language = 'en'  -- fallback
WHERE a.id = :attribute_id;
*/




 --=========================
/*
WITH desired AS (
    SELECT 
        (SELECT code FROM languages WHERE is_default = TRUE) AS default_lang,
        'fr'::varchar AS lang  -- язык, который хотим показать
)

SELECT 
    pm.id AS product_model_id,
    pm.sku,
    pm.slug,
    
    -- атрибут
    a.id AS attribute_id,
    COALESCE(at_desired.name, at_default.name) AS attribute_name,
    a.is_variant,
    a.data_type AS attribute_data_type,

    -- значение атрибута
    av.id AS attribute_value_id,
    COALESCE(avt_desired.value, avt_default.value, av.value) AS attribute_value,
    av.meta,
    av.meta_type

FROM product_models pm
JOIN product_model_attributes pma 
    ON pm.id = pma.product_model_id
JOIN attributes a 
    ON pma.attribute_id = a.id
JOIN attribute_values av 
    ON pma.attribute_value_id = av.id

-- переводы атрибутов
LEFT JOIN attribute_translations at_desired
    ON at_desired.attribute_id = a.id
    AND at_desired.language_code = (SELECT lang FROM desired)
LEFT JOIN attribute_translations at_default
    ON at_default.attribute_id = a.id
    AND at_default.language_code = (SELECT default_lang FROM desired)

-- переводы значений атрибутов
LEFT JOIN attribute_value_translations avt_desired
    ON avt_desired.attribute_value_id = av.id
    AND avt_desired.language_code = (SELECT lang FROM desired)
LEFT JOIN attribute_value_translations avt_default
    ON avt_default.attribute_value_id = av.id
    AND avt_default.language_code = (SELECT default_lang FROM desired)

WHERE pm.id = :product_model_id
ORDER BY a.is_variant DESC, a.id;


Результат

| product_model_id | sku      | attribute_name | is_variant | attribute_value | meta    | meta_type |
| ---------------- | -------- | -------------- | ---------- | --------------- | ------- | --------- |
| 101              | SHIRT-RM | Цвет           | true       | Красный         | #FF0000 | color   |
| 101              | SHIRT-RM | Размер         | true       | M               | 48      | number    |
| 101              | SHIRT-RM | Материал       | false      | Хлопок          | NULL    | string    |


*/


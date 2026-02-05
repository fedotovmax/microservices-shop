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


create table if not exists events (
  id uuid primary key default gen_random_uuid(),
  aggregate_id varchar(100) not null,
  event_topic varchar(100) not null,
  event_type varchar(100) not null, 
  payload jsonb not null,
  status varchar not null default 'new' check(status in ('new', 'done')),
  created_at timestamp not null,
  reserved_to timestamp default null
);




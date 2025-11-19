
create domain email_d as text
check (value ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

create domain phone_e164_d as text
check (VALUE ~ '^\+[1-9][0-9]{7,14}$');

create type gender_enum as enum ('male', 'female');

create type history_status_type as enum ('deleted', 'blocked', 'restored', 'unblocked', 'muted');

create table if not exists users (

  id uuid primary key default gen_random_uuid(),

  email email_d not null unique,

  phone phone_e164_d unique,

  password_hash text not null,

  is_email_verified boolean default false,

  is_phone_verified boolean default false,

  created_at timestamp default now(),

  updated_at timestamp default now()
);

create table if not exists status_history (
    id bigserial primary key,
    user_id uuid not null references users (id),
    status_type history_status_type not null,
    reason text,
    period_start timestamp not null default now(),  -- начало действия статуса
    period_end timestamp null,                      -- конец периода (null = бессрочно)
    performed_by uuid null                          -- кто сделал действие
);



create table if not exists profiles (

  user_id uuid primary key references users (id),

  last_name varchar(100),

  first_name varchar(100),

  middle_name varchar(100),

  birth_date DATE,

  gender gender_enum,

  avatar_url TEXT,

  updated_at timestamp default now()

);


create table if not exists change_password_codes (

  id serial primary key,

  user_id uuid not null references users (id) unique,

  code varchar(6) not null,

  validity_period timestamp not null default now() + interval '2 minutes',

  attempts smallint not null default 0,

  blocked_until timestamp default null

);
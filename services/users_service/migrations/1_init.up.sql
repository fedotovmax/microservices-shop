
create domain email_d as text
check (value ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

create domain phone_e164_d as text
check (VALUE ~ '^\+[1-9][0-9]{7,14}$');

create table if not exists users (

  id uuid primary key default gen_random_uuid(),

  email email_d not null unique,

  phone phone_e164_d unique,

  password_hash text not null,

  is_email_verified boolean default false,

  is_phone_verified boolean default false,

  created_at timestamp not null,

  updated_at timestamp not null,

  deleted_at timestamp null
);

create table if not exists profiles (

  user_id uuid primary key references users (id) on delete cascade,

  last_name varchar(100),

  first_name varchar(100),

  middle_name varchar(100),

  birth_date DATE,

  gender smallint default 1 check(gender in (1,2,3)),

  avatar_url TEXT,

  updated_at timestamp not null

);

create table if not exists change_password_codes (

  user_id uuid primary key references users(id) on delete cascade,

  code varchar(6) not null,

  code_expires_at timestamp not null,

  attempts smallint not null default 0,

  blocked_until timestamp default null

);

create table if not exists email_verification (
  
  link uuid primary key default gen_random_uuid(),

  user_id uuid not null references users (id) unique on delete cascade,

  link_expires_at timestamp not null
);

create table if not exists phone_verification (

  user_id uuid primary key references users(id) on delete cascade,

  code varchar(6) not null,

  code_expires_at timestamp not null,

  attempts smallint not null default 0,

  blocked_until timestamp default null
);

create domain email_d as text
check (value ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

create table if not exists sessions_users (
  uid uuid primary key,
  email email_d not null unique
);

create table if not exists sessions (
  id uuid primary key,
  uid uuid references sessions_users (uid) on delete cascade not null,
  refresh_hash text unique not null,
  ip inet not null,
  browser varchar(48) not null,
  browser_version varchar(24) not null,
  os varchar(48) not null,
  device varchar(24) not null,
  created_at timestamp not null,
  updated_at timestamp not null,
  revoked_at timestamp,
  expires_at timestamp not null
);

create index idx_sessions_uid on sessions(uid);


create table blacklist (
  uid uuid primary key references sessions_users(uid) on delete cascade,
  code varchar(6) not null, 
  code_expires_at timestamp not null
);


create table bypass (
  uid uuid primary key references sessions_users(uid) on delete cascade,
  code varchar(12) not null, 
  bypass_expires_at timestamp not null
);

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
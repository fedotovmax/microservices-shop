create table if not exists events (
  id uuid primary key default gen_random_uuid(),
  aggregate_id varchar(100) not null,
  event_topic varchar(100) not null,
  event_type varchar(100) not null, 
  payload jsonb not null,
  status varchar not null default 'new' check(status in ('new', 'done')),
  created_at timestamp default now(),
  reserved_to timestamp default null
);

-- TODO: add index on reserved_to, status
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

create index concurrently idx_events_new_unreserved_created_at
on events (created_at)
where status = 'new'
and reserved_to is null;
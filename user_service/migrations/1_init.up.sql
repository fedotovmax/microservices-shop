
create domain email_d as text
check (value ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

create table if not exists public.user (
  id UUID primary key default gen_random_uuid(),
  email email_d not null unique,
  first_name varchar(40) not null check (length(first_name) >= 2),
  last_name varchar(40) not null check (length(last_name) >= 2)
);
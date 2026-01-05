drop index if exists idx_sessions_user_id;


drop table if exists blacklist;
drop table if exists bypass;
drop table if exists sessions;
drop table if exists sessions_users;

drop table if exists events;

drop domain if exists email_d;
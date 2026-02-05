drop index if exists idx_sessions_user_id;

drop index if exists idx_trust_tokens_uid;

drop table if exists blacklist;
drop table if exists bypass;
drop table if exists trust_tokens;
drop table if exists sessions;
drop table if exists sessions_users;

drop index concurrently if exists idx_events_new_unreserved_created_at;

drop table if exists events;

drop domain if exists email_d;
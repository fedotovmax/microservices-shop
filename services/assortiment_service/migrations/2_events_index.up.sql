create index concurrently idx_events_new_unreserved_created_at
on events (created_at)
where status = 'new'
and reserved_to is null;
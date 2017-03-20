
-- +migrate Up
create table instances(
	id uuid primary key,
	name text,
	region text,
	attributes jsonb
);

-- +migrate Down
drop table instances;


-- +migrate Up

create table s3(
id uuid primary key,
bucket_name varchar,
attributes jsonb
);

-- +migrate Down
drop table s3;

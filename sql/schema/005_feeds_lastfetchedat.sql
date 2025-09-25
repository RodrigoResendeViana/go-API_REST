-- +goose Up

alter table feeds ADD column last_fetched_at TIMESTAMP;

-- +goose Down

alter table feeds DROP column last_fetched_at;
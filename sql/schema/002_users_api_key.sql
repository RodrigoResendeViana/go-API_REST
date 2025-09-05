-- +goose Up

alter table users ADD column api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT(
    encode(sha256(random()::text::bytea), 'hex')
);

-- +goose Down

alter table users DROP column api_key;
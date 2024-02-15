CREATE TABLE IF NOT EXISTS "login_log" (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL,
    is_authorized boolean NOT NULL,
    ip_address varchar(255) NOT NULL,
    timezone varchar NOT NULL,
    lat numeric NOT NULL,
    lon numeric NOT NULL,
    access_time timestamp(0) NOT NULL
);

ALTER TABLE "login_log" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
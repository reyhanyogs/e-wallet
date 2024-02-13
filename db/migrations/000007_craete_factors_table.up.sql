CREATE TABLE "factors" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "pin" varchar(100) NOT NULL
);

ALTER TABLE "factors" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "user_id" int NOT NULL,
  "account_number" varchar(100),
  "balance" numeric(19,2)
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
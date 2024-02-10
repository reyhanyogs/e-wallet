CREATE TABLE "topup" (
    "id" varchar(100) PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "status" int NOT NULL,
    "snap_url" varchar(255)
);

ALTER TABLE "topup" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
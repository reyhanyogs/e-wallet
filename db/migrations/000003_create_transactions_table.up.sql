CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "sof_number" varchar(100) NOT NULL,
  "dof_number" varchar(100) NOT NULL,
  "amount" numeric(19,2),
  "transaction_type" varchar(1),
  "account_id" int NOT NULL,
  "transaction_datetime" timestamp
);

ALTER TABLE "transactions" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
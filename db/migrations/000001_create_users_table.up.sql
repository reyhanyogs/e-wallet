CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "full_name" varchar NOT NULL,
    "phone" varchar NOT NULL,
    "username" varchar NOT NULL,
    "password" varchar NOT NULL,
    "email_verified_at" timestamp
)
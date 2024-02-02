CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "full_name" varchar(100) NOT NULL,
    "phone" varchar(100) NOT NULL,
    "email" varchar(100) NOT NULL,
    "username" varchar(100) NOT NULL,
    "password" varchar(100) NOT NULL,
    "email_verified_at" timestamp
)
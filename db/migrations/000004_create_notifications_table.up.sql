CREATE TABLE "notifications" (
    "id" bigserial PRIMARY KEY,
    "user_id" int NOT NULL,
    "status" int,
    "title" text,
    "body" text,
    "is_read" int,
    "created_at" timestamp(0)
)
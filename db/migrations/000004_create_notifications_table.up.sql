CREATE TABLE "notifications" (
    "id" bigserial PRIMARY KEY,
    "user_id" int NOT NULL,
    "status" int NOT NULL,
    "title" text NOT NULL,
    "body" text NOT NULL,
    "is_read" int NOT NULL,
    "created_at" timestamp(0) NOT NULL
)
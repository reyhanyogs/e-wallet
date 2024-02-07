CREATE TABLE "notifications" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "status" bigint NOT NULL,
    "title" text NOT NULL,
    "body" text NOT NULL,
    "is_read" bigint NOT NULL,
    "created_at" timestamp(0) NOT NULL
)
CREATE TABLE "requests" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" varchar NOT NULL,
    "video_url" varchar,
    "zip_output_url" varchar,
    "status" varchar,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    finished_at timestamp
)
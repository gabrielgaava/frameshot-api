CREATE TABLE "requests" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" varchar NOT NULL,
    "user_email" varchar NOT NULL,
    "video_size" int,
    "video_key" varchar,
    "zip_output_key" varchar,
    "status" varchar,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    finished_at timestamp
)
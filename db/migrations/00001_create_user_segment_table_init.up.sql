CREATE TABLE IF NOT EXISTS "user-segment"(
    "id" bigserial PRIMARY KEY
        constraint banner_pkey
            NOT NULL,
    "created_at" timestamp DEFAULT (now()),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    "user_id" varchar NOT NULL,
    "segment" varchar NOT NULL
    );

CREATE TABLE if NOT EXISTS "archived-user-segment" (
    "id" bigserial PRIMARY KEY
        constraint banner_pkey
            NOT NULL,
    "created_at" timestamp DEFAULT (now()),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    "user_id" varchar NOT NULL,
    "segment" varchar NOT NULL
);

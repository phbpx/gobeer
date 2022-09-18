CREATE TABLE IF NOT EXISTS "beers" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "brewery" VARCHAR(255) NOT NULL,
    "style" VARCHAR(255) NOT NULL,
    "abv" FLOAT NOT NULL,
    "short_desc" VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS "reviews" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,
    "beer_id" UUID NOT NULL REFERENCES "beers" ("id") ON DELETE CASCADE,
    "user_id" UUID NOT NULL,
    "comment" TEXT NOT NULL,
    "score" INTEGER NOT NULL
);

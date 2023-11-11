CREATE TABLE "users" (
                         "username" varchar PRIMARY KEY,
                         "hashed_password" varchar NOT NULL,
                         "full_name" varchar NOT NULL,
                         "email" varchar UNIQUE NOT NULL,
                         "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00+00',
                         "created_at" timestamptz DEFAULT 'now()'
);

CREATE TABLE "accounts" (
                            "id" BIGSERIAL PRIMARY KEY,
                            "owner" varchar NOT NULL,
                            "balance" bigint NOT NULL,
                            "currency" varchar NOT NULL,
                            "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "entries" (
                           "id" BIGSERIAL PRIMARY KEY,
                           "account_id" bigint NOT NULL,
                           "amount" bigint NOT NULL,
                           "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "transfers" (
                             "id" BIGSERIAL PRIMARY KEY,
                             "from_account_id" bigint NOT NULL,
                             "to_account_id" bigint NOT NULL,
                             "amount" bigint NOT NULL,
                             "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "managers" (
                            "id" BIGSERIAL PRIMARY KEY,
                            "full_name" varchar NOT NULL,
                            "username" varchar NOT NULL,
                            "hashed_password" varchar NOT NULL
);

CREATE TABLE "req_queue" (
                             "req_id" BIGSERIAL PRIMARY KEY,
                             "req" json NOT NULL
);

CREATE TABLE "manager_decision" (
                                    "dec_id" BIGSERIAL PRIMARY KEY,
                                    "man_id" bigint NOT NULL,
                                    "decision" boolean NOT NULL,
                                    "message" varchar
);

CREATE INDEX ON "accounts" ("owner");

CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "manager_decision" ADD FOREIGN KEY ("dec_id") REFERENCES "managers" ("id");

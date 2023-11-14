
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


ALTER TABLE "manager_decision" ADD FOREIGN KEY ("man_id") REFERENCES "managers" ("id");

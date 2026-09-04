ALTER TABLE "accounts" ADD COLUMN "balance" bigint NOT NULL DEFAULT 0;

UPDATE "accounts" a
SET "balance" = COALESCE(sub.total, 0)
FROM (
  SELECT "account_id", SUM("amount") AS "total"
  FROM "entries"
  GROUP BY "account_id"
) sub
WHERE a."id" = sub."account_id";

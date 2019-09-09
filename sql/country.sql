CREATE TABLE "country" (
    "id" text,
    "name" text,
    PRIMARY KEY (id)
);

INSERT INTO country (id, name)
    VALUES ('sk', 'slovensko');

-- hset currency-codes eur euro
-- hset currency-codes usd dollar
-- hgetall currency-codes

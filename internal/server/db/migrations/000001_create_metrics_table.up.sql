BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS metrics(
   id varchar PRIMARY KEY,
   type VARCHAR (50) NOT NULL,
   delta BIGINT NULL,
   value DOUBLE PRECISION NULL
);

COMMIT;

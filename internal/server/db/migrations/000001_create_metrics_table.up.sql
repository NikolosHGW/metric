BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS metrics(
   id varchar PRIMARY KEY,
   type VARCHAR (50) NOT NULL,
   delta INTEGER NULL,
   value DOUBLE PRECISION NULL
);

COMMIT;

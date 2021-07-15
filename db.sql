DROP TABLE IF EXISTS "accounts" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "users_auth" CASCADE;
DROP TABLE IF EXISTS "transactions" CASCADE;

CREATE TABLE IF NOT EXISTS "users" (
  "user_id" int GENERATED ALWAYS AS IDENTITY (START WITH 100) PRIMARY KEY,
  "username" varchar(50) UNIQUE NOT NULL,
  "password" varchar(255) NOT NULL,
  "name" varchar(50) NOT NULL,
  "pan" varchar(11) UNIQUE NOT NULL,
  "location" varchar(50) NOT NULL,
  "address" varchar(255) NOT NULL,
  "gender" varchar(10) NOT NULL,
  "nationality" varchar(50) NOT NULL,
  "contact_number" BIGINT UNIQUE NOT NULL,
  "status" smallint NOT NULL DEFAULT '1'
);

CREATE TABLE IF NOT EXISTS "accounts" (
	"account_id" int GENERATED ALWAYS AS IDENTITY (START WITH 1000) PRIMARY KEY,
	"user_id" int NOT NULL,
	"amount" decimal NOT NULL,
	"created_on" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"status" smallint NOT NULL DEFAULT '1',
	CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
	  REFERENCES users(user_id)
);


CREATE TABLE IF NOT EXISTS "transactions" (
  "transaction_id" int GENERATED ALWAYS AS IDENTITY (START WITH 10000) PRIMARY KEY,
  "account_id" int NOT NULL,
  "amount" decimal(10,2) NOT NULL,
  "transaction_type" varchar(10) NOT NULL,
  "transaction_date" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_account
      FOREIGN KEY(account_id) 
	  REFERENCES accounts(account_id)
)

CREATE TABLE IF NOT EXISTS "accounts" (
  "account_id" serial PRIMARY KEY,
  "name" varchar(50) NOT NULL,
  "username" varchar(50) UNIQUE NOT NULL,
  "password" varchar(255) NOT NULL,
  "pan" varchar(11) UNIQUE NOT NULL,
  "location" varchar(50) NOT NULL,
  "address" varchar(255) NOT NULL,
  "gender" varchar(10) NOT NULL,
  "nationality" varchar(50) NOT NULL,
  "contact_number" BIGINT UNIQUE NOT NULL
);

INSERT into "accounts" ("name","username","password","location","pan","address","contact_number","gender","nationality")
values 
('Chakra', 'chakrasm','pass123','Bangalore', 'AFGYT4351K', '#256, 2nd Cross, Vijayanagar', 8867508500, 'Male', 'Indian'),
('Mithun', 'mithu','pass456', 'Dvg', 'GFNYH1242G', '#2, 7th Cross, Some Extension', 9108999966, 'Male', 'Indian')
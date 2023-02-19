CREATE TABLE account(
	acc_number BIGSERIAL PRIMARY KEY,
	name varchar(150) NOT NULL UNIQUE,
	balance bigint NOT NULL,
	currency varchar(3) NOT NULL,
	created_at timestamptz NOT NULL DEFAULT NOW()
);


CREATE TABLE transactions(
	id bigint primary key,
	acc_number bigint,
	amount bigint NOT NULL,
	date timestamptz NOT NULL default NOW()
);

CREATE TABLE transfers(
	id bigint primary key,
	amount bigint NOT NULL,
	debit bigint NOT NULL,
	credit bigint NOT NULL,
	date timestamptz NOT NULL default NOW()
);

ALTER TABLE transactions
ADD FOREIGN KEY (acc_number)
REFERENCES account(acc_number);

ALTER TABLE transfers
ADD FOREIGN KEY (debit)
REFERENCES account(acc_number);

ALTER TABLE transfers
ADD FOREIGN KEY (credit)
REFERENCES account(acc_number);

CREATE INDEX idx_acc_acc_num
ON "account"("acc_number");

CREATE INDEX idx_transaction
ON "transactions"("acc_number");

CREATE INDEX idx_transfers_sender
ON "transfers"("debit");

CREATE INDEX idx_transfers_receiver
ON "transfers"("credit");

CREATE INDEX idx_transfers
ON transfers("debit", "credit");

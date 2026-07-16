-- +goose Up

------------------------------------------------------------
-- Accounts
------------------------------------------------------------
INSERT INTO accounts (name)
VALUES
('HDFC Savings'),
('Cash Wallet');

------------------------------------------------------------
-- Jars
------------------------------------------------------------
INSERT INTO jars (name, allocation_type, value)
VALUES
('Necessities', 'percentage', 55),
('Leisure', 'percentage', 20),
('Investments', 'percentage', 15),
('Giving', 'remainder', 0);

------------------------------------------------------------
-- Transactions
------------------------------------------------------------

-- June Salary
INSERT INTO transactions
(type, amount, name, date, to_account_id, category, is_master_income)
VALUES
('income',85000,'June Salary','2026-06-01',1,'Income',TRUE);

-- June Allocations
INSERT INTO transactions
(type, amount, name, date, from_account_id, to_account_id, category, jar_id)
VALUES
('transfer',46750,'Necessities Allocation','2026-06-01',1,1,'Transfers',1),
('transfer',17000,'Leisure Allocation','2026-06-01',1,1,'Transfers',2),
('transfer',12750,'Investment Allocation','2026-06-01',1,1,'Transfers',3),
('transfer',8500,'Giving Allocation','2026-06-01',1,1,'Transfers',4),
('transfer',4000,'Cash Withdrawal','2026-06-03',1,2,'Transfers',NULL);

-- June Expenses
INSERT INTO transactions
(type, amount, name, date, from_account_id, category, jar_id)
VALUES
('expense',18000,'House Rent','2026-06-02',1,'Housing',1),
('expense',3200,'Groceries','2026-06-05',1,'Food',1),
('expense',500,'Coffee','2026-06-08',2,'Food',2),
('expense',5000,'Nifty 50 SIP','2026-06-15',1,'Investment',3),
('expense',1000,'Church Offering','2026-06-20',1,'Donations',4);

-- July Salary
INSERT INTO transactions
(type, amount, name, date, to_account_id, category, is_master_income)
VALUES
('income',90000,'July Salary','2026-07-01',1,'Income',TRUE);

-- July Allocations
INSERT INTO transactions
(type, amount, name, date, from_account_id, to_account_id, category, jar_id)
VALUES
('transfer',49500,'Necessities Allocation','2026-07-01',1,1,'Transfers',1),
('transfer',18000,'Leisure Allocation','2026-07-01',1,1,'Transfers',2),
('transfer',13500,'Investment Allocation','2026-07-01',1,1,'Transfers',3),
('transfer',9000,'Giving Allocation','2026-07-01',1,1,'Transfers',4),
('transfer',5000,'Cash Withdrawal','2026-07-03',1,2,'Transfers',NULL);

-- July Expenses
INSERT INTO transactions
(type, amount, name, date, from_account_id, category, jar_id)
VALUES
('expense',18000,'House Rent','2026-07-02',1,'Housing',1),
('expense',3500,'Groceries','2026-07-04',1,'Food',1),
('expense',1200,'Fuel','2026-07-06',1,'Transport',1),
('expense',450,'Cafe Coffee Day','2026-07-08',2,'Food',2),
('expense',1200,'Movie Night','2026-07-12',2,'Entertainment',2),
('expense',5000,'Nifty 50 SIP','2026-07-15',1,'Investment',3),
('expense',1000,'Church Offering','2026-07-20',1,'Donations',4);

------------------------------------------------------------
-- Goals
------------------------------------------------------------
INSERT INTO goals
(name, target_amount, saved_amount, deadline, allocation_percentage)
VALUES
('Emergency Fund',300000,35000,'2027-12-31',20),
('Japan Trip',150000,10000,'2027-06-01',10);

------------------------------------------------------------
-- Goal Transactions
------------------------------------------------------------
INSERT INTO goal_transactions
(goal_id, amount, transaction_type, notes)
VALUES
(1,17000,'allocation','June allocation'),
(1,18000,'allocation','July allocation'),
(2,10000,'manual_contribution','Initial contribution');

-- +goose Down

DELETE FROM goal_transactions;
DELETE FROM goals;
DELETE FROM transactions;
DELETE FROM jars;
DELETE FROM accounts;

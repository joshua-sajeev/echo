package db

import (
	"context"
	"log"
	"os"
)

// IsDemo returns true when APP_ENV=demo.
func IsDemo() bool {
	return os.Getenv("APP_ENV") == "demo"
}

// SeedDemo wipes all user data and inserts realistic dummy data.
// Only callable in demo mode. Safe to call repeatedly.
func SeedDemo(ctx context.Context) {
	if !IsDemo() {
		log.Println("SeedDemo: not in demo mode, skipping")
		return
	}

	steps := []struct {
		name string
		sql  string
	}{
		// wipe in FK-safe order
		{"wipe tx_templates", `DELETE FROM tx_templates`},
		{"wipe transactions", `DELETE FROM transactions`},
		{"wipe accounts", `DELETE FROM accounts`},
		// reset sequences so IDs start from 1 again
		{"reset account seq", `ALTER SEQUENCE accounts_id_seq RESTART WITH 1`},
		{"reset tx seq", `ALTER SEQUENCE transactions_id_seq RESTART WITH 1`},
		{"reset tmpl seq", `ALTER SEQUENCE tx_templates_id_seq RESTART WITH 1`},

		// accounts
		{"seed accounts", `
			INSERT INTO accounts (name) VALUES
				('HDFC Savings'),
				('ICICI Credit Card'),
				('Cash')
		`},

		// transactions — mix of income/expense/transfer across ~3 months
		// uses subqueries to resolve account IDs by name (safe if names are unique)
		{"seed transactions", `
			INSERT INTO transactions (type, amount, name, date, from_account_id, to_account_id, jar_id, is_master_income)
			SELECT * FROM (VALUES
				-- master income (May)
				('income'::text,   85000, 'May Salary',          '2026-05-01'::date, NULL, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, true),
				-- expenses May
				('expense'::text,   4200, 'Rent',                '2026-05-02'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='Necessities'), false),
				('expense'::text,    850, 'Electricity Bill',    '2026-05-04'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='Necessities'), false),
				('expense'::text,    320, 'Zomato',              '2026-05-05'::date, (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, (SELECT id FROM jars WHERE name='Leisure'),      false),
				('expense'::text,   1000, 'SIP - Nifty 50',     '2026-05-06'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='SIP'),           false),
				('expense'::text,   8500, 'Chitty Instalment',   '2026-05-07'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='Chitty'),        false),
				('expense'::text,    540, 'Grocery - DMart',     '2026-05-09'::date, (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, (SELECT id FROM jars WHERE name='Necessities'), false),
				('expense'::text,    199, 'Netflix',             '2026-05-10'::date, (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, (SELECT id FROM jars WHERE name='Leisure'),      false),
				('expense'::text,    650, 'Petrol',              '2026-05-12'::date, (SELECT id FROM accounts WHERE name='Cash'),             NULL, (SELECT id FROM jars WHERE name='Necessities'), false),
				('expense'::text,   8500, 'Charity - CRY',       '2026-05-13'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='Charity'),       false),
				('expense'::text,    420, 'Swiggy',              '2026-05-15'::date, (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, (SELECT id FROM jars WHERE name='Leisure'),      false),
				('expense'::text,   1200, 'Medicine',            '2026-05-17'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='Necessities'), false),
				('expense'::text,    890, 'Uber rides',          '2026-05-20'::date, (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, (SELECT id FROM jars WHERE name='Leisure'),      false),
				-- transfer
				('transfer'::text,  5000, 'CC Bill Payment',     '2026-05-22'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'), (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, false),
				-- income (Apr)
				('income'::text,   85000, 'April Salary',        '2026-04-01'::date, NULL, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, true),
				-- expenses Apr
				('expense'::text,   4200, 'Rent',                '2026-04-02'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='Necessities'), false),
				('expense'::text,   1000, 'SIP - Nifty 50',     '2026-04-05'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='SIP'),           false),
				('expense'::text,    780, 'Grocery - BigBasket', '2026-04-07'::date, (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, (SELECT id FROM jars WHERE name='Necessities'), false),
				('expense'::text,    450, 'Zomato',              '2026-04-10'::date, (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, (SELECT id FROM jars WHERE name='Leisure'),      false),
				('expense'::text,   8500, 'Chitty Instalment',   '2026-04-07'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='Chitty'),        false),
				('expense'::text,   8500, 'Charity - GiveIndia', '2026-04-13'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='Charity'),       false),
				('expense'::text,    199, 'Netflix',             '2026-04-10'::date, (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, (SELECT id FROM jars WHERE name='Leisure'),      false),
				('expense'::text,    620, 'Petrol',              '2026-04-14'::date, (SELECT id FROM accounts WHERE name='Cash'),             NULL, (SELECT id FROM jars WHERE name='Necessities'), false),
				('transfer'::text,  4500, 'CC Bill Payment',     '2026-04-22'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'), (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, false),
				-- income (Mar)
				('income'::text,   85000, 'March Salary',        '2026-03-01'::date, NULL, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, true),
				('expense'::text,   4200, 'Rent',                '2026-03-02'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='Necessities'), false),
				('expense'::text,   1000, 'SIP - Nifty 50',     '2026-03-05'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='SIP'),           false),
				('expense'::text,   8500, 'Chitty Instalment',   '2026-03-07'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='Chitty'),        false),
				('expense'::text,   8500, 'Charity - Akshaya',   '2026-03-13'::date, (SELECT id FROM accounts WHERE name='HDFC Savings'),    NULL, (SELECT id FROM jars WHERE name='Charity'),       false),
				('expense'::text,    560, 'Grocery',             '2026-03-08'::date, (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, (SELECT id FROM jars WHERE name='Necessities'), false),
				('expense'::text,    380, 'Swiggy',              '2026-03-12'::date, (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, (SELECT id FROM jars WHERE name='Leisure'),      false),
				('expense'::text,    199, 'Netflix',             '2026-03-10'::date, (SELECT id FROM accounts WHERE name='ICICI Credit Card'), NULL, (SELECT id FROM jars WHERE name='Leisure'),      false)
			) AS v(type, amount, name, date, from_account_id, to_account_id, jar_id, is_master_income)
		`},

		// a few templates
		{"seed templates", `
			INSERT INTO tx_templates (name, type, amount, jar_id, from_account_id, to_account_id, is_master_income)
			VALUES
				('Salary',         'income',   85000, NULL,
					NULL,
					(SELECT id FROM accounts WHERE name='HDFC Savings'),
					true),
				('Rent',           'expense',   4200,
					(SELECT id FROM jars WHERE name='Necessities'),
					(SELECT id FROM accounts WHERE name='HDFC Savings'),
					NULL, false),
				('SIP',            'expense',   1000,
					(SELECT id FROM jars WHERE name='SIP'),
					(SELECT id FROM accounts WHERE name='HDFC Savings'),
					NULL, false),
				('Netflix',        'expense',    199,
					(SELECT id FROM jars WHERE name='Leisure'),
					(SELECT id FROM accounts WHERE name='ICICI Credit Card'),
					NULL, false),
				('Zomato',         'expense',      0,
					(SELECT id FROM jars WHERE name='Leisure'),
					(SELECT id FROM accounts WHERE name='ICICI Credit Card'),
					NULL, false),
				('CC Bill Payment','transfer',     0,
					NULL,
					(SELECT id FROM accounts WHERE name='HDFC Savings'),
					(SELECT id FROM accounts WHERE name='ICICI Credit Card'),
					false)
		`},
	}

	for _, s := range steps {
		if _, err := Pool.Exec(ctx, s.sql); err != nil {
			log.Fatalf("seed step %q failed: %v", s.name, err)
		}
		log.Printf("seed ok: %s", s.name)
	}
	log.Println("demo data seeded successfully")
}

// NeedsSeeding returns true if the accounts table is empty.
// Used to auto-seed on first startup in demo mode.
func NeedsSeeding(ctx context.Context) bool {
	var count int
	Pool.QueryRow(ctx, `SELECT COUNT(*) FROM accounts`).Scan(&count)
	return count == 0
}

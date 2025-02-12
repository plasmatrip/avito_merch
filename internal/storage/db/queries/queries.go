package queries

const (
	InsertUser = `
		WITH insertIntoUsers AS (
			INSERT INTO users (id, date, login, password)
			VALUES (gen_random_uuid (), @date,@login, @password)
			RETURNING id
		)
		INSERT INTO accounts (id, user_id, amount)
		VALUES (gen_random_uuid (), (SELECT id FROM insertIntoUsers), 1000)
		RETURNING user_id;
	`

	SelectItem = `
		SELECT id, price::money::numeric FROM merch WHERE name = @item_name
	`

	SelectAccount = `
		SELECT amount::money::numeric FROM accounts WHERE user_id = @user_id
	`

	SelectUser = `
		SELECT id FROM users WHERE login = @login
	`

	BuyItem = `
		WITH insertPurchase AS (
			INSERT INTO purchases (id, date, user_id, merch_id)
			values (
				gen_random_uuid (),
				CURRENT_TIMESTAMP,
				@user_id,
				@merch_id
			)
		) UPDATE accounts
			SET amount = accounts.amount - @price
			WHERE user_id = @user_id;
	`

	UpdateCoin = `
		UPDATE accounts
		SET amount = accounts.amount + @amount
		WHERE user_id = @user_id
	`

	InsertTransaction = `
		INSERT INTO transactions (id, date, from_user_id, to_user_id, amount)
		VALUES (gen_random_uuid (), CURRENT_TIMESTAMP, @from_user_id, @to_user_id, @amount)
	`

	SelectPurchases = `
		SELECT m.name, COUNT(p.merch_id)
		FROM merch m 
		RIGHT JOIN purchases p ON m.id = p.merch_id 
		WHERE p.user_id =@user_id
		GROUP BY m."name" 
	`

	SelectRecicedCoins = `
		SELECT u.login, sum(t.amount)::money::numeric
		FROM users u 
		RIGHT JOIN transactions t ON u.id = t.from_user_id 
		WHERE t.to_user_id = @user_id
		GROUP BY u.login
	`
	SelectSentCoins = `
		SELECT u.login, sum(t.amount)::money::numeric
		FROM users u 
		RIGHT JOIN transactions t ON u.id = t.to_user_id 
		WHERE t.from_user_id = @user_id
		GROUP BY u.login
	`
	SelectTransactions = `
		SELECT u.login, sum(t.amount)::money::numeric AS amount, 'received' AS type
			FROM users u 
			RIGHT JOIN transactions t ON u.id = t.from_user_id 
			WHERE t.to_user_id = @user_id
			GROUP BY u.login
		UNION ALL
		SELECT u.login, sum(t.amount)::money::numeric AS amount, 'sent' AS type
			FROM users u 
			RIGHT JOIN transactions t ON u.id = t.to_user_id 
			WHERE t.from_user_id = @user_id
			GROUP BY u.login;
	`
)

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

	// BuyItem = `
	// 	WITH select_item AS (
	// 		SELECT merch.id, merch.price
	// 		FROM merch
	// 		INNER JOIN accounts ON accounts.user_id = @user_id
	// 		WHERE name = @item_name
	// 			AND accounts.amount::money::numeric> 0
	// 			AND accounts.amount::money::numeric  - merch.price::money::numeric > 0
	// 	), insertPurchase AS (
	// 		INSERT INTO purchases (id, date, user_id, merch_id)
	// 		values (
	// 			gen_random_uuid (),
	// 			CURRENT_TIMESTAMP,
	// 			@user_id,
	// 			(SELECT id FROM select_item)
	// 		)
	// 	) UPDATE accounts
	// 		SET amount = accounts.amount - (SELECT price FROM select_item)
	// 		WHERE user_id = @user_id;
	// `

	AddSong = `
		INSERT INTO music_library (group_name, song_name, release_date, lyrics, link)
		VALUES (@group_name, @song_name, @release_date, @lyrics, @link);
	`
	DeleteSong = `
		DELETE FROM music_library
		WHERE group_name = @group_name AND song_name = @song_name;
	`

	UpdateSong = `
		UPDATE music_library
		SET group_name = @group_name,
			song_name = @song_name,
			release_date = COALESCE(@release_date, release_date),
			lyrics = CASE WHEN TRIM(@lyrics) != '' THEN @lyrics ELSE lyrics END,
			link = CASE WHEN TRIM(@link) != '' THEN @link ELSE link END
		WHERE group_name = @group_name AND song_name = @song_name;
	`

	SelectSongs = `
		SELECT group_name, song_name, release_date, lyrics, link
		FROM music_library
		WHERE 1=1
	`

	SelectSong = `
		SELECT lyrics
		FROM music_library
		WHERE group_name = @group AND song_name = @song;
	`
)

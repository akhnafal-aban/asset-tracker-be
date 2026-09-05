package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Connected to SQLite database", slog.String("path", dbPath))

	err = createSchema(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	err = seedData(db)
	if err != nil {
		return nil, fmt.Errorf("failed to seed data: %w", err)
	}

	return db, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS asset_categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		description TEXT
	);

	CREATE TABLE IF NOT EXISTS assets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		purchase_price DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
		current_value DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
		purchase_date DATE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (category_id) REFERENCES asset_categories(id)
	);

	CREATE INDEX IF NOT EXISTS idx_assets_user_id ON assets(user_id);
	CREATE INDEX IF NOT EXISTS idx_assets_category_id ON assets(category_id);
	`

	_, err := db.Exec(schema)
	return err
}

func seedData(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO users (id, email, name) 
		VALUES (1, 'default@example.com', 'Default User')
		ON CONFLICT(id) DO NOTHING;
	`)
	if err != nil {
		return err
	}

	categories := []string{"Cash", "Physical", "Stocks", "Crypto", "Other"}
	for _, cat := range categories {
		_, err := db.Exec(`
			INSERT INTO asset_categories (name, description) 
			VALUES (?, 'Default category')
			ON CONFLICT(name) DO NOTHING;
		`, cat)
		if err != nil {
			return err
		}
	}

	var assetCount int
	err = db.QueryRow("SELECT COUNT(*) FROM assets").Scan(&assetCount)
	if err != nil {
		return err
	}

	if assetCount == 0 {
		dummyAssets := []struct {
			CategoryID    int
			Name          string
			Description   string
			PurchasePrice float64
			CurrentValue  float64
		}{
			{1, "Emergency Fund", "Savings account in Bank Mandiri", 2000000.00, 2000000.00},
			{2, "MacBook Pro", "Work laptop", 35000000.00, 30000000.00},
			{2, "Iphone 17", "Personal Phone", 17000000.00, 14000000.00},
			{2, "Tecno Camon 40 Pro 5G", "Second Phone that i want to sell this", 3800000.00, 2500000.00},
		}

		for _, a := range dummyAssets {
			_, err := db.Exec(`
				INSERT INTO assets (user_id, category_id, name, description, purchase_price, current_value, purchase_date)
				VALUES (1, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
			`, a.CategoryID, a.Name, a.Description, a.PurchasePrice, a.CurrentValue)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

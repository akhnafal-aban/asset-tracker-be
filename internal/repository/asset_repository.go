package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/akhnafal-aban/asset_tracker_be/internal/model"
)

type AssetRepository interface {
	Create(ctx context.Context, asset *model.Asset) error
	GetByID(ctx context.Context, id int64) (*model.Asset, error)
	ListByUserID(ctx context.Context, userID int64) ([]*model.Asset, error)
	Update(ctx context.Context, asset *model.Asset) error
	Delete(ctx context.Context, id int64) error
	
	ListCategories(ctx context.Context) ([]*model.AssetCategory, error)
}

type sqliteAssetRepository struct {
	db *sql.DB
}

func NewAssetRepository(db *sql.DB) AssetRepository {
	return &sqliteAssetRepository{db: db}
}

func (r *sqliteAssetRepository) Create(ctx context.Context, asset *model.Asset) error {
	query := `
		INSERT INTO assets (user_id, category_id, name, description, purchase_price, current_value, purchase_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	if asset.PurchaseDate.IsZero() {
		asset.PurchaseDate = now
	}
	
	res, err := r.db.ExecContext(ctx, query,
		asset.UserID,
		asset.CategoryID,
		asset.Name,
		asset.Description,
		asset.PurchasePrice,
		asset.CurrentValue,
		asset.PurchaseDate,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create asset: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	
	asset.ID = id
	asset.CreatedAt = now
	asset.UpdatedAt = now

	return nil
}

func (r *sqliteAssetRepository) GetByID(ctx context.Context, id int64) (*model.Asset, error) {
	query := `
		SELECT a.id, a.user_id, a.category_id, a.name, a.description, a.purchase_price, a.current_value, a.purchase_date, a.created_at, a.updated_at,
		       c.id, c.name, c.description
		FROM assets a
		LEFT JOIN asset_categories c ON a.category_id = c.id
		WHERE a.id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var asset model.Asset
	var category model.AssetCategory
	var desc sql.NullString
	var catDesc sql.NullString
	var purDate sql.NullTime

	err := row.Scan(
		&asset.ID, &asset.UserID, &asset.CategoryID, &asset.Name, &desc, &asset.PurchasePrice, &asset.CurrentValue, &purDate, &asset.CreatedAt, &asset.UpdatedAt,
		&category.ID, &category.Name, &catDesc,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("asset not found")
		}
		return nil, fmt.Errorf("failed to scan asset: %w", err)
	}

	if desc.Valid {
		asset.Description = desc.String
	}
	if purDate.Valid {
		asset.PurchaseDate = purDate.Time
	}
	if catDesc.Valid {
		category.Description = catDesc.String
	}
	
	asset.Category = &category

	return &asset, nil
}

func (r *sqliteAssetRepository) ListByUserID(ctx context.Context, userID int64) ([]*model.Asset, error) {
	query := `
		SELECT a.id, a.user_id, a.category_id, a.name, a.description, a.purchase_price, a.current_value, a.purchase_date, a.created_at, a.updated_at,
		       c.id, c.name, c.description
		FROM assets a
		LEFT JOIN asset_categories c ON a.category_id = c.id
		WHERE a.user_id = ?
		ORDER BY a.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	var assets []*model.Asset
	for rows.Next() {
		var asset model.Asset
		var category model.AssetCategory
		var desc sql.NullString
		var catDesc sql.NullString
		var purDate sql.NullTime

		err := rows.Scan(
			&asset.ID, &asset.UserID, &asset.CategoryID, &asset.Name, &desc, &asset.PurchasePrice, &asset.CurrentValue, &purDate, &asset.CreatedAt, &asset.UpdatedAt,
			&category.ID, &category.Name, &catDesc,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset row: %w", err)
		}

		if desc.Valid {
			asset.Description = desc.String
		}
		if purDate.Valid {
			asset.PurchaseDate = purDate.Time
		}
		if catDesc.Valid {
			category.Description = catDesc.String
		}
		asset.Category = &category

		assets = append(assets, &asset)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return assets, nil
}

func (r *sqliteAssetRepository) Update(ctx context.Context, asset *model.Asset) error {
	query := `
		UPDATE assets
		SET category_id = ?, name = ?, description = ?, purchase_price = ?, current_value = ?, purchase_date = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`
	now := time.Now()
	
	res, err := r.db.ExecContext(ctx, query,
		asset.CategoryID,
		asset.Name,
		asset.Description,
		asset.PurchasePrice,
		asset.CurrentValue,
		asset.PurchaseDate,
		now,
		asset.ID,
		asset.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update asset: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("asset not found or not owned by user")
	}

	asset.UpdatedAt = now
	return nil
}

func (r *sqliteAssetRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM assets WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("asset not found")
	}

	return nil
}

func (r *sqliteAssetRepository) ListCategories(ctx context.Context) ([]*model.AssetCategory, error) {
	query := `SELECT id, name, description FROM asset_categories ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []*model.AssetCategory
	for rows.Next() {
		var cat model.AssetCategory
		var desc sql.NullString
		
		if err := rows.Scan(&cat.ID, &cat.Name, &desc); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		if desc.Valid {
			cat.Description = desc.String
		}
		categories = append(categories, &cat)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("category row iteration error: %w", err)
	}

	return categories, nil
}

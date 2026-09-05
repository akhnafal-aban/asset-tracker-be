package service

import (
	"context"
	"fmt"
	"time"

	"github.com/akhnafal-aban/asset_tracker_be/internal/model"
	"github.com/akhnafal-aban/asset_tracker_be/internal/repository"
)

// AssetService defines the interface for business logic related to assets.
type AssetService interface {
	CreateAsset(ctx context.Context, userID int64, req CreateAssetRequest) (*model.Asset, error)
	GetAsset(ctx context.Context, userID int64, assetID int64) (*model.Asset, error)
	ListAssets(ctx context.Context, userID int64) ([]*model.Asset, error)
	UpdateAsset(ctx context.Context, userID int64, assetID int64, req UpdateAssetRequest) (*model.Asset, error)
	DeleteAsset(ctx context.Context, userID int64, assetID int64) error
	
	ListCategories(ctx context.Context) ([]*model.AssetCategory, error)
}

// CreateAssetRequest represents the payload for creating an asset.
type CreateAssetRequest struct {
	CategoryID    int64     `json:"category_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	PurchasePrice float64   `json:"purchase_price"`
	CurrentValue  float64   `json:"current_value"`
	PurchaseDate  time.Time `json:"purchase_date"`
}

// UpdateAssetRequest represents the payload for updating an asset.
type UpdateAssetRequest struct {
	CategoryID    *int64     `json:"category_id"`
	Name          *string    `json:"name"`
	Description   *string    `json:"description"`
	PurchasePrice *float64   `json:"purchase_price"`
	CurrentValue  *float64   `json:"current_value"`
	PurchaseDate  *time.Time `json:"purchase_date"`
}

type assetService struct {
	repo repository.AssetRepository
}

// NewAssetService creates a new AssetService.
func NewAssetService(repo repository.AssetRepository) AssetService {
	return &assetService{repo: repo}
}

func (s *assetService) CreateAsset(ctx context.Context, userID int64, req CreateAssetRequest) (*model.Asset, error) {
	// Basic validation
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.CategoryID <= 0 {
		return nil, fmt.Errorf("valid category_id is required")
	}
	if req.PurchasePrice < 0 || req.CurrentValue < 0 {
		return nil, fmt.Errorf("price and value cannot be negative")
	}

	asset := &model.Asset{
		UserID:        userID,
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		Description:   req.Description,
		PurchasePrice: req.PurchasePrice,
		CurrentValue:  req.CurrentValue,
		PurchaseDate:  req.PurchaseDate,
	}

	err := s.repo.Create(ctx, asset)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	return s.GetAsset(ctx, userID, asset.ID) // Fetch with joins (category info)
}

func (s *assetService) GetAsset(ctx context.Context, userID int64, assetID int64) (*model.Asset, error) {
	asset, err := s.repo.GetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}
	// Verify ownership
	if asset.UserID != userID {
		return nil, fmt.Errorf("unauthorized: asset does not belong to user")
	}
	return asset, nil
}

func (s *assetService) ListAssets(ctx context.Context, userID int64) ([]*model.Asset, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *assetService) UpdateAsset(ctx context.Context, userID int64, assetID int64, req UpdateAssetRequest) (*model.Asset, error) {
	asset, err := s.GetAsset(ctx, userID, assetID)
	if err != nil {
		return nil, err // Handles ownership check too
	}

	// Apply updates conditionally
	if req.CategoryID != nil {
		asset.CategoryID = *req.CategoryID
	}
	if req.Name != nil && *req.Name != "" {
		asset.Name = *req.Name
	}
	if req.Description != nil {
		asset.Description = *req.Description
	}
	if req.PurchasePrice != nil {
		if *req.PurchasePrice < 0 {
			return nil, fmt.Errorf("purchase price cannot be negative")
		}
		asset.PurchasePrice = *req.PurchasePrice
	}
	if req.CurrentValue != nil {
		if *req.CurrentValue < 0 {
			return nil, fmt.Errorf("current value cannot be negative")
		}
		asset.CurrentValue = *req.CurrentValue
	}
	if req.PurchaseDate != nil {
		asset.PurchaseDate = *req.PurchaseDate
	}

	err = s.repo.Update(ctx, asset)
	if err != nil {
		return nil, err
	}

	return s.GetAsset(ctx, userID, assetID) // Fetch updated state with category
}

func (s *assetService) DeleteAsset(ctx context.Context, userID int64, assetID int64) error {
	// Verify existence and ownership first
	_, err := s.GetAsset(ctx, userID, assetID)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, assetID)
}

func (s *assetService) ListCategories(ctx context.Context) ([]*model.AssetCategory, error) {
	return s.repo.ListCategories(ctx)
}

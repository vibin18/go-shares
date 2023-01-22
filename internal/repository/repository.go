package repository

import (
	"github.com/vibin18/go-shares/internal/models"
)

type DatabaseRepo interface {
	InsertNewShare(res models.Share) error
	GetAllShares() ([]models.Share, error)
	BuyShare(res models.SellBuyShare) error
	SellShare(res models.SellBuyShare) error
	GetAllSharesWithData() ([]models.TotalShare, error)
	GetAllPurchases() ([]models.SellBuyShare, error)
	GetAllSales() ([]models.SellBuyShare, error)
	GetAllSalesReport() ([]models.ShareReport, error)
	GetAllPurchaseReport() ([]models.ShareReport, error)
}

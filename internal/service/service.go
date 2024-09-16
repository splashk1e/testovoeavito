package service

import (
	"testovoe.com/internal/models"
	"testovoe.com/internal/repository"
)

type Tenders interface {
	GetTenders(serviceTypes []string) ([]models.Tender, error)
	CreateTender(item models.Tender, authorId string) (*models.Tender, error)
	GetTendersByUsername(organizationId string) ([]models.Tender, error)
	GetTenderById(id string, organizationId string) (*models.Tender, error)
	GetTenderStatusById(id string, username string) (string, error)
	EditTenderById(id string, changes map[string]string, organizationId string) (*models.Tender, error)
	CheckTenderByOrganizationId(tender *models.Tender, organizationId string) error
	EditTenderStatusById(id string, username string, status string) (*models.Tender, error)
	TenderRollBack(id string, version int, username string) (*models.Tender, error)
}
type Bids interface {
	CreateBid(item models.Bid, username string) (*models.Bid, error)
	GetBidById(id, username string) (*models.Bid, error)
	CheckPermissionForBid(id string, userId string) error
	GetBidsByUsername(username string) ([]models.Bid, error)
	GetBidByTenderId(tenderId string, username string) ([]models.Bid, error)
	EditBidById(id string, changes map[string]string, username string) (*models.Bid, error)

	ChangeBidStatus(id string, status string, username string) (*models.Bid, error)
	BidRollBack(id string, version int, username string) (*models.Bid, error)
}

type Service struct {
	Tenders
	Bids
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Tenders: NewTendersService(repos.Tenders),
		Bids:    NewBidsService(repos.Bids),
	}
}

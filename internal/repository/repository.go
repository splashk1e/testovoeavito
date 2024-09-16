package repository

import (
	"testovoe.com/bootstrap"
	"testovoe.com/internal/models"
)

type Users interface {
	CheckUserOrganization(userId string) (string, error)
	GetIdByUsername(username string) (string, error)
}
type Organization interface {
	GetUsersFromOrganization(organizationId string) ([]string, error)
}
type Tenders interface {
	GetTendersWithServiceType(serviceType string) ([]models.Tender, error)
	GetTenders() ([]models.Tender, error)
	CreateTender(item models.Tender) (string, error)
	GetTendersByOrganizationId(organizationId string) ([]models.Tender, error)
	GetTenderById(id string) (*models.Tender, error)
	EditTenderById(id string, changes map[string]string) error
	ChangeTenderStatus(id string, status string) error
	TenderRollBack(id string, version int) error
	Users
}
type Bids interface {
	CreateBid(item models.Bid) (string, error)
	GetBidsByUserId(AuthorId string) ([]models.Bid, error)
	GetBidById(id string) (*models.Bid, error)
	GetBidByTenderId(tenderId string) ([]models.Bid, error)
	EditBidById(id string, changes map[string]string) error
	ChangeBidStatus(id string, status string) error
	BidRollBack(id string, version int) error
	Users
	Organization
}

type Repository struct {
	Tenders
	Bids
}

func NewRepository(repos *bootstrap.PostgresClient) *Repository {
	return &Repository{
		Tenders: NewTendersPostgres(repos),
		Bids:    NewBidsPostgres(repos),
	}
}

package service

import (
	"errors"
	"slices"

	"testovoe.com/internal/models"
	"testovoe.com/internal/repository"
)

type BidsService struct {
	Repository repository.Bids
}

func NewBidsService(repo repository.Bids) *BidsService {
	return &BidsService{Repository: repo}
}
func (service *BidsService) CreateBid(item models.Bid, username string) (*models.Bid, error) {
	authorId, err := service.Repository.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}
	if _, err = service.Repository.CheckUserOrganization(authorId); err != nil {
		return nil, err
	}
	id, err := service.Repository.CreateBid(item)
	if err != nil {
		return nil, err
	}
	return service.Repository.GetBidById(id)
}
func (service *BidsService) GetBidById(id, username string) (*models.Bid, error) {
	userId, err := service.Repository.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}
	bid, err := service.Repository.GetBidById(id)
	if err != nil {
		return nil, err
	}
	if err = service.CheckPermissionForBid(id, userId); err != nil {
		return nil, err
	}
	return bid, nil
}
func (service *BidsService) CheckPermissionForBid(id string, userId string) error {
	bid, err := service.Repository.GetBidById(id)
	if err != nil {
		return err
	}
	if bid.AuthorType == `User` {
		organizationId, err := service.Repository.CheckUserOrganization(bid.AuthorId)
		if err != nil {
			return err
		}
		users, err := service.Repository.GetUsersFromOrganization(organizationId)
		if err != nil {
			return err
		}
		if !slices.Contains(users, userId) {
			return errors.New("user have no permission")
		}
		return nil
	}
	if bid.AuthorType == "Organization" {
		users, err := service.Repository.GetUsersFromOrganization(bid.AuthorId)
		if err != nil {
			return err
		}
		if !slices.Contains(users, userId) {
			return errors.New("user have no permission")
		}
		return nil
	}
	return errors.New("wrong bid authortype")
}

func (service *BidsService) GetBidsByUsername(username string) ([]models.Bid, error) {
	userId, err := service.Repository.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}
	return service.Repository.GetBidsByUserId(userId)

}

func (service *BidsService) GetBidByTenderId(tenderId string, username string) ([]models.Bid, error) {
	var output []models.Bid
	userId, err := service.Repository.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}
	bids, err := service.Repository.GetBidByTenderId(tenderId)
	if err != nil {
		return nil, err
	}
	for _, bid := range bids {
		if err = service.CheckPermissionForBid(bid.Id, userId); err == nil {
			output = append(output, bid)
		}
	}
	if len(output) == 0 {
		return nil, errors.New("user have no permission")
	}
	return output, nil
}

func (service *BidsService) EditBidById(id string, changes map[string]string, username string) (*models.Bid, error) {
	userId, err := service.Repository.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}
	err = service.CheckPermissionForBid(id, userId)
	if err != nil {
		return nil, err
	}
	bid, err := service.GetBidById(id, username)
	if err != nil {
		return nil, err
	}
	bidmap := map[string]string{
		"name":        bid.Name,
		"description": bid.Description,
	}
	for key, val := range changes {
		if val == "" {
			changes[key] = bidmap[key]
		}
	}
	if err = service.Repository.EditBidById(id, changes); err != nil {
		return nil, err
	}
	return service.GetBidById(id, username)
}
func (service *BidsService) ChangeBidStatus(id string, status string, username string) (*models.Bid, error) {
	userId, err := service.Repository.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}
	err = service.CheckPermissionForBid(id, userId)
	if err != nil {
		return nil, err
	}
	if err = service.Repository.ChangeBidStatus(id, status); err != nil {
		return nil, err
	}
	return service.GetBidById(id, username)
}
func (service *BidsService) BidRollBack(id string, version int, username string) (*models.Bid, error) {
	userId, err := service.Repository.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}
	err = service.CheckPermissionForBid(id, userId)
	if err != nil {
		return nil, err
	}
	if err = service.Repository.BidRollBack(id, version); err != nil {
		return nil, err
	}
	return service.Repository.GetBidById(id)
}

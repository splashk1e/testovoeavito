package service

import (
	"errors"

	"testovoe.com/internal/models"
	"testovoe.com/internal/repository"
)

type TendersService struct {
	Repository repository.Tenders
}

func NewTendersService(repo repository.Tenders) *TendersService {
	return &TendersService{Repository: repo}
}
func (service *TendersService) GetTenders(serviceTypes []string) ([]models.Tender, error) {
	var tenders []models.Tender
	if serviceTypes != nil {
		for _, serviceType := range serviceTypes {
			response, err := service.Repository.GetTendersWithServiceType(serviceType)
			if err != nil {
				return nil, err
			}
			tenders = append(tenders, response...)
		}
		return tenders, nil
	}
	return service.Repository.GetTenders()
}

func (service *TendersService) CreateTender(item models.Tender, username string) (*models.Tender, error) {
	authorId, err := service.Repository.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}
	if _, err = service.Repository.CheckUserOrganization(authorId); err != nil {
		return nil, err
	}
	id, err := service.Repository.CreateTender(item)
	if err != nil {
		return nil, err
	}
	return service.Repository.GetTenderById(id)
}
func (service *TendersService) GetTendersByUsername(username string) ([]models.Tender, error) {
	authorId, err := service.Repository.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}
	organizationId, err := service.Repository.CheckUserOrganization(authorId)
	if err != nil {
		return nil, err
	}
	return service.Repository.GetTendersByOrganizationId(organizationId)
}
func (service *TendersService) GetTenderById(id string, username string) (*models.Tender, error) {
	authorId, err := service.Repository.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}
	organizationId, err := service.Repository.CheckUserOrganization(authorId)
	if err != nil {
		return nil, err
	}
	tender, err := service.Repository.GetTenderById(id)
	if err != nil {
		return nil, err
	}
	if err = service.CheckTenderByOrganizationId(tender, organizationId); err != nil {
		return nil, err
	}
	return tender, nil
}
func (service *TendersService) GetTenderStatusById(id string, username string) (string, error) {
	tender, err := service.GetTenderById(id, username)
	if err != nil {
		return "", err
	}
	return tender.Status, nil
}
func (service *TendersService) EditTenderStatusById(id string, username string, status string) (*models.Tender, error) {
	_, err := service.GetTenderById(id, username)
	if err != nil {
		return nil, err
	}
	if err := service.Repository.ChangeTenderStatus(id, status); err != nil {
		return nil, err
	}
	return service.GetTenderById(id, username)
}
func (service *TendersService) EditTenderById(id string, changes map[string]string, organizationId string) (*models.Tender, error) {
	tender, err := service.GetTenderById(id, organizationId)
	tendermap := map[string]string{
		"name":         tender.Name,
		"decription":   tender.Description,
		"service_type": tender.ServiceType,
	}
	if err != nil {
		return nil, err
	}
	for key, val := range changes {
		if val == "" {
			changes[key] = tendermap[key]
		}
	}
	if err := service.Repository.EditTenderById(id, changes); err != nil {
		return nil, err
	}
	return service.Repository.GetTenderById(id)
}
func (service *TendersService) CheckTenderByOrganizationId(tender *models.Tender, organizationId string) error {
	if tender.OrganizationId == organizationId {
		return nil
	}
	return errors.New("user have no permission")
}
func (service *TendersService) TenderRollBack(id string, version int, username string) (*models.Tender, error) {
	_, err := service.GetTenderById(id, username)
	if err != nil {
		return nil, err
	}
	if err = service.Repository.TenderRollBack(id, version); err != nil {
		return nil, err
	}
	return service.Repository.GetTenderById(id)
}

package repository

import (
	"context"
	"errors"

	"testovoe.com/bootstrap"
	"testovoe.com/internal/models"
)

type TendersPostgres struct {
	db *bootstrap.PostgresClient
	UsersPostgres
}

func NewTendersPostgres(db *bootstrap.PostgresClient) *TendersPostgres {
	return &TendersPostgres{db: db}
}
func (repo *TendersPostgres) CreateTender(item models.Tender) (string, error) {
	repo.db.Mu.RLock()
	defer repo.db.Mu.RUnlock()
	var id string
	row, err := repo.db.Pool.Query(context.Background(),
		"INSERT INTO tender (name, description, service_type, organization_id) VALUES ($1, $2, $3, $4, $5,$6) RETURNING id;",
		item.Name, item.Description, item.ServiceType, item.OrganizationId)
	if err != nil {
		return "", err
	}
	defer row.Conn()
	if row.Next() {
		if err := row.Scan(&id); err != nil {
			return "", err
		}
		return id, nil
	}
	return "", errors.New("psql don't return id")

}

func (repo *TendersPostgres) GetTendersWithServiceType(serviceTypes string) ([]models.Tender, error) {

	repo.db.Mu.RLock()
	defer repo.db.Mu.RUnlock()
	rows, err := repo.db.Pool.Query(context.Background(),
		"SELECT id, name, description, service_type, status, organization_id, version, created_at FROM tender WHERE service_type=$1;",
		serviceTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Conn()
	var data []models.Tender
	for rows.Next() {
		var item models.Tender
		err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.ServiceType,
			&item.Status,
			&item.OrganizationId,
			&item.Version,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return data, nil
}
func (repo *TendersPostgres) GetTenders() ([]models.Tender, error) {

	repo.db.Mu.RLock()
	defer repo.db.Mu.RUnlock()
	rows, err := repo.db.Pool.Query(context.Background(),
		"SELECT id, name, description, service_type, status, organization_id, version, created_at FROM tender;")
	if err != nil {
		return nil, err
	}
	defer rows.Conn()
	var data []models.Tender
	for rows.Next() {
		var item models.Tender
		err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.ServiceType,
			&item.Status,
			&item.OrganizationId,
			&item.Version,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return data, nil
}
func (repo *TendersPostgres) GetTendersByOrganizationId(organizationId string) ([]models.Tender, error) {
	repo.db.Mu.RLock()
	defer repo.db.Mu.RUnlock()
	rows, err := repo.db.Pool.Query(context.Background(),
		"SELECT id, name, description, service_type, status, organization_id, version, created_at FROM tender WHERE organization_id=$1;",
		organizationId)
	if err != nil {
		return nil, err
	}
	defer rows.Conn()
	var data []models.Tender
	for rows.Next() {
		var item models.Tender
		err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.ServiceType,
			&item.Status,
			&item.OrganizationId,
			&item.Version,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return data, nil
}
func (repo *TendersPostgres) GetTenderById(id string) (*models.Tender, error) {
	repo.db.Mu.RLock()
	defer repo.db.Mu.RUnlock()
	row, err := repo.db.Pool.Query(context.Background(),
		"SELECT id, name, description, service_type, status, organization_id, version, created_at FROM tender WHERE id=$1;",
		id)
	if err != nil {
		return nil, err
	}
	defer row.Conn()
	var item models.Tender
	if row.Next() {

		err := row.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.ServiceType,
			&item.Status,
			&item.OrganizationId,
			&item.Version,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	return &item, nil
}

func (repo *TendersPostgres) EditTenderById(id string, changes map[string]string) error {
	repo.db.Mu.Lock()
	defer repo.db.Mu.Unlock()
	_, err := repo.db.Pool.Exec(context.Background(),
		"UPDATE tender SET name=$1, description=$2, service_type=$3, WHERE id=$4",
		changes["name"], changes["description"], changes["service_type"], id)
	if err != nil {
		return err
	}
	return nil
}
func (repo *TendersPostgres) ChangeTenderStatus(id string, status string) error {
	repo.db.Mu.Lock()
	defer repo.db.Mu.Unlock()
	_, err := repo.db.Pool.Exec(context.Background(),
		"UPDATE tender SET status=$1 WHERE id=$2",
		status, id)
	if err != nil {
		return err
	}
	return nil
}
func (repo *TendersPostgres) TenderRollBack(id string, version int) error {
	repo.db.Mu.Lock()
	defer repo.db.Mu.Unlock()
	_, err := repo.db.Pool.Exec(context.Background(), `INSERT INTO tender (id, name, description, service_type, version)
	SELECT tender_id, name, description, service_type, version
	FROM tender_history
	WHERE tender_id = $1 AND version = %2
	ORDER BY updated_at DESC
	LIMIT 1;`, id, version)
	return err
}

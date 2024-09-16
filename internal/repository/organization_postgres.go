package repository

import (
	"context"

	"testovoe.com/bootstrap"
)

type OrganizationPostgres struct {
	db *bootstrap.PostgresClient
}

func (repo *OrganizationPostgres) GetUsersFromOrganization(organizationId string) ([]string, error) {
	repo.db.Mu.RLock()
	defer repo.db.Mu.RUnlock()
	var data []string
	var item string
	rows, err := repo.db.Pool.Query(context.Background(),
		"SELECT user_id FROM organization_responsible WHERE organization_id=$1",
		organizationId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&item)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return data, nil
}

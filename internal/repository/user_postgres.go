package repository

import (
	"context"
	"errors"

	"testovoe.com/bootstrap"
)

type UsersPostgres struct {
	db *bootstrap.PostgresClient
}

func (repo *UsersPostgres) GetIdByUsername(username string) (string, error) {
	repo.db.Mu.RLock()
	defer repo.db.Mu.RUnlock()
	row, err := repo.db.Pool.Query(context.Background(),
		"SELECT id FROM employee WHERE username=$1;",
		username)
	if err != nil {
		return "", err
	}
	defer row.Conn()
	var userId string
	if row.Next() {
		err := row.Scan(
			&userId,
		)
		if err != nil {
			return "", err
		}
		return userId, nil
	}
	return "", errors.New("user not found")
}
func (repo *UsersPostgres) CheckUserOrganization(userId string) (string, error) {
	repo.db.Mu.RLock()
	defer repo.db.Mu.RUnlock()
	row, err := repo.db.Pool.Query(context.Background(),
		"SELECT organization_id FROM organization_responsible WHERE user_id=$1;",
		userId)
	if err != nil {
		return "", err
	}
	defer row.Conn()
	var organizationId string
	if row.Next() {
		err := row.Scan(
			&organizationId,
		)
		if err != nil {
			return "", err
		}
		return organizationId, nil
	}
	return "", errors.New("user have no organization")
}

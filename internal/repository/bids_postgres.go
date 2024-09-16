package repository

import (
	"context"
	"errors"

	"testovoe.com/bootstrap"
	"testovoe.com/internal/models"
)

type BidsPostgres struct {
	db *bootstrap.PostgresClient
	UsersPostgres
	OrganizationPostgres
}

func NewBidsPostgres(db *bootstrap.PostgresClient) *BidsPostgres {
	return &BidsPostgres{db: db}
}
func (repo *BidsPostgres) CreateBid(item models.Bid) (string, error) {
	repo.db.Mu.RLock()
	defer repo.db.Mu.RUnlock()
	var id string
	row, err := repo.db.Pool.Query(context.Background(),
		"INSERT INTO bid (name, decription, status, tender_id, author_type, author_id) VALUES ($1, $2, $3, $4, $5,$6) RETURNING id;",
		item.Name, item.Description, item.Status, item.TenderId, item.AuthorType, item.AuthorId)
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
	return "", errors.New("postgres don't return id")

}

func (repo *BidsPostgres) GetBidsByUserId(userId string) ([]models.Bid, error) {

	repo.db.Mu.RLock()
	defer repo.db.Mu.RUnlock()
	rows, err := repo.db.Pool.Query(context.Background(),
		`SELECT id, name, decription, status, tender_id, author_type, author_id, version, created_at FROM bid WHERE author_id=$1 AND author_type="Organization";`,
		userId)
	if err != nil {
		return nil, err
	}
	defer rows.Conn()
	var data []models.Bid
	for rows.Next() {
		var item models.Bid
		err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.Status,
			&item.TenderId,
			&item.AuthorType,
			&item.AuthorId,
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

func (repo *BidsPostgres) GetBidById(id string) (*models.Bid, error) {
	repo.db.Mu.RLock()
	defer repo.db.Mu.RUnlock()
	row, err := repo.db.Pool.Query(context.Background(),
		"SELECT id, name, decription, status, tender_id, author_type, author_id, version, created_at FROM bid WHERE id=$1;",
		id)
	if err != nil {
		return nil, err
	}
	defer row.Conn()
	var item models.Bid
	if row.Next() {

		err := row.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.Status,
			&item.TenderId,
			&item.AuthorType,
			&item.AuthorId,
			&item.Version,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	return &item, nil
}

func (repo *BidsPostgres) GetBidByTenderId(tenderId string) ([]models.Bid, error) {
	repo.db.Mu.RLock()
	defer repo.db.Mu.RUnlock()
	rows, err := repo.db.Pool.Query(context.Background(),
		"SELECT id, name, decription, status, tender_id, author_type, author_id, version, created_at FROM bid WHERE tender_id=$1;",
		tenderId)
	if err != nil {
		return nil, err
	}
	defer rows.Conn()
	var item models.Bid
	var data []models.Bid
	for rows.Next() {
		err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.Status,
			&item.TenderId,
			&item.AuthorType,
			&item.AuthorId,
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

func (repo *BidsPostgres) EditBidById(id string, changes map[string]string) error {
	repo.db.Mu.Lock()
	defer repo.db.Mu.Unlock()
	_, err := repo.db.Pool.Exec(context.Background(),
		"UPDATE bid SET name=$1, description=$2, WHERE id=$3",
		changes["name"], changes["description"], id)
	if err != nil {
		return err
	}
	return nil
}
func (repo *BidsPostgres) ChangeBidStatus(id string, status string) error {
	repo.db.Mu.Lock()
	defer repo.db.Mu.Unlock()
	_, err := repo.db.Pool.Exec(context.Background(),
		"UPDATE bid SET status=$1 WHERE id=$2",
		status, id)
	if err != nil {
		return err
	}
	return nil
}
func (repo *BidsPostgres) BidRollBack(id string, version int) error {
	repo.db.Mu.Lock()
	defer repo.db.Mu.Unlock()
	_, err := repo.db.Pool.Exec(context.Background(), `INSERT INTO bid (id, name, description, version)
	SELECT bid_id, name, description, version
	FROM bid_history
	WHERE bid_id = $1 AND version = $2
	ORDER BY updated_at DESC
	LIMIT 1;`, id, version)
	return err
}

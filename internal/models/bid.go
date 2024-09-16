package models

import "time"

type Bid struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	TenderId    string    `json:"tenderId"`
	Version     int       `json:"version"`
	AuthorType  string    `json:"authorType"`
	AuthorId    string    `json:"authorId"`
	CreatedAt   time.Time `json:"createdAt"`
}

package repo

import "database/sql"

type RecommendationRepo struct {
	DB *sql.DB
}

func NewRecommendationRepo(db *sql.DB) *RecommendationRepo {
	return &RecommendationRepo{DB: db}
}


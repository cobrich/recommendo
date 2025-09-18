package service

import "github.com/cobrich/recommendo/repo"

type RecommendationService struct {
	r *repo.RecommendationRepo
}

func NewRecommmendationService(r *repo.RecommendationRepo) *RecommendationService {
	return &RecommendationService{r: r}
}
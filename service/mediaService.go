package service

import "github.com/cobrich/recommendo/repo"

type MediaService struct{
	r *repo.MediaRepo
}

func NewMediaService(r *repo.MediaRepo) *MediaService {
	return &MediaService{r: r}
}


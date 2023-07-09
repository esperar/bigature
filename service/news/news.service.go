package news

import (
	"bitature/repository"
	"bitature/repository/news"
)

type NewsService struct {
	Repository *news.NewsRepository
}

var Service NewsService

func (s *NewsService) InitService() error {
	db, err := repository.OpenWithMemory()

	if err != nil {
		return err
	}

	s.Repository = &news.Repository
	s.Repository.AssignDB(db)

	return nil
}

func (s *NewsService) GetAllNews() (*[]news.NewsRaw, error) {
	raws, err := s.Repository.GetAllNews()

	return raws, err
}

func (s *NewsService) GetOneNews(id string) (*news.NewsRaw, error) {
	raw, err := s.Repository.GetOneNews(id)

	return raw, err
}

func (s *NewsService) CreateNews(n news.NewsDto) error {
	_, err := s.Repository.InsertNews(n)

	return err
}

func (s *NewsService) UpdateNews(id string, n news.NewsDto) error {
	_, err := s.Repository.UpdateOneNews(id, n)

	return err
}

func (s *NewsService) DeleteNews(id string) error {
	_, err := s.Repository.DeleteOneNews(id)

	return err
}

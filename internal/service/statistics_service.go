package service

import (
	"context"

	"frontdev333/gym/internal/domain"
	"frontdev333/gym/internal/repository"
)

type StatisticsService struct {
	statisticsRepo repository.StatisticsRepository
}

func NewStatisticsService(statisticsRepo repository.StatisticsRepository) *StatisticsService {
	return &StatisticsService{statisticsRepo: statisticsRepo}
}

func (s *StatisticsService) GetUserStatistics(ctx context.Context, userID string) (*domain.Statistics, error) {
	return s.statisticsRepo.GetUserStatistics(ctx, userID)
}

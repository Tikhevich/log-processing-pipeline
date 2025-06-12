package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"ui_api.go/internal/repositories"
)

type StatsService struct {
	repo  repositories.LogRepository
	redis *redis.Client
}

type TrafficStats struct {
	TotalRequest int64 `json:"total_request"`
	UniqueIps    int64 `json:"unique_ips"`
}

type LatencyStats struct {
	Avg float64 `json:"avg_latency"`
	Max int64   `json:"max_latency"`
}

type StatsStruct struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func NewStatsService(repo repositories.LogRepository, redis *redis.Client) *StatsService {
	return &StatsService{repo: repo, redis: redis}
}

func (s *StatsService) GetStatsByStatsTypeAndRange(ctx context.Context, statsType, rangeType string) (StatsStruct, error) {
	var statsStruct StatsStruct

	cacheKey := fmt.Sprintf("stats:%v:%v", statsType, rangeType)
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		err := json.Unmarshal([]byte(cachedData), &statsStruct)
		if err == nil {
			return statsStruct, nil
		}
	}

	switch statsType {
	case "errors":
		statsStruct.Type = statsType
		stats, _ := s.getErrorStats(ctx, rangeType)
		statsStruct.Data = stats

	case "traffic":
		statsStruct.Type = statsType
		stats, _ := s.getTrafficStats(ctx, rangeType)
		statsStruct.Data = stats

	case "latency":
		statsStruct.Type = statsType
		stats, _ := s.getLatencyStats(ctx, rangeType)
		statsStruct.Data = stats
	default:
		panic("Invalid statsType")
	}

	s.saveToCache(ctx, statsStruct, cacheKey, rangeType)

	return statsStruct, nil
}

func (s *StatsService) getErrorStats(ctx context.Context, rangeType string) (map[int]int64, error) {
	from, to := getTimeRange(rangeType)

	data, err := s.repo.GetErrorStats(ctx, from, to)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *StatsService) getTrafficStats(ctx context.Context, rangeType string) (TrafficStats, error) {
	from, to := getTimeRange(rangeType)
	var trafficStats TrafficStats

	total, unique, err := s.repo.GetTrafficStats(ctx, from, to)
	if err != nil {
		return trafficStats, err
	}

	trafficStats.TotalRequest = total
	trafficStats.UniqueIps = unique

	return trafficStats, nil
}

func (s *StatsService) getLatencyStats(ctx context.Context, rangeType string) (LatencyStats, error) {
	from, to := getTimeRange(rangeType)
	var latencyStats LatencyStats

	avg, max, err := s.repo.GetLatencyStats(ctx, from, to)

	if err != nil {
		return latencyStats, err
	}

	latencyStats.Avg = avg
	latencyStats.Max = max

	return latencyStats, nil
}

func getTimeRange(rangeType string) (from, to time.Time) {
	now := time.Now().UTC()
	switch rangeType {
	case "hour":
		return now.Add(-1 * time.Hour), now
	case "day":
		return now.Add(-24 * time.Hour), now
	case "week":
		return now.Add(-7 * 24 * time.Hour), now
	default:
		return now.Add(-1 * time.Hour), now
	}
}

func getTTLByRange(rangeType string) time.Duration {
	switch rangeType {
	case "hour":
		return 2 * time.Minute
	case "day":
		return 10 * time.Minute
	case "week":
		return 30 * time.Minute
	default:
		return 2 * time.Minute
	}
}

func (s *StatsService) saveToCache(ctx context.Context, statsData StatsStruct, key string, rangeType string) {
	data, err := json.Marshal(statsData)
	if err == nil {
		ttl := getTTLByRange(rangeType)
		_ = s.redis.Set(ctx, key, data, ttl).Err()
	}
}

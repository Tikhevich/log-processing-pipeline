package service

import (
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

func (s *StatsService) GetStatsByStatsTypeAndRange(statsType, rangeType string) (StatsStruct, error) {
	var statsStruct StatsStruct

	// cacheKey := fmt.Sprintf("stats:%v:%v", statsType, rangeType)

	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr: "localhost:6379",
	// 	DB:   0,
	// })

	// cachedData, err := redisClient.Get(ctx, cacheKey).Result()
	// if err == nil {
	// 	var response any
	// 	if err := json.Unmarshal([]byte(cachedData), &response); err == nil {
	// 		c.JSON(http.StatusOK, response)
	// 		return
	// 	}
	// }

	switch statsType {
	case "errors":
		statsStruct.Type = statsType
		stats, _ := s.GetErrorStats(rangeType)
		statsStruct.Data = stats

	case "traffic":
		statsStruct.Type = statsType
		stats, _ := s.GetTrafficStats(rangeType)
		statsStruct.Data = stats

	case "latency":
		statsStruct.Type = statsType
		stats, _ := s.GetLatencyStats(rangeType)
		statsStruct.Data = stats
	default:
		panic("Invalid statsType")
	}

	return statsStruct, nil
}

func (s *StatsService) GetErrorStats(rangeType string) (map[int]int64, error) {
	from, to := getTimeRange(rangeType)

	data, err := s.repo.GetErrorStats(from, to)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *StatsService) GetTrafficStats(rangeType string) (TrafficStats, error) {
	from, to := getTimeRange(rangeType)
	var trafficStats TrafficStats

	total, unique, err := s.repo.GetTrafficStats(from, to)
	if err != nil {
		return trafficStats, err
	}

	trafficStats.TotalRequest = total
	trafficStats.UniqueIps = unique

	return trafficStats, nil
}

func (s *StatsService) GetLatencyStats(rangeType string) (LatencyStats, error) {
	from, to := getTimeRange(rangeType)
	var latencyStats LatencyStats

	avg, max, err := s.repo.GetLatencyStats(from, to)

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

// func getTTLByRange(rangeType string) time.Duration {
// 	switch rangeType {
// 	case "hour":
// 		return 2 * time.Minute
// 	case "day":
// 		return 10 * time.Minute
// 	case "week":
// 		return 30 * time.Minute
// 	default:
// 		return 2 * time.Minute
// 	}
// }

// func cacheAndReturn[T any](ctx context.Context, c *gin.Context, redisClient *redis.Client, key string, response T, rangeType string) {
// 	data, err := json.Marshal(response)
// 	if err == nil {
// 		_ = redisClient.Set(ctx, key, data, getTTLByRange(rangeType)).Err()
// 	}
// 	c.JSON(http.StatusOK, response)
// }

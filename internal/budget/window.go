package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/arche/sentinelmesh/pkg/cache"
	"github.com/redis/go-redis/v9"
)

type WindowCalculator struct {
	cache *cache.RedisCache
}

func NewWindowCalculator(cache *cache.RedisCache) *WindowCalculator {
	return &WindowCalculator{cache: cache}
}

// AddMetric adds a latency or error metric to a Redis sorted set
func (w *WindowCalculator) AddMetric(ctx context.Context, serviceID string, isError bool, latency float64) error {
	now := time.Now().UnixMilli()
	
	pipeline := w.cache.Client.Pipeline()
	
	// Add to total requests sliding window (last 5 minutes)
	reqKey := fmt.Sprintf("req:%s", serviceID)
	pipeline.ZAdd(ctx, reqKey, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d-total", now)})
	pipeline.ZRemRangeByScore(ctx, reqKey, "0", fmt.Sprintf("%d", now-(5*60*1000)))
	pipeline.Expire(ctx, reqKey, 6*time.Minute)

	if isError {
		errKey := fmt.Sprintf("err:%s", serviceID)
		pipeline.ZAdd(ctx, errKey, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d-error", now)})
		pipeline.ZRemRangeByScore(ctx, errKey, "0", fmt.Sprintf("%d", now-(5*60*1000)))
		pipeline.Expire(ctx, errKey, 6*time.Minute)
	}

	latKey := fmt.Sprintf("lat:%s", serviceID)
	pipeline.ZAdd(ctx, latKey, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d-%.2f", now, latency)})
	pipeline.ZRemRangeByScore(ctx, latKey, "0", fmt.Sprintf("%d", now-(5*60*1000)))
	pipeline.Expire(ctx, latKey, 6*time.Minute)

	_, err := pipeline.Exec(ctx)
	return err
}

func (w *WindowCalculator) GetErrorRate(ctx context.Context, serviceID string) (float64, error) {
	reqKey := fmt.Sprintf("req:%s", serviceID)
	errKey := fmt.Sprintf("err:%s", serviceID)

	totalReqs, err := w.cache.Client.ZCard(ctx, reqKey).Result()
	if err != nil {
		return 0, err
	}

	if totalReqs == 0 {
		return 0, nil
	}

	errReqs, err := w.cache.Client.ZCard(ctx, errKey).Result()
	if err != nil {
		return 0, err
	}

	return float64(errReqs) / float64(totalReqs), nil
}

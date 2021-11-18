package limiter

import (
	"sync"
	"time"
)

type Limiter struct {
	mu       sync.Mutex
	capacity int
	duration time.Duration
	buckets  map[string]IBucket
}

func NewLimiter(capacity int, duration time.Duration) *Limiter {
	return &Limiter{
		mu:       sync.Mutex{},
		capacity: capacity,
		duration: duration,
		buckets:  make(map[string]IBucket),
	}
}

func (l *Limiter) GetBucket(key string) IBucket {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.buckets[key]; !ok {
		l.buckets[key] = &Bucket{
			mu:              sync.Mutex{},
			availableTokens: l.capacity,
			capacity:        l.capacity,
			duration:        l.duration,
			updatedAt:       time.Now(),
		}
	}

	return l.buckets[key]
}

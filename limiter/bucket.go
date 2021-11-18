package limiter

import (
	"sync"
	"time"
)

type IBucket interface {
	AddTokens()
	DeleteToken()
	ClearBucket()
	CheckTokensExist() bool
	CountAvailableTokens() int
}

type Bucket struct {
	mu              sync.Mutex
	availableTokens int
	capacity        int
	duration        time.Duration
	updatedAt       time.Time
}

func (b *Bucket) AddTokens() {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	if now.After(b.updatedAt.Add(b.duration)) {
		b.availableTokens = b.capacity
		b.updatedAt = now
	}
}

func (b *Bucket) DeleteToken() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.availableTokens > 0 {
		b.availableTokens--
	}
}

func (b *Bucket) CountAvailableTokens() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.availableTokens
}

func (b *Bucket) CheckTokensExist() bool {
	if b.CountAvailableTokens() <= 0 {
		return false
	}

	return true
}

func (b *Bucket) ClearBucket() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.availableTokens = 0
}

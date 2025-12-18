package auth

import (
	"context"
	"crypto/sha256"
	"sijil-core/internals/database"
	"sijil-core/internals/utils"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ShardCount = 256
	TTL        = 5 * time.Minute
)

type ProjectCacheEntry struct {
	ProjectID int
	Expires   int64
}

// The Room (shard)
type CasheShard struct {
	sync.RWMutex
	items map[string]ProjectCacheEntry
}

// The Hallway (The manager)
type AuthCache struct {
	shards [ShardCount]*CasheShard
	db     *pgxpool.Pool
}

func NewAuthCache(db *pgxpool.Pool) *AuthCache {
	c := &AuthCache{db: db}
	for i := range ShardCount {
		c.shards[i] = &CasheShard{
			items: make(map[string]ProjectCacheEntry),
		}
	}
	return c
}

func (c *AuthCache) getShard(apiKey string) *CasheShard {
	hash := sha256.Sum256([]byte(apiKey))

	shardIndex := hash[0]
	return c.shards[shardIndex]
}

func (c *AuthCache) ValidateAPIKey(ctx context.Context, apiKey, apiSecret string) (int, bool) {

	// Step 1: Check RAM (fast path)
	shard := c.getShard(apiKey)

	shard.RLock()
	entry, exists := shard.items[apiKey]
	shard.RUnlock()

	if exists && time.Now().Unix() > entry.Expires {
		return entry.ProjectID, true
	}

	// Step 2: Check DB (slow path - Cache Miss)
	project, err := database.GetProductByApiKey(ctx, c.db, apiKey)
	if err != nil {
		return 0, false
	}

	// Step 3: Check secret (slow math)
	if !utils.ComparePasswordHash(apiSecret, project.ApiSecretHash) {
		return 0, false
	}

	// Step 4: Update cache
	shard.Lock()
	shard.items[apiKey] = ProjectCacheEntry{
		ProjectID: project.ID,
		Expires:   time.Now().Add(TTL).Unix(),
	}
	shard.Unlock()

	return project.ID, true
}

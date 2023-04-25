package auth

import (
	"sync"
	"time"
)

type GroupCache struct {
	gg    GroupGetter
	cache map[string]map[string]bool
	stamp map[string]int64
	mux   sync.RWMutex
}

func NewCache(gg GroupGetter) *GroupCache {
	gcache := GroupCache{
		gg:    gg,
		cache: make(map[string]map[string]bool),
		stamp: make(map[string]int64),
	}
	return &gcache
}

func (gc *GroupCache) GroupsForUser(username string) (map[string]bool, error) {

	// check the cache first
	groups := gc.checkCache(username)
	if groups != nil {
		return groups, nil
	}

	// load from getter
	groups, err := gc.loadGroups(username)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (gc *GroupCache) checkCache(username string) map[string]bool {
	gc.mux.RLock()
	defer gc.mux.RUnlock()

	groups, ok := gc.cache[username]
	if ok {
		stamp := gc.stamp[username]
		if time.Now().Unix()-stamp < 600 {
			return groups
		}
	}

	return nil
}

func (gc *GroupCache) loadGroups(username string) (map[string]bool, error) {
	gc.mux.Lock()
	defer gc.mux.Unlock()

	now := time.Now().Unix()

	// prune all expired cache entries
	for user, stamp := range gc.stamp {
		if now-stamp > 600 {
			delete(gc.stamp, user)
			delete(gc.cache, user)
		}
	}

	// either we don't have the user cached, or the cache is stale
	groups, err := gc.gg.GroupsForUser(username)
	if err != nil {
		return nil, err
	}
	gc.cache[username] = groups
	gc.stamp[username] = now

	return groups, nil
}

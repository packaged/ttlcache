package ttlmap

import "sync"

var backgroundMutex = newKeyMutex()

type keyMutex struct {
	c *sync.Cond
	l sync.Locker
	s map[string]struct{}
}

// Create new keyMutex
func newKeyMutex() *keyMutex {
	l := sync.Mutex{}
	return &keyMutex{c: sync.NewCond(&l), l: &l, s: make(map[string]struct{})}
}

// Unlock keyMutex by unique ID
func (km *keyMutex) Unlock(key string) {
	km.l.Lock()
	defer km.l.Unlock()
	delete(km.s, key)
	km.c.Broadcast()
}

// Lock keyMutex by unique ID
func (km *keyMutex) Lock(key string) bool {
	km.l.Lock()
	defer km.l.Unlock()
	if _, ok := km.s[key]; ok {
		return false
	}
	km.s[key] = struct{}{}
	return true
}

func (m CacheMap) BackgroundUpdate(key string, updater func() (interface{}, error)) {
	//Lock the key from writes
	locked := backgroundMutex.Lock(key)
	if locked {
		//Defer release write lock
		defer backgroundMutex.Unlock(key)
		value, err := updater()
		if err == nil {
			m.Set(key, value, nil)
		}
	}
}

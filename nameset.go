package targets

import (
	"sync"
)

type NameSet struct {
	names map[string]struct{}
	lock  sync.RWMutex
}

func (t *NameSet) AddOnce(name string) bool {
	t.lock.RLock()
	_, ok := t.names[name]
	t.lock.RUnlock()

	if ok {
		return false
	}

	t.lock.Lock()
	defer t.lock.Unlock()
	_, ok = t.names[name]
	if ok {
		return false
	}
	t.names[name] = struct{}{}
	return true
}

func MakeNameSet() *NameSet {
	out := new(NameSet)
	out.names = make(map[string]struct{})
	return out
}

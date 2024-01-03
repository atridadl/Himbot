package lib

import (
	"sync"
	"time"
)

type CooldownManager struct {
	cooldowns map[string]time.Time
	mu        sync.Mutex
}

func NewCooldownManager() *CooldownManager {
	return &CooldownManager{
		cooldowns: make(map[string]time.Time),
	}
}

func (m *CooldownManager) StartCooldown(key string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cooldowns[key] = time.Now().Add(duration)
}

func (m *CooldownManager) IsOnCooldown(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	cooldownEnd, exists := m.cooldowns[key]
	if !exists {
		return false
	}

	if time.Now().After(cooldownEnd) {
		delete(m.cooldowns, key)
		return false
	}

	return true
}

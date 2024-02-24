package lib

import (
	"sync"
	"time"
)

var (
	mu       sync.Mutex
	instance *CooldownManager
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

func GetInstance() *CooldownManager {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		instance = &CooldownManager{
			cooldowns: make(map[string]time.Time),
		}
	}

	return instance
}

func (m *CooldownManager) StartCooldown(userID string, key string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cooldowns[userID+":"+key] = time.Now().Add(duration)
}

func (m *CooldownManager) IsOnCooldown(userID string, key string) (bool, time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cooldownEnd, exists := m.cooldowns[userID+":"+key]
	if !exists {
		return false, 0
	}

	if time.Now().After(cooldownEnd) {
		delete(m.cooldowns, userID+":"+key)
		return false, 0
	}

	return true, time.Until(cooldownEnd)
}

func CancelCooldown(userID string, key string) {
	manager := GetInstance()

	// Handle non-existent keys gracefully
	if _, exists := manager.cooldowns[userID+":"+key]; !exists {
		return
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()
	delete(manager.cooldowns, userID+":"+key)
}

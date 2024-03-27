package lib

import (
	"sync"
	"time"
)

var (
	mu       sync.Mutex
	instance *TimerManager
)

type TimerManager struct {
	timers map[string]time.Time
	mu     sync.Mutex
}

func NewTimerManager() *TimerManager {
	return &TimerManager{
		timers: make(map[string]time.Time),
	}
}

func GetInstance() *TimerManager {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		instance = &TimerManager{
			timers: make(map[string]time.Time),
		}
	}

	return instance
}

func (m *TimerManager) StartTimer(userID string, key string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.timers[userID+":"+key] = time.Now().Add(duration)
}

func (m *TimerManager) TimerRunning(userID string, key string) (bool, time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	timerEnd, exists := m.timers[userID+":"+key]
	if !exists {
		return false, 0
	}

	if time.Now().After(timerEnd) {
		delete(m.timers, userID+":"+key)
		return false, 0
	}

	return true, time.Until(timerEnd)
}

func CancelTimer(userID string, key string) {
	manager := GetInstance()

	// Handle non-existent keys gracefully
	if _, exists := manager.timers[userID+":"+key]; !exists {
		return
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()
	delete(manager.timers, userID+":"+key)
}

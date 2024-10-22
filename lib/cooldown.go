package lib

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	once     sync.Once
	instance *CooldownManager
)

type CooldownManager struct {
	cooldowns map[string]time.Time
	mu        sync.Mutex
}

func GetCooldownManager() *CooldownManager {
	once.Do(func() {
		instance = &CooldownManager{
			cooldowns: make(map[string]time.Time),
		}
	})
	return instance
}

func (cm *CooldownManager) SetCooldown(userID, commandName string, duration time.Duration) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cooldowns[userID+":"+commandName] = time.Now().Add(duration)
}

func (cm *CooldownManager) CheckCooldown(userID, commandName string) (bool, time.Duration) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	key := userID + ":" + commandName
	if cooldownEnd, exists := cm.cooldowns[key]; exists {
		if time.Now().Before(cooldownEnd) {
			return false, time.Until(cooldownEnd)
		}
		delete(cm.cooldowns, key)
	}
	return true, 0
}

func CheckAndApplyCooldown(s *discordgo.Session, i *discordgo.InteractionCreate, commandName string, duration time.Duration) bool {
	cooldownManager := GetCooldownManager()
	user, err := GetUser(i)
	if err != nil {
		RespondWithError(s, i, "Error processing command: "+err.Error())
		return false
	}

	canUse, remaining := cooldownManager.CheckCooldown(user.ID, commandName)
	if !canUse {
		RespondWithError(s, i, fmt.Sprintf("You can use this command again in %v", remaining.Round(time.Second)))
		return false
	}

	cooldownManager.SetCooldown(user.ID, commandName, duration)
	return true
}

package config

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// RedirectType 跳转类型
type RedirectType int

const (
	RedirectType301 RedirectType = 301 // HTTP 301 永久重定向
	RedirectType302 RedirectType = 302 // HTTP 302 临时重定向
	RedirectType307 RedirectType = 307 // HTTP 307 临时重定向（保持方法）
	RedirectTypeJS  RedirectType = 4   // JavaScript 跳转
)

// RedirectRule 跳转规则
type RedirectRule struct {
	ID          string       `json:"id"`           // 规则ID
	Domain      string       `json:"domain"`       // 域名
	Path        string       `json:"path"`         // 路径（可选）
	Target      string       `json:"target"`       // 目标URL
	Type        RedirectType `json:"type"`         // 跳转类型
	ExpiresAt   *time.Time   `json:"expires_at"`   // 过期时间（nil表示永不过期）
	CreatedAt   time.Time    `json:"created_at"`   // 创建时间
	Description string       `json:"description"`  // 描述
}

// IsExpired 检查规则是否已过期
func (r *RedirectRule) IsExpired() bool {
	if r.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*r.ExpiresAt)
}

// Config 配置管理
type Config struct {
	Port           int    `json:"port"`             // 服务端口
	LogFile        string `json:"log_file"`         // 日志文件路径
	ConfigFile     string `json:"config_file"`      // 配置文件路径
	LogBufferSize  int    `json:"log_buffer_size"`  // 日志缓冲大小
	LogFlushInterval int  `json:"log_flush_interval"` // 日志刷新间隔（秒）
	rules          *sync.Map
	mu             sync.RWMutex
}

var defaultConfig = &Config{
	Port:            8080,
	LogFile:         "access.log",
	ConfigFile:      "rules.json",
	LogBufferSize:   1000,
	LogFlushInterval: 180,
	rules:           &sync.Map{},
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return defaultConfig
}

// GetRule 获取跳转规则
func (c *Config) GetRule(key string) (*RedirectRule, bool) {
	value, ok := c.rules.Load(key)
	if !ok {
		return nil, false
	}
	rule := value.(*RedirectRule)
	if rule.IsExpired() {
		c.rules.Delete(key)
		return nil, false
	}
	return rule, true
}

// SetRule 设置跳转规则
func (c *Config) SetRule(rule *RedirectRule) {
	key := c.generateKey(rule.Domain, rule.Path)
	c.rules.Store(key, rule)
}

// DeleteRule 删除跳转规则
func (c *Config) DeleteRule(domain, path string) {
	key := c.generateKey(domain, path)
	c.rules.Delete(key)
}

// GetAllRules 获取所有规则
func (c *Config) GetAllRules() []*RedirectRule {
	var rules []*RedirectRule
	c.rules.Range(func(key, value interface{}) bool {
		rule := value.(*RedirectRule)
		if !rule.IsExpired() {
			rules = append(rules, rule)
		}
		return true
	})
	return rules
}

// generateKey 生成规则键
func (c *Config) generateKey(domain, path string) string {
	if path == "" {
		return domain
	}
	return domain + "|" + path
}

// LoadFromFile 从文件加载配置
func (c *Config) LoadFromFile() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := os.ReadFile(c.ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，使用空配置
		}
		return err
	}

	var rules []*RedirectRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return err
	}

	c.rules = &sync.Map{}
	for _, rule := range rules {
		if !rule.IsExpired() {
			key := c.generateKey(rule.Domain, rule.Path)
			c.rules.Store(key, rule)
		}
	}

	return nil
}

// SaveToFile 保存配置到文件
func (c *Config) SaveToFile() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	rules := c.GetAllRules()
	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.ConfigFile, data, 0644)
}

// FindRule 查找匹配的规则（优先精确路径匹配，再域名匹配）
func (c *Config) FindRule(domain, path string) (*RedirectRule, bool) {
	// 先尝试精确匹配（域名+路径）
	if path != "" {
		key := c.generateKey(domain, path)
		if rule, ok := c.GetRule(key); ok {
			return rule, true
		}
	}

	// 再尝试域名匹配
	key := c.generateKey(domain, "")
	if rule, ok := c.GetRule(key); ok {
		return rule, true
	}

	return nil, false
}

// CheckConflict 检查规则冲突
// 返回冲突的规则列表和冲突描述
func (c *Config) CheckConflict(rule *RedirectRule, excludeID string) ([]*RedirectRule, string) {
	var conflicts []*RedirectRule
	var conflictMsg string

	key := c.generateKey(rule.Domain, rule.Path)
	existingRule, exists := c.GetRule(key)
	
	if exists && existingRule.ID != excludeID {
		conflicts = append(conflicts, existingRule)
		conflictMsg = "存在完全相同的规则（域名和路径都相同）"
		return conflicts, conflictMsg
	}

	// 检查是否有更具体的规则冲突
	// 如果新规则是域名级别，检查是否有该域名的路径规则
	if rule.Path == "" {
		rules := c.GetAllRules()
		for _, r := range rules {
			if r.ID != excludeID && r.Domain == rule.Domain && r.Path != "" {
				conflicts = append(conflicts, r)
				conflictMsg = "存在该域名的路径级别规则，域名级别规则会覆盖所有路径规则"
			}
		}
		if len(conflicts) > 0 {
			return conflicts, conflictMsg
		}
	}

	// 检查是否有域名级别的规则会覆盖当前路径规则
	if rule.Path != "" {
		domainKey := c.generateKey(rule.Domain, "")
		domainRule, exists := c.GetRule(domainKey)
		if exists && domainRule.ID != excludeID {
			conflicts = append(conflicts, domainRule)
			conflictMsg = "存在该域名的域名级别规则，会优先匹配并覆盖此路径规则"
			return conflicts, conflictMsg
		}
	}

	return nil, ""
}

package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"mini_jump/config"
)

// API API 管理接口
type API struct {
	config *config.Config
}

// NewAPI 创建 API 处理器
func NewAPI(cfg *config.Config) *API {
	return &API{
		config: cfg,
	}
}

// RegisterRoutes 注册 API 路由
func (a *API) RegisterRoutes(r *mux.Router) {
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/rules", a.ListRules).Methods("GET")
	apiRouter.HandleFunc("/rules", a.CreateRule).Methods("POST")
	apiRouter.HandleFunc("/rules/{id}", a.GetRule).Methods("GET")
	apiRouter.HandleFunc("/rules/{id}", a.UpdateRule).Methods("PUT")
	apiRouter.HandleFunc("/rules/{id}", a.DeleteRule).Methods("DELETE")
	apiRouter.HandleFunc("/reload", a.ReloadConfig).Methods("POST")
	apiRouter.HandleFunc("/save", a.SaveConfig).Methods("POST")
}

// ListRules 列出所有规则
func (a *API) ListRules(w http.ResponseWriter, r *http.Request) {
	rules := a.config.GetAllRules()
	respondJSON(w, http.StatusOK, rules)
}

// CreateRule 创建规则
func (a *API) CreateRule(w http.ResponseWriter, r *http.Request) {
	var rule config.RedirectRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 验证必填字段
	if rule.Domain == "" || rule.Target == "" {
		respondError(w, http.StatusBadRequest, "域名和目标URL不能为空")
		return
	}

	// 生成 ID
	if rule.ID == "" {
		rule.ID = a.generateID(rule.Domain, rule.Path)
	}

	// 检查冲突
	conflicts, conflictMsg := a.config.CheckConflict(&rule, "")
	if len(conflicts) > 0 {
		respondJSON(w, http.StatusConflict, map[string]interface{}{
			"error":     conflictMsg,
			"conflicts": conflicts,
		})
		return
	}

	rule.CreatedAt = time.Now()
	a.config.SetRule(&rule)
	a.config.SaveToFile()

	respondJSON(w, http.StatusCreated, rule)
}

// GetRule 获取规则
func (a *API) GetRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	rules := a.config.GetAllRules()
	for _, rule := range rules {
		if rule.ID == id {
			respondJSON(w, http.StatusOK, rule)
			return
		}
	}

	respondError(w, http.StatusNotFound, "Rule not found")
}

// UpdateRule 更新规则
func (a *API) UpdateRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedRule config.RedirectRule
	if err := json.NewDecoder(r.Body).Decode(&updatedRule); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 验证必填字段
	if updatedRule.Domain == "" || updatedRule.Target == "" {
		respondError(w, http.StatusBadRequest, "域名和目标URL不能为空")
		return
	}

	rules := a.config.GetAllRules()
	for _, rule := range rules {
		if rule.ID == id {
			updatedRule.ID = id
			if updatedRule.CreatedAt.IsZero() {
				updatedRule.CreatedAt = rule.CreatedAt
			}

			// 检查冲突（排除当前规则）
			conflicts, conflictMsg := a.config.CheckConflict(&updatedRule, id)
			if len(conflicts) > 0 {
				respondJSON(w, http.StatusConflict, map[string]interface{}{
					"error":     conflictMsg,
					"conflicts": conflicts,
				})
				return
			}

			a.config.SetRule(&updatedRule)
			a.config.SaveToFile()
			respondJSON(w, http.StatusOK, updatedRule)
			return
		}
	}

	respondError(w, http.StatusNotFound, "Rule not found")
}

// DeleteRule 删除规则
func (a *API) DeleteRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	rules := a.config.GetAllRules()
	for _, rule := range rules {
		if rule.ID == id {
			a.config.DeleteRule(rule.Domain, rule.Path)
			a.config.SaveToFile()
			respondJSON(w, http.StatusOK, map[string]string{"message": "Rule deleted"})
			return
		}
	}

	respondError(w, http.StatusNotFound, "Rule not found")
}

// ReloadConfig 重新加载配置
func (a *API) ReloadConfig(w http.ResponseWriter, r *http.Request) {
	if err := a.config.LoadFromFile(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to reload config: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Config reloaded"})
}

// SaveConfig 保存配置
func (a *API) SaveConfig(w http.ResponseWriter, r *http.Request) {
	if err := a.config.SaveToFile(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save config: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Config saved"})
}

// generateID 生成规则 ID
func (a *API) generateID(domain, path string) string {
	if path == "" {
		return strings.ReplaceAll(domain, ".", "_")
	}
	return strings.ReplaceAll(domain+"_"+path, ".", "_")
}

// respondJSON 返回 JSON 响应
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError 返回错误响应
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

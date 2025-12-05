package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"mini_jump/config"
	"mini_jump/logger"
)

// Handler HTTP 请求处理器
type Handler struct {
	config *config.Config
	logger *logger.Logger
}

// NewHandler 创建处理器
func NewHandler(cfg *config.Config, log *logger.Logger) *Handler {
	return &Handler{
		config: cfg,
		logger: log,
	}
}

// HandleRedirect 处理跳转请求
func (h *Handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	domain := r.Host
	path := r.URL.Path

	// 查找匹配的规则
	rule, found := h.config.FindRule(domain, path)
	if !found {
		http.NotFound(w, r)
		return
	}

	// 记录访问日志
	accessLog := &logger.AccessLog{
		Timestamp:    time.Now(),
		IP:           h.getClientIP(r),
		UserAgent:    r.UserAgent(),
		Method:       r.Method,
		Domain:       domain,
		Path:         path,
		Target:       rule.Target,
		RedirectType: int(rule.Type),
		StatusCode:   int(rule.Type),
	}

	// 执行跳转
	switch rule.Type {
	case config.RedirectType301:
		http.Redirect(w, r, rule.Target, http.StatusMovedPermanently)
		accessLog.StatusCode = http.StatusMovedPermanently
	case config.RedirectType302:
		http.Redirect(w, r, rule.Target, http.StatusFound)
		accessLog.StatusCode = http.StatusFound
	case config.RedirectType307:
		http.Redirect(w, r, rule.Target, http.StatusTemporaryRedirect)
		accessLog.StatusCode = http.StatusTemporaryRedirect
	case config.RedirectTypeJS:
		h.writeJavaScriptRedirect(w, rule.Target)
		accessLog.StatusCode = http.StatusOK
	default:
		http.Redirect(w, r, rule.Target, http.StatusFound)
		accessLog.StatusCode = http.StatusFound
	}

	// 异步记录日志
	go h.logger.Log(accessLog)
}

// writeJavaScriptRedirect 写入 JavaScript 跳转
func (h *Handler) writeJavaScriptRedirect(w http.ResponseWriter, target string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta http-equiv="refresh" content="0;url=%s">
<script>window.location.href="%s";</script>
</head>
<body>正在跳转到 %s...</body>
</html>`, target, target, target)
	w.Write([]byte(html))
}

// getClientIP 获取客户端 IP
func (h *Handler) getClientIP(r *http.Request) string {
	// 尝试从 X-Forwarded-For 获取
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 尝试从 X-Real-IP 获取
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// 使用 RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

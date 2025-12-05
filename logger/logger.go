package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// AccessLog 访问日志
type AccessLog struct {
	Timestamp   time.Time `json:"timestamp"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	Method      string    `json:"method"`
	Domain      string    `json:"domain"`
	Path        string    `json:"path"`
	Target      string    `json:"target"`
	RedirectType int      `json:"redirect_type"`
	StatusCode  int       `json:"status_code"`
}

// Logger 日志管理器
type Logger struct {
	logFile       string
	bufferSize    int
	flushInterval time.Duration
	buffer        []*AccessLog
	mu            sync.Mutex
	file          *os.File
	ticker        *time.Ticker
	done          chan struct{}
}

// NewLogger 创建日志管理器
func NewLogger(logFile string, bufferSize int, flushIntervalSeconds int) (*Logger, error) {
	l := &Logger{
		logFile:       logFile,
		bufferSize:    bufferSize,
		flushInterval: time.Duration(flushIntervalSeconds) * time.Second,
		buffer:        make([]*AccessLog, 0, bufferSize),
		done:          make(chan struct{}),
	}

	// 打开或创建日志文件
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	l.file = file

	// 启动定时刷新
	l.ticker = time.NewTicker(l.flushInterval)
	go l.flushPeriodically()

	return l, nil
}

// Log 记录访问日志
func (l *Logger) Log(log *AccessLog) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.buffer = append(l.buffer, log)

	// 达到缓冲大小，立即刷新
	if len(l.buffer) >= l.bufferSize {
		l.flushUnlocked()
	}
}

// flushUnlocked 刷新日志到文件（不加锁版本，需要在锁内调用）
func (l *Logger) flushUnlocked() {
	if len(l.buffer) == 0 {
		return
	}

	for _, log := range l.buffer {
		data, err := json.Marshal(log)
		if err != nil {
			fmt.Printf("Error marshaling log: %v\n", err)
			continue
		}
		l.file.Write(append(data, '\n'))
	}

	l.buffer = l.buffer[:0] // 清空缓冲
}

// Flush 刷新日志到文件
func (l *Logger) Flush() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flushUnlocked()
	return nil
}

// flushPeriodically 定期刷新日志
func (l *Logger) flushPeriodically() {
	for {
		select {
		case <-l.ticker.C:
			l.Flush()
		case <-l.done:
			return
		}
	}
}

// Close 关闭日志管理器
func (l *Logger) Close() error {
	close(l.done)
	if l.ticker != nil {
		l.ticker.Stop()
	}
	l.Flush()
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

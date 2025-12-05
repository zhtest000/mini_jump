package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"

	"syscall"

	"github.com/gorilla/mux"

	"mini_jump/api"
	"mini_jump/config"
	"mini_jump/handler"
	"mini_jump/logger"
	"mini_jump/manager"
	"mini_jump/service"
)

func main() {
	// 检查是否为 install 或 uninstall 命令
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "install":
			handleInstall()
			return
		case "uninstall":
			handleUninstall()
			return
		}
	}

	// 解析命令行参数
	port := flag.Int("port", 18082, "服务端口")
	configFile := flag.String("config", "rules.json", "配置文件路径")
	logFile := flag.String("log", "access.log", "日志文件路径")
	logBufferSize := flag.Int("log-buffer", 1000, "日志缓冲大小")
	logFlushInterval := flag.Int("log-flush", 180, "日志刷新间隔（秒）")
	flag.Parse()

	// 初始化配置
	cfg := config.GetDefaultConfig()
	cfg.Port = *port
	cfg.ConfigFile = *configFile
	cfg.LogFile = *logFile
	cfg.LogBufferSize = *logBufferSize
	cfg.LogFlushInterval = *logFlushInterval

	// 加载配置
	if err := cfg.LoadFromFile(); err != nil {
		log.Printf("Warning: Failed to load config: %v\n", err)
	}

	// 初始化日志
	accessLogger, err := logger.NewLogger(cfg.LogFile, cfg.LogBufferSize, cfg.LogFlushInterval)
	if err != nil {
		log.Fatalf("Failed to create logger: %v\n", err)
	}
	defer accessLogger.Close()

	// 初始化处理器
	redirectHandler := handler.NewHandler(cfg, accessLogger)

	// 初始化 API
	apiHandler := api.NewAPI(cfg)

	// 初始化管理页面
	managerHandler := manager.NewManager()

	// 设置路由
	router := mux.NewRouter()

	// 管理页面路由
	router.HandleFunc("/manager558630", managerHandler.ServeManager)

	// API 路由
	apiHandler.RegisterRoutes(router)

	// 跳转路由（所有其他请求）
	router.PathPrefix("/").HandlerFunc(redirectHandler.HandleRedirect)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")

		// 保存配置
		if err := cfg.SaveToFile(); err != nil {
			log.Printf("Failed to save config: %v\n", err)
		}

		// 刷新日志
		if err := accessLogger.Flush(); err != nil {
			log.Printf("Failed to flush logs: %v\n", err)
		}

		os.Exit(0)
	}()

	log.Printf("MiniJump HTTP Redirect Service starting on port %d\n", cfg.Port)
	log.Printf("Config file: %s\n", cfg.ConfigFile)
	log.Printf("Log file: %s\n", cfg.LogFile)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v\n", err)
	}
}

// handleInstall 处理安装服务命令
func handleInstall() {
	// 解析安装参数
	installFlags := flag.NewFlagSet("install", flag.ExitOnError)
	port := installFlags.Int("port", 8080, "服务端口")
	configFile := installFlags.String("config", "rules.json", "配置文件路径")
	logFile := installFlags.String("log", "access.log", "日志文件路径")
	logBufferSize := installFlags.Int("log-buffer", 1000, "日志缓冲大小")
	logFlushInterval := installFlags.Int("log-flush", 180, "日志刷新间隔（秒）")
	serviceName := installFlags.String("name", "MiniJump", "服务名称")
	installFlags.Parse(os.Args[2:])

	// 构建服务启动参数
	var args []string
	if *port != 8080 {
		args = append(args, fmt.Sprintf("-port=%d", *port))
	}
	if *configFile != "rules.json" {
		args = append(args, fmt.Sprintf("-config=%s", *configFile))
	}
	if *logFile != "access.log" {
		args = append(args, fmt.Sprintf("-log=%s", *logFile))
	}
	if *logBufferSize != 1000 {
		args = append(args, fmt.Sprintf("-log-buffer=%d", *logBufferSize))
	}
	if *logFlushInterval != 180 {
		args = append(args, fmt.Sprintf("-log-flush=%d", *logFlushInterval))
	}

	// 检查权限
	if runtime.GOOS == "windows" {
		if !isAdminWindows() {
			fmt.Println("错误: 需要管理员权限才能安装服务")
			fmt.Println("请以管理员身份运行此命令")
			os.Exit(1)
		}
	} else {
		if os.Geteuid() != 0 {
			fmt.Println("错误: 需要 root 权限才能安装服务")
			fmt.Println("请使用 sudo 运行此命令")
			os.Exit(1)
		}
	}

	// 创建服务管理器
	sm, err := service.NewServiceManager(
		*serviceName,
		"MiniJump HTTP Redirect Service",
		"MiniJump 是一个轻量级的 HTTP 跳转服务，支持基于域名和路径的智能跳转配置",
	)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	sm.SetArgs(args)

	// 安装服务
	if err := sm.Install(); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
}

// handleUninstall 处理卸载服务命令
func handleUninstall() {
	// 解析卸载参数
	uninstallFlags := flag.NewFlagSet("uninstall", flag.ExitOnError)
	serviceName := uninstallFlags.String("name", "MiniJump", "服务名称")
	uninstallFlags.Parse(os.Args[2:])

	// 检查权限
	if runtime.GOOS == "windows" {
		if !isAdminWindows() {
			fmt.Println("错误: 需要管理员权限才能卸载服务")
			fmt.Println("请以管理员身份运行此命令")
			os.Exit(1)
		}
	} else {
		if os.Geteuid() != 0 {
			fmt.Println("错误: 需要 root 权限才能卸载服务")
			fmt.Println("请使用 sudo 运行此命令")
			os.Exit(1)
		}
	}

	// 创建服务管理器
	sm, err := service.NewServiceManager(
		*serviceName,
		"MiniJump HTTP Redirect Service",
		"MiniJump HTTP Redirect Service",
	)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 卸载服务
	if err := sm.Uninstall(); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
}

// isAdminWindows 检查是否有 Windows 管理员权限
func isAdminWindows() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	cmd := exec.Command("net", "session")
	err := cmd.Run()
	return err == nil
}

package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ServiceManager 服务管理器
type ServiceManager struct {
	ServiceName string
	DisplayName string
	Description string
	ExecPath    string
	WorkDir     string
	Args        []string
}

// NewServiceManager 创建服务管理器
func NewServiceManager(serviceName, displayName, description string) (*ServiceManager, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("获取可执行文件路径失败: %v", err)
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(execPath)
	if err != nil {
		return nil, fmt.Errorf("获取绝对路径失败: %v", err)
	}

	workDir := filepath.Dir(absPath)

	return &ServiceManager{
		ServiceName: serviceName,
		DisplayName: displayName,
		Description: description,
		ExecPath:    absPath,
		WorkDir:     workDir,
	}, nil
}

// SetArgs 设置服务启动参数
func (sm *ServiceManager) SetArgs(args []string) {
	sm.Args = args
}

// Install 安装服务
func (sm *ServiceManager) Install() error {
	if runtime.GOOS == "windows" {
		return sm.installWindows()
	} else {
		return sm.installLinux()
	}
}

// Uninstall 卸载服务
func (sm *ServiceManager) Uninstall() error {
	if runtime.GOOS == "windows" {
		return sm.uninstallWindows()
	} else {
		return sm.uninstallLinux()
	}
}

// installWindows 在 Windows 上安装服务
func (sm *ServiceManager) installWindows() error {
	// 检查管理员权限
	if !isAdmin() {
		return fmt.Errorf("需要管理员权限才能安装服务。请以管理员身份运行")
	}

	// 检查服务是否已存在
	if sm.isServiceExists() {
		return fmt.Errorf("服务 '%s' 已存在", sm.ServiceName)
	}

	// 构建命令参数
	args := sm.Args
	argsStr := ""
	if len(args) > 0 {
		argsStr = strings.Join(args, " ")
	}

	// 使用 sc.exe 创建服务
	cmd := exec.Command("sc.exe", "create", sm.ServiceName,
		fmt.Sprintf("binPath= \"%s\" %s", sm.ExecPath, argsStr),
		fmt.Sprintf("DisplayName= %s", sm.DisplayName),
		"start= auto",
		fmt.Sprintf("obj= %s", "LocalSystem"),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("创建服务失败: %v, 输出: %s", err, string(output))
	}

	// 设置服务描述
	cmd = exec.Command("sc.exe", "description", sm.ServiceName, sm.Description)
	output, err = cmd.CombinedOutput()
	if err != nil {
		// 描述设置失败不影响服务创建
		fmt.Printf("警告: 设置服务描述失败: %v\n", err)
	}

	fmt.Printf("✓ 服务 '%s' 安装成功\n", sm.ServiceName)
	fmt.Printf("  服务名称: %s\n", sm.ServiceName)
	fmt.Printf("  显示名称: %s\n", sm.DisplayName)
	fmt.Printf("  可执行文件: %s\n", sm.ExecPath)
	fmt.Printf("\n使用以下命令管理服务:\n")
	fmt.Printf("  启动服务: sc start %s\n", sm.ServiceName)
	fmt.Printf("  停止服务: sc stop %s\n", sm.ServiceName)
	fmt.Printf("  删除服务: sc delete %s\n", sm.ServiceName)

	return nil
}

// uninstallWindows 在 Windows 上卸载服务
func (sm *ServiceManager) uninstallWindows() error {
	// 检查管理员权限
	if !isAdmin() {
		return fmt.Errorf("需要管理员权限才能卸载服务。请以管理员身份运行")
	}

	// 检查服务是否存在
	if !sm.isServiceExists() {
		return fmt.Errorf("服务 '%s' 不存在", sm.ServiceName)
	}

	// 先停止服务
	fmt.Printf("正在停止服务 '%s'...\n", sm.ServiceName)
	stopCmd := exec.Command("sc.exe", "stop", sm.ServiceName)
	stopCmd.Run() // 忽略错误，服务可能已经停止

	// 删除服务
	fmt.Printf("正在删除服务 '%s'...\n", sm.ServiceName)
	cmd := exec.Command("sc.exe", "delete", sm.ServiceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("删除服务失败: %v, 输出: %s", err, string(output))
	}

	fmt.Printf("✓ 服务 '%s' 卸载成功\n", sm.ServiceName)
	return nil
}

// installLinux 在 Linux 上安装服务（systemd）
func (sm *ServiceManager) installLinux() error {
	// 检查是否为 root 用户
	if os.Geteuid() != 0 {
		return fmt.Errorf("需要 root 权限才能安装服务。请使用 sudo 运行")
	}

	// 检查 systemd 是否存在
	if _, err := exec.LookPath("systemctl"); err != nil {
		return fmt.Errorf("未找到 systemctl，请确保系统使用 systemd")
	}

	// 构建服务单元文件内容
	serviceContent := sm.generateSystemdUnit()

	// systemd 服务文件路径
	serviceFile := fmt.Sprintf("/etc/systemd/system/%s.service", sm.ServiceName)

	// 检查服务是否已存在
	if _, err := os.Stat(serviceFile); err == nil {
		return fmt.Errorf("服务 '%s' 已存在", sm.ServiceName)
	}

	// 写入服务文件
	err := os.WriteFile(serviceFile, []byte(serviceContent), 0644)
	if err != nil {
		return fmt.Errorf("创建服务文件失败: %v", err)
	}

	// 重新加载 systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("重新加载 systemd 失败: %v", err)
	}

	// 启用服务（开机自启）
	cmd = exec.Command("systemctl", "enable", sm.ServiceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启用服务失败: %v", err)
	}

	fmt.Printf("✓ 服务 '%s' 安装成功\n", sm.ServiceName)
	fmt.Printf("  服务文件: %s\n", serviceFile)
	fmt.Printf("  可执行文件: %s\n", sm.ExecPath)
	fmt.Printf("\n使用以下命令管理服务:\n")
	fmt.Printf("  启动服务: sudo systemctl start %s\n", sm.ServiceName)
	fmt.Printf("  停止服务: sudo systemctl stop %s\n", sm.ServiceName)
	fmt.Printf("  查看状态: sudo systemctl status %s\n", sm.ServiceName)
	fmt.Printf("  删除服务: sudo systemctl disable %s && sudo rm %s\n", sm.ServiceName, serviceFile)

	return nil
}

// uninstallLinux 在 Linux 上卸载服务
func (sm *ServiceManager) uninstallLinux() error {
	// 检查是否为 root 用户
	if os.Geteuid() != 0 {
		return fmt.Errorf("需要 root 权限才能卸载服务。请使用 sudo 运行")
	}

	serviceFile := fmt.Sprintf("/etc/systemd/system/%s.service", sm.ServiceName)

	// 检查服务文件是否存在
	if _, err := os.Stat(serviceFile); os.IsNotExist(err) {
		return fmt.Errorf("服务 '%s' 不存在", sm.ServiceName)
	}

	// 先停止和禁用服务
	fmt.Printf("正在停止并禁用服务 '%s'...\n", sm.ServiceName)
	exec.Command("systemctl", "stop", sm.ServiceName).Run()
	exec.Command("systemctl", "disable", sm.ServiceName).Run()

	// 删除服务文件
	err := os.Remove(serviceFile)
	if err != nil {
		return fmt.Errorf("删除服务文件失败: %v", err)
	}

	// 重新加载 systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	cmd.Run()

	fmt.Printf("✓ 服务 '%s' 卸载成功\n", sm.ServiceName)
	return nil
}

// generateSystemdUnit 生成 systemd 服务单元文件内容
func (sm *ServiceManager) generateSystemdUnit() string {
	args := strings.Join(sm.Args, " ")
	if args != "" {
		args = " " + args
	}

	return fmt.Sprintf(`[Unit]
Description=%s
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=%s
ExecStart=%s%s
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
`, sm.Description, sm.WorkDir, sm.ExecPath, args)
}

// isServiceExists 检查 Windows 服务是否存在
func (sm *ServiceManager) isServiceExists() bool {
	cmd := exec.Command("sc.exe", "query", sm.ServiceName)
	err := cmd.Run()
	return err == nil
}

// isAdmin 检查是否有管理员权限（Windows）
func isAdmin() bool {
	if runtime.GOOS != "windows" {
		return os.Geteuid() == 0
	}

	cmd := exec.Command("net", "session")
	err := cmd.Run()
	return err == nil
}

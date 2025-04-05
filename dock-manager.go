package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

type DockApp struct {
	BundleID string
	Label    string
	Position int
}

func main() {
	// 获取当前用户
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("Error getting current user: %v\n", err)
		return
	}

	// 获取 Dock 应用信息
	apps, err := getDockApps(currentUser.Username)
	if err != nil {
		fmt.Printf("Error getting dock apps: %v\n", err)
		return
	}

	// 准备多选项
	var options []string
	for _, app := range apps {
		options = append(options, fmt.Sprintf("%s (%s)", app.Label, app.BundleID))
	}

	// 创建多选提示
	var selected []string
	prompt := &survey.MultiSelect{
		Message: "选择要移除的应用（空格选择，回车确认）:",
		Options: options,
	}

	// 获取用户选择
	err = survey.AskOne(prompt, &selected)
	if err != nil {
		fmt.Printf("Error during selection: %v\n", err)
		return
	}

	if len(selected) == 0 {
		fmt.Println("未选择任何应用")
		return
	}

	// 显示确认信息
	fmt.Println("\n将移除以下应用：")
	for _, sel := range selected {
		fmt.Println(sel)
	}

	// 确认操作
	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: "确认移除这些应用吗？",
		Default: true,
	}
	survey.AskOne(confirmPrompt, &confirm)

	if !confirm {
		fmt.Println("操作已取消")
		return
	}

	// 移除选中的应用
	for _, sel := range selected {
		// 从选项中提取 bundle ID
		bundleID := strings.Split(strings.Split(sel, "(")[1], ")")[0]
		err := removeDockApp(currentUser.Username, bundleID)
		if err != nil {
			fmt.Printf("Error removing %s: %v\n", bundleID, err)
		}
	}

	// 重启 Dock
	restartDock(currentUser.Username)
	fmt.Println("操作完成")
}

func getDockApps(username string) ([]DockApp, error) {
	cmd := exec.Command("sudo", "-u", username, "defaults", "read", "com.apple.dock", "persistent-apps")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var apps []DockApp
	scanner := bufio.NewScanner(bytes.NewReader(output))
	var currentApp DockApp
	position := 0

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.Contains(line, "bundle-identifier") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 4 {
				currentApp.BundleID = parts[3]
			}
		} else if strings.Contains(line, "file-label") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 4 {
				label := parts[3]
				// 使用正则表达式替换 Unicode 转义序列
				re := regexp.MustCompile(`\\U([0-9a-fA-F]{4})`)
				decodedLabel := re.ReplaceAllStringFunc(label, func(match string) string {
					// 将 \Uxxxx 转换为实际的 Unicode 字符
					hex := match[2:]
					if i, err := strconv.ParseInt(hex, 16, 32); err == nil {
						return string(rune(i))
					}
					return match
				})
				// 移除多余的转义反斜杠
				decodedLabel = strings.ReplaceAll(decodedLabel, "\\", "")
				currentApp.Label = decodedLabel
				currentApp.Position = position
				apps = append(apps, currentApp)
				position++
			}
		}
	}

	return apps, nil
}

func removeDockApp(username, bundleID string) error {
	// 使用 PlistBuddy 删除指定位置的 Dock 图标
	cmd := exec.Command("sudo", "-u", username, "/usr/libexec/PlistBuddy", "-c",
		fmt.Sprintf("Delete persistent-apps:%d", getAppPosition(username, bundleID)),
		fmt.Sprintf("/Users/%s/Library/Preferences/com.apple.dock.plist", username))
	return cmd.Run()
}

func getAppPosition(username, bundleID string) int {
	cmd := exec.Command("sudo", "-u", username, "defaults", "read", "com.apple.dock", "persistent-apps")
	output, err := cmd.Output()
	if err != nil {
		return -1
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	position := 0

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, bundleID) {
			return position
		}
		if strings.Contains(line, "bundle-identifier") {
			position++
		}
	}

	return -1
}

func restartDock(username string) error {
	cmd := exec.Command("sudo", "-u", username, "osascript", "-e", "delay 3", "-e", "tell Application \"Dock\"", "-e", "quit", "-e", "end tell", "-e", "delay 3")
	return cmd.Run()
}

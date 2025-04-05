# macOS Dock Manager

一个用于管理 macOS Dock 栏应用的命令行工具，可以方便地移除不需要的 Dock 图标。本项目提供了两个实现版本：Go 版本和 Shell 版本。

## 功能特性

- 显示当前 Dock 栏中的所有应用
- 支持多选移除应用
- 自动重启 Dock 以应用更改
- 支持 Unicode 字符显示
- 交互式命令行界面

## 系统要求

- macOS 操作系统
- 管理员权限（用于修改 Dock 设置）

### Go 版本额外要求
- Go 1.16 或更高版本

## 安装

### Go 版本安装
1. 确保已安装 Go 环境
2. 克隆仓库：
```bash
git clone https://github.com/yourusername/macos-dock-manager.git
cd macos-dock-manager
```
3. 安装依赖：
```bash
go mod download
```
4. 编译项目：
```bash
make
```

### Shell 版本安装
1. 克隆仓库：
```bash
git clone https://github.com/yourusername/macos-dock-manager.git
cd macos-dock-manager
```
2. 添加执行权限：
```bash
chmod +x docker-mamger-use-sh.sh
```

## 使用方法

### Go 版本
1. 运行程序：
```bash
./dock-manager
```

2. 使用空格键选择要移除的应用
3. 按回车键确认选择
4. 确认是否要移除选中的应用
5. 程序会自动重启 Dock 以应用更改

### Shell 版本
1. 运行脚本：
```bash
./docker-mamger-use-sh.sh
```

2. 输入要删除的应用编号（多个编号用空格分隔）
3. 按回车键确认选择
4. 输入 'y' 确认删除
5. 脚本会自动重启 Dock 以应用更改

## 版本比较

| 特性 | Go 版本 | Shell 版本 |
|------|---------|------------|
| 安装复杂度 | 较高（需要 Go 环境） | 较低（直接运行） |
| 交互方式 | 图形化选择界面 | 命令行输入编号 |
| 执行效率 | 较高 | 较低 |
| 依赖项 | 需要 Go 环境 | 仅需系统自带工具 |

## 注意事项

- 程序需要管理员权限来修改 Dock 设置
- 移除的应用可以通过 Finder 重新添加到 Dock
- 建议在操作前备份重要的 Dock 设置

## 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进这个项目。

## 许可证

MIT License 
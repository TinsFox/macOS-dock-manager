#!/bin/bash

# 获取当前用户
CURRENT_USER=$(whoami)

# 获取 Dock 应用信息并格式化输出
get_dock_apps() {
    local apps=()
    local position=0
    local bundle_id=""
    local label=""
    
    while IFS= read -r line; do
        if [[ $line =~ "bundle-identifier" ]]; then
            bundle_id=$(echo "$line" | sed -E 's/.*"([^"]+)".*/\1/')
        elif [[ $line =~ "file-label" ]]; then
            label=$(echo "$line" | sed -E 's/.*"([^"]+)".*/\1/')
            # 处理 Unicode 转义序列
            label=$(echo -e "$label")
            apps+=("$position: $label ($bundle_id)")
            ((position++))
        fi
    done < <(defaults read com.apple.dock persistent-apps)
    
    printf '%s\n' "${apps[@]}"
}

# 从 Dock 中移除应用
remove_dock_app() {
    local position=$1
    /usr/libexec/PlistBuddy -c "Delete persistent-apps:$position" ~/Library/Preferences/com.apple.dock.plist
}

# 重启 Dock
restart_dock() {
    osascript -e 'tell application "Dock" to quit'
    sleep 2
}

# 主程序
main() {
    # 获取并显示所有 Dock 应用
    echo "当前 Dock 中的应用:"
    echo "==================="
    apps=($(get_dock_apps))
    
    # 显示应用列表
    for i in "${!apps[@]}"; do
        echo "$((i+1)). ${apps[$i]}"
    done
    
    # 提示用户选择要删除的应用
    echo -e "\n请输入要删除的应用编号（多个编号用空格分隔，直接回车确认）:"
    read -r selections
    
    if [ -z "$selections" ]; then
        echo "未选择任何应用"
        exit 0
    fi
    
    # 确认删除
    echo -e "\n将删除以下应用："
    for num in $selections; do
        echo "${apps[$((num-1))]}"
    done
    
    echo -e "\n确认删除这些应用吗？(y/N)"
    read -r confirm
    
    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        echo "操作已取消"
        exit 0
    fi
    
    # 执行删除操作
    # 注意：需要倒序删除以避免位置改变影响后续删除
    for num in $(echo "$selections" | tr ' ' '\n' | sort -nr); do
        position=$((num-1))
        echo "正在删除位置 $position 的应用..."
        remove_dock_app "$position"
    done
    
    # 重启 Dock
    echo "重启 Dock..."
    restart_dock
    
    echo "操作完成"
}

# 运行主程序
main


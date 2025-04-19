#!/bin/bash

# 设置记录文件的路径
LOGFILE="docker_usage_log.txt"

# 每秒钟记录一次Docker CPU和内存使用情况
while true; do
    # 获取当前时间
    TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
    
    # 获取所有正在运行的Docker容器的CPU和内存使用情况
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" | while read -r line ; do
        # 如果是表头跳过
        if [[ $line == "NAME"* ]]; then
            continue
        fi
        
        # 将结果写入日志文件
        echo "$TIMESTAMP $line" >> $LOGFILE
    done
    
    # 等待1秒
    sleep 1
done

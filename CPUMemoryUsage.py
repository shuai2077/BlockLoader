#encoding=utf-8
"""
@author=gang wang
"""
import matplotlib.pyplot as plt
from matplotlib import rc
rc('font',**{'family':'serif','serif':['Palatino']})
plt.rcParams['pdf.fonttype'] = 42


def parse_log_file(file_path):
    cpu_usage = []
    memory_usage = []
    timestamps = []

    with open(file_path, 'r') as file:
        lines = file.readlines()

        for i in range(1, len(lines), 2):  # 读取偶数行
            line = lines[i].strip()
            if 'peer0.org2.example.com' in line:
                parts = line.split()
                try:
                    timestamps.append(parts[0] + ' ' + parts[1])  # 提取时间戳
                    cpu_value = parts[3].replace('%', '').strip()
                    memory_value = parts[4].replace('MiB', '').strip()
                    cpu_usage.append(float(cpu_value))  # 提取CPU使用率
                    memory_usage.append(float(memory_value))  # 提取内存使用率（以MiB为单位）
                except ValueError as e:
                    print(f"Error converting values: {e}")
                    continue

    return timestamps, cpu_usage, memory_usage


def plot_metrics(timestamps, cpu_usage, memory_usage):
    #plt.figure(figsize=(10, 6))
    x_values = list(range(1, len(cpu_usage) + 1))  # 将x轴设置为次数

    plt.figure(figsize=(10, 6))

    plt.plot(x_values, cpu_usage, marker='d', color='#FA7F6F', label='CPU Usage (%)')
    #plt.title('CPU Usage Over Time')
    plt.xlabel('Count',fontsize=20)
    plt.ylabel('CPU Usage (%)', fontsize=20)
    plt.xticks(fontsize=15)  # 设置x轴刻度字体大小
    plt.yticks(fontsize=15)
    plt.grid(True)
    plt.show()

    plt.figure(figsize=(10, 6))

    plt.plot(x_values, memory_usage, marker='p', color='#FFBE7A', label='Memory Usage (MiB)')
    #plt.title('Memory Usage Over Time')
    plt.xlabel('Count' ,fontsize=20)
    plt.ylabel('Memory Usage (MB)', fontsize=20)
    plt.xticks(fontsize=15)  # 设置x轴刻度字体大小
    plt.yticks(fontsize=15)
    plt.grid(True)
    plt.show()

    # plt.subplot(2, 1, 1)
    # plt.plot(timestamps, cpu_usage, marker='o', color='b', label='CPU Usage (%)')
    # plt.title('CPU and Memory Usage Over Time')
    # plt.ylabel('CPU Usage (%)')
    # plt.xticks(rotation=45)
    # plt.grid(True)
    #
    # plt.subplot(2, 1, 2)
    # plt.plot(timestamps, memory_usage, marker='o', color='r', label='Memory Usage (MiB)')
    # plt.xlabel('Time')
    # plt.ylabel('Memory Usage (MiB)')
    # plt.xticks(rotation=45)
    # plt.grid(True)
    #
    # plt.tight_layout()
    # plt.show()


# 使用示例
file_path = 'F:/WeChat Files/wxid_35yudjrz27xy22/FileStorage/File/2024-08/linear-rate/docker_usage_log.txt'   # 替换为你的文件路径
timestamps, cpu_usage, memory_usage = parse_log_file(file_path)
plot_metrics(timestamps, cpu_usage, memory_usage)

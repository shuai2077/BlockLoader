#encoding=utf-8
"""
@author=gang wang
"""
import matplotlib.pyplot as plt
from datetime import datetime
import pandas as pd


# 读取文件并解析数据
def parse_transaction_data(file_path):
    with open(file_path, 'r') as file:
        lines = file.readlines()

    transactions = []
    for i in range(0, len(lines), 2):
        transaction_id = lines[i].strip().split(': ')[1]
        timestamp = lines[i + 1].strip().split(': ')[1]
        transactions.append((transaction_id, timestamp))
    return transactions


# 计算TPS
def calculate_tps(transactions):
    time_format = '%Y-%m-%dT%H:%M:%S%z'
    timestamps = [datetime.strptime(t[1], time_format) for t in transactions]

    tps_dict = {}
    for timestamp in timestamps:
        time_key = timestamp.replace(microsecond=0)  # 只舍入到秒
        if time_key not in tps_dict:
            tps_dict[time_key] = 0
        tps_dict[time_key] += 1

    tps_series = pd.Series(tps_dict)
    return tps_series

    # time_format = '%Y-%m-%dT%H:%M:%S%z'
    # timestamps = [datetime.strptime(t[1], time_format) for t in transactions]
    #
    # tps_dict = {}
    # for timestamp in timestamps:
    #     time_key = timestamp.replace(second=0, microsecond=0)
    #     if time_key not in tps_dict:
    #         tps_dict[time_key] = 0
    #     tps_dict[time_key] += 1
    #
    # tps_series = pd.Series(tps_dict)
    # return tps_series


# 绘制TPS折线图
def plot_tps(tps_series):
    tps_series.plot(kind='line', title='TPS over Time', xlabel='Time', ylabel='Transactions per Second')
    plt.show()


# 计算统计信息
def calculate_statistics(tps_series):
    max_tps = tps_series.max()
    min_tps = tps_series.min()
    avg_tps = tps_series.mean()
    return max_tps, min_tps, avg_tps


# 主函数
def main(file_path):
    transactions = parse_transaction_data(file_path)
    tps_series = calculate_tps(transactions)
    plot_tps(tps_series)

    max_tps, min_tps, avg_tps = calculate_statistics(tps_series)
    print(f"Max TPS: {max_tps}")
    print(f"Min TPS: {min_tps}")
    print(f"Average TPS: {avg_tps:.2f}")


# 示例用法
file_path = 'F:/WeChat Files/wxid_35yudjrz27xy22/FileStorage/File/2024-08/data_v2/transaction_info.txt'  # 替换为你的文件路径
main(file_path)

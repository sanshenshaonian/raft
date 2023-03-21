import re
import os
import time

# 用于存储每个分组的最后5条日志和状态
group_data = {}

# 用于匹配方括号内的数字和follower/leader/client/kvserver
pattern = re.compile(r"\[(\d+)\].*(follower|leader|client|kvserver)?")

# 遍历日志文件
with open('d:/Code/go/src/新建文件夹/raft-go/src/kvraft/cs_logger.txt', 'r', encoding='utf-8') as f:
    for line in f:
        # 匹配方括号内的数字和follower/leader/client/kvserver
        match = pattern.search(line)
        if match:
            # 获取分组名称
            if 'client' in line:
                group_name = 'client'
            elif 'kvserver' in line:
                group_name = 'kvserver'
            else:
                # 如果没有匹配到client和kvserver，则使用方括号内数字作为分组名称
                group_number = int(match.group(1))
                group_name = f"group{group_number}"

            # 获取分组状态（如果有）
            status = match.group(2)

            # 更新分组状态（如果有）
            if status and group_name not in group_data:
                group_data[group_name] = {'logs': [], 'status': ''}
            if status:
                group_data[group_name]['status'] = status

            # 将日志添加到分组中
            if group_name not in group_data:
                group_data[group_name] = {'logs': [], 'status': ''}
            group_data[group_name]['logs'].append(line.strip())

            # 只保留最后5条日志
            if len(group_data[group_name]['logs']) > 10:
                group_data[group_name]['logs'] = group_data[group_name]['logs'][-10:]

            # 输出每个分组的状态和最后5条日志
            os.system("cls")
            for group_name, data in group_data.items():
                print(f"{group_name} ({data['status']}):")
                # 对于不是client和kvserver的分组，使用方括号内数字进行二次分组
                if 'client' not in group_name and 'kvserver' not in group_name:
                    group_number = int(re.findall(r"\[(\d+)\]", data['logs'][0])[0])
                    for log in data['logs']:
                        if int(re.findall(r"\[(\d+)\]", log)[0]) == group_number:
                            print(f"  {log}")
                # 对于client和kvserver分组，直接输出日志
                else:
                    for log in data['logs']:
                        print(f"  {log}")
            time.sleep(0.5)

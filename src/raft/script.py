import re
import os
import time

# 用于存储每个分组的最后5条日志和状态
group_data = {}

# 用于匹配方括号内的数字和follower/leader
pattern = re.compile(r"\[(\d+)\].*(follower|leader)?")

# 遍历日志文件
with open('d:/Code/go/src/新建文件夹/raft-go/src/raft/raft_logger.txt','r',encoding='utf-8') as f:
    for line in f:
        # 匹配方括号内的数字和follower/leader
        match = pattern.search(line)
        if match:
            # 获取分组编号和状态（如果有）
            group_number = int(match.group(1))
            status = match.group(2)

            # 更新分组状态（如果有）
            if status and group_number not in group_data:
                group_data[group_number] = {'logs': [], 'status': ''}
            if status:
                group_data[group_number]['status'] = status

            # 将日志添加到分组中
            if group_number not in group_data:
                group_data[group_number] = {'logs': [], 'status': ''}
            group_data[group_number]['logs'].append(line.strip())

            # 只保留最后5条日志
            if len(group_data[group_number]['logs']) > 10:
                group_data[group_number]['logs'] = group_data[group_number]['logs'][-10:]

            # 输出每个分组的状态和最后5条日志
            os.system("cls")
            for group_number, data in group_data.items():
                print(f"Group {group_number} ({data['status']}):")
                for log in data['logs']:
                    print(f"  {log}")
            time.sleep(0.5)


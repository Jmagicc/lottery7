import requests
from bs4 import BeautifulSoup
from datetime import datetime
import pymysql
from pymysql import Error

def get_latest_draw():
    """抓取最新开奖信息"""
    url = "https://www.vipc.cn/results/qxc?in=result_content"
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
    }
    
    try:
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        
        soup = BeautifulSoup(response.text, 'html.parser')
        
        # 获取开奖信息
        draw_info_div = soup.select_one('.vResult_listDigit_title')
        if not draw_info_div:
            raise Exception("未找到开奖信息区域")
        
        # 获取期号
        period_text = draw_info_div.find('i').text.strip()  # "七星彩开奖结果第24151期"
        draw_no = period_text.split('第')[1].replace('期', '')  # 提取期号 "24151"
        
        # 获取开奖时间
        date_text = draw_info_div.find('span').text.strip()  # "开奖时间: 2024-12-31 20:30:00"
        draw_date = datetime.strptime(date_text.split(': ')[1][:10], '%Y-%m-%d').date()
        
        # 获取开奖号码
        number_elements = soup.select('.vRes_lottery_ball b.red')[:4]  # 只取前4个数字
        if len(number_elements) < 4:
            raise Exception(f"开奖号码不完整，仅找到 {len(number_elements)} 个号码")
        
        num_list = [int(num.text.strip()) for num in number_elements]
        
        result = {
            'draw_date': draw_date,
            'draw_no': draw_no,
            'numbers': num_list
        }
        
        print(f"抓取到的数据：期号={draw_no}, 日期={draw_date}, 号码={num_list}")
        return result
        
    except Exception as e:
        print(f"抓取数据时出错: {e}")
        print(f"请求URL: {url}")
        return None

def create_table_if_not_exists(cursor):
    """Create the lottery_results table if it doesn't exist"""
    create_table_query = """
    CREATE TABLE IF NOT EXISTS lottery_results (
        id bigint NOT NULL AUTO_INCREMENT,
        draw_no varchar(20) NOT NULL COMMENT '期号',
        draw_date date NOT NULL COMMENT '开奖日期',
        num1 tinyint unsigned NOT NULL COMMENT '第1位数字',
        num2 tinyint unsigned NOT NULL COMMENT '第2位数字',
        num3 tinyint unsigned NOT NULL COMMENT '第3位数字',
        num4 tinyint unsigned NOT NULL COMMENT '第4位数字',
        created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (id),
        UNIQUE KEY draw_no (draw_no),
        KEY idx_draw_no (draw_no),
        KEY idx_draw_date (draw_date)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='七星彩开奖结果';
    """
    cursor.execute(create_table_query)

def main():
    # 获取当前日期
    current_date = datetime.now().date()
    print(f"当前日期: {current_date}")
    
    # 获取最新开奖信息
    draw_info = get_latest_draw()
    
    if not draw_info:
        print("获取开奖信息失败")
        return
        
    if draw_info['draw_date'] != current_date:
        print(f"今日暂无新开奖结果（最新开奖日期: {draw_info['draw_date']}）")
        return
    
    # 数据库配置
    db_config = {
        'host': '192.168.0.200',
        'port': 3306,
        'user': 'lottery7',
        'password': 'EniQXpY6x8rMjMsz',
        'database': 'lottery7',
        'charset': 'utf8mb4'
    }
    
    try:
        # 建立数据库连接
        connection = pymysql.connect(**db_config)
        cursor = connection.cursor()
        
        # 创建表（如果不存在）
        create_table_if_not_exists(cursor)
        
        # 插入开奖结果
        insert_query = """
        INSERT INTO lottery_results 
        (draw_no, draw_date, num1, num2, num3, num4)
        VALUES (%s, %s, %s, %s, %s, %s)
        """
        
        values = (
            draw_info['draw_no'],
            draw_info['draw_date'],
            draw_info['numbers'][0],
            draw_info['numbers'][1],
            draw_info['numbers'][2],
            draw_info['numbers'][3]
        )
        
        cursor.execute(insert_query, values)
        connection.commit()
        print(f"成功插入开奖结果：期号={draw_info['draw_no']}, 日期={draw_info['draw_date']}, 号码={draw_info['numbers']}")
        
    except Error as e:
        print(f"数据库错误: {e}")
    finally:
        if 'connection' in locals():  # 如果连接存在
            cursor.close()
            connection.close()
            print("数据库连接已关闭")

if __name__ == "__main__":
    main() 
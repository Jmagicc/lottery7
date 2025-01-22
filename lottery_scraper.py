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
        number_elements = soup.select('.vRes_lottery_ball b.red')[:5]  # 只取前5个数字
        if len(number_elements) != 5:
            raise Exception(f"开奖号码数量不正确，期望5个但找到 {len(number_elements)} 个号码")

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
        num5 tinyint unsigned NOT NULL COMMENT '第5位数字',
        created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (id),
        UNIQUE KEY draw_no (draw_no),
        KEY idx_draw_no (draw_no),
        KEY idx_draw_date (draw_date)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='七星彩开奖结果';
    """
    cursor.execute(create_table_query)

def get_latest_draw_from_db(cursor):
    """获取数据库中最新的开奖期号"""
    query = """
    SELECT draw_no 
    FROM lottery_results 
    ORDER BY draw_no DESC 
    LIMIT 1
    """
    cursor.execute(query)
    result = cursor.fetchone()
    return result[0] if result else None

def get_draw_by_period(period):
    """抓取指定期号的开奖信息"""
    url = f"https://www.vipc.cn/results/qxc/{period}?in=result_content"
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
            return None

        # 获取期号
        period_text = draw_info_div.find('i').text.strip()
        draw_no = period_text.split('第')[1].replace('期', '')

        # 如果返回的期号与请求的不一致，说明该期尚未开奖
        if draw_no != str(period):
            return None

        # 获取开奖时间
        date_text = draw_info_div.find('span').text.strip()
        draw_date = datetime.strptime(date_text.split(': ')[1][:10], '%Y-%m-%d').date()

        # 获取开奖号码
        number_elements = soup.select('.vRes_lottery_ball b.red')[:5]
        if len(number_elements) != 5:
            return None

        num_list = [int(num.text.strip()) for num in number_elements]

        return {
            'draw_date': draw_date,
            'draw_no': draw_no,
            'numbers': num_list
        }

    except Exception as e:
        print(f"抓取期号 {period} 数据时出错: {e}")
        return None

def main():
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

        # 获取数据库中最新的期号
        latest_db_draw = get_latest_draw_from_db(cursor)

        # 获取网站最新开奖信息
        latest_web_draw = get_latest_draw()
        if not latest_web_draw:
            print("获取网站最新开奖信息失败")
            return

        latest_web_period = int(latest_web_draw['draw_no'])
        start_period = int(latest_db_draw) + 1 if latest_db_draw else latest_web_period

        print(f"开始同步期号范围: {start_period} - {latest_web_period}")

        sync_count = 0
        # 从数据库最新期号的下一期开始，直到网站最新期号
        for period in range(start_period, latest_web_period + 1):
            draw_info = get_draw_by_period(period)
            if not draw_info:
                continue

            insert_query = """
            INSERT IGNORE INTO lottery_results 
            (draw_no, draw_date, num1, num2, num3, num4, num5)
            VALUES (%s, %s, %s, %s, %s, %s, %s)
            """

            values = (
                draw_info['draw_no'],
                draw_info['draw_date'],
                draw_info['numbers'][0],
                draw_info['numbers'][1],
                draw_info['numbers'][2],
                draw_info['numbers'][3],
                draw_info['numbers'][4]
            )

            result = cursor.execute(insert_query, values)
            if result == 1:
                sync_count += 1
                print(f"成功同步期号 {draw_info['draw_no']}: 日期={draw_info['draw_date']}, 号码={draw_info['numbers']}")
                connection.commit()

        print(f"\n同步完成！共同步了 {sync_count} 期开奖数据")

    except pymysql.err.OperationalError as e:
        print(f"数据库连接错误: {e}")
    except pymysql.err.IntegrityError as e:
        print(f"数据完整性错误（可能是重复数据）: {e}")
    except Error as e:
        print(f"其他数据库错误: {e}")
    finally:
        if 'connection' in locals():
            cursor.close()
            connection.close()
            print("数据库连接已关闭")

if __name__ == "__main__":
    main()

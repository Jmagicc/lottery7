七星彩数据采集工具 - 自动获取并存储七星彩开奖结果到MySQL数据库

## 功能特点

- 自动抓取七星彩开奖数据（前4位数字）
- MySQL数据库存储支持
- 自动判断开奖日期，避免重复采集
- 详细的运行日志
- 使用Poetry进行依赖管理
- 支持编译为独立二进制文件


## 生成密钥
```
 go test -v ./service -run TestGenerateKey
```

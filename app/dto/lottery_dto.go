package dto

// LotteryResultResponse 开奖结果响应
type LotteryResultResponse struct {
	DrawNo   string `json:"draw_no"`
	DrawDate string `json:"draw_date"`
	DrawDay  string `json:"draw_day"`
	Num1     uint8  `json:"num1"`
	Num2     uint8  `json:"num2"`
	Num3     uint8  `json:"num3"`
	Num4     uint8  `json:"num4"`
	Num5     uint8  `json:"num5"`
}

// UnusedNumbersResponse 未使用号码响应
type UnusedNumbersResponse struct {
	UnusedNumbers []string `json:"unusedNumbers"`
	Count         int      `json:"count"`
}

// MatrixResponse 矩阵响应
type MatrixResponse struct {
	Matrix []string `json:"matrix"` // 开奖号码矩阵
	Prompt string   `json:"prompt"` // AI提示文本
}

// DeepSeekMessage AI消息结构
type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeekRequest AI请求结构
type DeepSeekRequest struct {
	Model    string            `json:"model"`
	Messages []DeepSeekMessage `json:"messages"`
	Stream   bool              `json:"stream"`
}

// DeepSeekResponse AI响应结构
type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// RepeatNumbersResponse 重复号码统计响应
type RepeatNumbersResponse struct {
	DoubleNumbers []RepeatNumberStat `json:"doubleNumbers"` // 双重号码统计
	TripleNumbers []RepeatNumberStat `json:"tripleNumbers"` // 三重号码统计
	TotalPeriods  int                `json:"totalPeriods"`  // 总期数
}

// RepeatNumberStat 重复号码统计
type RepeatNumberStat struct {
	Number    string  `json:"number"`    // 重复号码
	Count     int     `json:"count"`     // 出现次数
	Frequency float64 `json:"frequency"` // 出现频率
}

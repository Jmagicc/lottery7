package service

import (
	"fmt"
	"lottery7/dto"
	"lottery7/models"
	"math/rand"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

type LotteryService struct {
	db *gorm.DB
}

func NewLotteryService(db *gorm.DB) *LotteryService {
	return &LotteryService{db: db}
}

func (s *LotteryService) GetLotteryResults() ([]dto.LotteryResultResponse, error) {
	var results []models.LotteryResult
	if err := s.db.Order("draw_date desc").Limit(100).Find(&results).Error; err != nil {
		return nil, err
	}

	var response []dto.LotteryResultResponse
	for _, result := range results {
		response = append(response, dto.LotteryResultResponse{
			DrawNo:   result.DrawNo,
			DrawDate: result.DrawDate.Format("2006-01-02"),
			DrawDay:  result.DrawDate.Weekday().String(),
			Num1:     result.Num1,
			Num2:     result.Num2,
			Num3:     result.Num3,
			Num4:     result.Num4,
			Num5:     result.Num5,
		})
	}

	return response, nil
}

func (s *LotteryService) GetUnusedNumbers(periods int) (*dto.UnusedNumbersResponse, error) {
	var results []struct {
		Num1Num2 string `gorm:"column:num1_num2"`
	}

	err := s.db.Raw(`
		SELECT 
			CONCAT(num1, num2) AS num1_num2
		FROM 
			lottery_results
		ORDER BY 
			draw_date DESC
		LIMIT ?
	`, periods).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	usedNumbers := make(map[string]bool)
	for _, result := range results {
		usedNumbers[result.Num1Num2] = true
	}

	var unusedNumbers []string
	for i := 0; i < 100; i++ {
		num := fmt.Sprintf("%02d", i)
		if !usedNumbers[num] {
			num += "xx"
			unusedNumbers = append(unusedNumbers, num)
		}
	}

	return &dto.UnusedNumbersResponse{
		UnusedNumbers: unusedNumbers,
		Count:         len(unusedNumbers),
	}, nil
}

func (s *LotteryService) GetNumberMatrix() (*dto.MatrixResponse, error) {
	var results []models.LotteryResult
	if err := s.db.Order("draw_date desc").Limit(6).Find(&results).Error; err != nil {
		return nil, err
	}

	var matrix []string
	for _, result := range results {
		row := fmt.Sprintf("%d %d %d %d %d",
			result.Num1, result.Num2, result.Num3, result.Num4, result.Num5)
		matrix = append(matrix, row)
	}
	// 翻转 matrix
	for i, j := 0, len(matrix)-1; i < j; i, j = i+1, j-1 {
		matrix[i], matrix[j] = matrix[j], matrix[i]
	}
	matrix = append(matrix, "? ? ? ? ?")
	matrixStr := strings.Join(matrix, "\n")

	//var tripletCounts []struct {
	//	NumTriplet   string `gorm:"column:num_triplet"`
	//	TripletCount int    `gorm:"column:triplet_count"`
	//}
	//if err := s.db.Raw(`
	//	SELECT num_triplet, triplet_count
	//	FROM lottery_number_triplet_counts
	//	WHERE triplet_count >= 1
	//`).Scan(&tripletCounts).Error; err != nil {
	//	return nil, err
	//}
	//
	//var tripletResults []string
	//for _, tc := range tripletCounts {
	//	tripletResults = append(tripletResults, tc.NumTriplet)
	//}
	//
	//var tripletStr strings.Builder
	//for _, triplet := range tripletResults {
	//	tripletStr.WriteString(triplet + "\n")
	//}

	return &dto.MatrixResponse{
		Matrix: matrix,
		Prompt: fmt.Sprintf(`[%s] 作为一位在%d年代成长起来的数学天才，你现在是世界顶级数学家和概率学家，擅长从复杂数据中识别隐藏规律。今天，你需要以全新的视角来分析一个特殊的数字矩阵。

		请以你独特的思维方式，仔细观察以下矩阵。这个矩阵的每一行都包含五个数字（0-9）：

		%s



		基于当前时间点[%s]的分析，请从以下几个方面进行思考：
		1. 这个矩阵中蕴含着怎样独特的数字排列规律？
		2. 根据你的专业判断，问号处最有可能出现的前三位定组合是什么？请确保给出具有唯一性的答案
		3. 请跳出固有思维模式，提供一个全新的分析视角，不要受到任何历史分析的影响

		注意：这是在%s时刻进行的分析，请确保你的回答反映当前时间点的独特见解。
		`,
			time.Now().Format("2006-01-02 15:04:05"),
			rand.Intn(3)+198, // 随机生成年代：1980-1982
			matrixStr,
			//tripletStr.String(),
			time.Now().Format("15:04:05"),
			time.Now().Format("2006-01-02 15:04:05")),
	}, nil
}

func (s *LotteryService) GetRepeatNumbers() (*dto.RepeatNumbersResponse, error) {
	var results []models.LotteryResult
	if err := s.db.Order("draw_date desc").Limit(200).Find(&results).Error; err != nil {
		return nil, err
	}

	doubleStats := make(map[string]int)
	tripleStats := make(map[string]int)

	for _, result := range results {
		// 创建数字计数map
		numCount := make(map[int]int)

		// 统计每个数字出现的次数
		nums := []int{int(result.Num1), int(result.Num2), int(result.Num3), int(result.Num4)}
		for _, num := range nums {
			numCount[num]++
		}

		// 检查重复
		for num, count := range numCount {
			if count == 2 {
				// 双重数：将数字转为两位相同数字的字符串
				doubleNum := fmt.Sprintf("%d%d", num, num)
				doubleStats[doubleNum]++
			}
			if count == 3 {
				// 三重数：将数字转为三位相同数字的字符串
				tripleNum := fmt.Sprintf("%d%d%d", num, num, num)
				tripleStats[tripleNum]++
			}
		}
	}

	totalPeriods := len(results)
	var response dto.RepeatNumbersResponse
	response.TotalPeriods = totalPeriods

	for num, count := range doubleStats {
		response.DoubleNumbers = append(response.DoubleNumbers, dto.RepeatNumberStat{
			Number:    num,
			Count:     count,
			Frequency: float64(count) / float64(totalPeriods),
		})
	}

	for num, count := range tripleStats {
		response.TripleNumbers = append(response.TripleNumbers, dto.RepeatNumberStat{
			Number:    num,
			Count:     count,
			Frequency: float64(count) / float64(totalPeriods),
		})
	}

	if response.TripleNumbers == nil {
		response.TripleNumbers = make([]dto.RepeatNumberStat, 0)
	}

	sort.Slice(response.DoubleNumbers, func(i, j int) bool {
		return response.DoubleNumbers[i].Number < response.DoubleNumbers[j].Number
	})

	sort.Slice(response.TripleNumbers, func(i, j int) bool {
		return response.TripleNumbers[i].Number < response.TripleNumbers[j].Number
	})

	return &response, nil
}

// reverseMatrix 翻转矩阵
func reverseMatrix(matrix [][]string) [][]string {
	// 获取矩阵的行数和列数
	rows := len(matrix)
	if rows == 0 {
		return matrix
	}
	cols := len(matrix[0])

	// 创建一个新的矩阵来存储翻转后的结果
	reversed := make([][]string, cols)
	for i := range reversed {
		reversed[i] = make([]string, rows)
	}

	// 进行翻转操作
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			reversed[j][rows-1-i] = matrix[i][j]
		}
	}

	return reversed
}

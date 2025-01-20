package service

import (
	"fmt"
	"lottery7/dto"
	"lottery7/models"
	"sort"
	"strings"

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
	matrix = append(matrix, "? ? ? ? ?")

	var tripletCounts []struct {
		NumTriplet   string `gorm:"column:num_triplet"`
		TripletCount int    `gorm:"column:triplet_count"`
	}

	if err := s.db.Raw(`
		SELECT num_triplet, triplet_count
		FROM lottery_number_triplet_counts
		WHERE triplet_count >= 1
	`).Scan(&tripletCounts).Error; err != nil {
		return nil, err
	}

	var tripletResults []string
	for _, tc := range tripletCounts {
		tripletResults = append(tripletResults, tc.NumTriplet)
	}

	matrixStr := strings.Join(matrix, "\n")
	var tripletStr strings.Builder
	for _, triplet := range tripletResults {
		tripletStr.WriteString(triplet + "\n")
	}

	return &dto.MatrixResponse{
		Matrix: matrix,
		Prompt: fmt.Sprintf(`你是一位顶级数学家和概率学家，擅长从复杂数据中识别隐藏规律。现在，你面对一个特殊的矩阵排列，其元素看似随机但实则遵循某种深层规律。

		请仔细观察以下矩阵，并运用你的专业知识和洞察力，预测出其中的规律。这个矩阵每一行有五位，每一位都是0-9：

		%s

		历史上出现过的三定组合如下：

		%s

		基于以上数据，请分析：
		1. 矩阵中的数字排列规律
		2. 算出问号中有可能出现的前三位的定组合
		`, matrixStr, tripletStr.String()),
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

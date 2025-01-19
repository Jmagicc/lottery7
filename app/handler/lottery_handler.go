package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"lottery7/crypto"
	"lottery7/dto"
	"lottery7/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type LotteryHandler struct {
	service *service.LotteryService
}

func NewLotteryHandler(service *service.LotteryService) *LotteryHandler {
	return &LotteryHandler{service: service}
}

func (h *LotteryHandler) GetLotteryResults(c *gin.Context) {
	results, err := h.service.GetLotteryResults()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *LotteryHandler) GetUnusedNumbers(c *gin.Context) {
	periodsStr := c.Query("periods")
	if periodsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing periods parameter"})
		return
	}

	periods, err := strconv.Atoi(periodsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid periods parameter"})
		return
	}

	response, err := h.service.GetUnusedNumbers(periods)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *LotteryHandler) GetNumberMatrix(c *gin.Context) {
	matrixResponse, err := h.service.GetNumberMatrix()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 从配置文件读取 AI API key
	apiKey := viper.GetString("deepseek.api_key")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DeepSeek API key not configured"})
		return
	}
	apiKey, err = crypto.AESDecrypt(apiKey, "testkey123456")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := "请你根据这个矩阵，找规律，要求只算出前三个?，大概率的组成会有哪些。请从历史上出现过的三定组合的范围中给出答案。不需要你说明为什么。请注意，只告诉我10个结果，并用逗号拼接。"

	aiResponse, err := h.callDeepSeekAPI(apiKey, matrixResponse.Prompt, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 处理 AI 响应，将逗号分隔的结果转换为数组
	predictions := strings.Split(aiResponse, ",")

	// 清理每个预测结果并添加 "x"
	cleanedPredictions := make([]string, 0, len(predictions))
	for _, p := range predictions {
		// 清理空格和其他无关字符
		p = strings.TrimSpace(p)
		if p != "" {
			// 确保每个预测后面都加上 "x"
			if !strings.HasSuffix(p, "x") {
				p = p + "x"
			}
			cleanedPredictions = append(cleanedPredictions, p)
		}
	}

	// 返回处理后的结果
	c.JSON(http.StatusOK, gin.H{
		"matrix":      matrixResponse.Matrix,
		"predictions": cleanedPredictions,
	})
}

func (h *LotteryHandler) callDeepSeekAPI(apiKey, prompt, query string) (string, error) {
	messages := []dto.DeepSeekMessage{
		{Role: "system", Content: prompt},
		{Role: "user", Content: query},
	}

	requestBody := dto.DeepSeekRequest{
		Model:    "deepseek-chat",
		Messages: messages,
		Stream:   false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var apiResponse dto.DeepSeekResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", err
	}

	if len(apiResponse.Choices) > 0 {
		return apiResponse.Choices[0].Message.Content, nil
	}

	return "", nil
}

func (h *LotteryHandler) GetRepeatNumbers(c *gin.Context) {
	response, err := h.service.GetRepeatNumbers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

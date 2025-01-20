package handler

import (
	"lottery7/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LicenseHandler struct {
	service *service.LicenseService
}

func NewLicenseHandler(service *service.LicenseService) *LicenseHandler {
	return &LicenseHandler{service: service}
}

func (h *LicenseHandler) ValidateKey(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效密钥"})
		return
	}

	createdAt, err := h.service.ValidateKey(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "无效密钥"})
		return
	}

	// 计算剩余天数
	expiryDate := createdAt.AddDate(0, 0, 30)
	remainingDays := int(time.Until(expiryDate).Hours() / 24)

	c.JSON(http.StatusOK, gin.H{
		"created_at":     createdAt,
		"expiry_date":    expiryDate,
		"remaining_days": remainingDays,
	})
}

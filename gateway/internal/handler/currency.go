package handler

import (
	"net/http"
	"time"

	"github.com/vctrl/currency-service/gateway/internal/dto"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetCurrencyRates(c *gin.Context) {
	var req dto.CurrencyRequest
	err := c.BindQuery(&req)
	if err != nil {
		s.logger.Error("Error binding request parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dateFrom, err := time.Parse("2006-01-02", req.DateFrom)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for date_from, expected YYYY-MM-DD"})
		return
	}

	dateTo, err := time.Parse("2006-01-02", req.DateTo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for date_to, expected YYYY-MM-DD"})
		return
	}

	parsedCurrencyRequest := dto.ParsedCurrencyRequest{
		Currency: req.Currency,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	data, err := s.CurrencyService.GetCurrencyRates(c.Request.Context(), parsedCurrencyRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

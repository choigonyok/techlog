package handler

import (
	"encoding/json"
	"fmt"

	"github.com/choigonyok/techlog/internal/usecase"
	"github.com/gin-gonic/gin"
)

type VisitorHandler struct {
	gin.HandlerFunc
	usecase *usecase.VisitorUsecase
}

func NewVisitorHandler() *VisitorHandler {
	return &VisitorHandler{
		usecase: usecase.NewVisitorUsecase(),
	}
}

// GetVisitorCounts returns today/total visitor counts
func (h *VisitorHandler) GetVisitorCounts(c *gin.Context) {
	today, err := h.usecase.GetTodayCount()
	if err != nil {
		fmt.Println("ERR GETTING TODAY COUNT:", err)
		// RETURN
		return
	}
	total, err := h.usecase.GetTotalCount()
	if err != nil {
		fmt.Println("ERR GETTOTALCOUNT:", err)
		// RETURN
		return
	}

	data := struct {
		Today int `json:"today"`
		Total int `json:"total"`
	}{
		Today: today,
		Total: total,
	}

	b, _ := json.Marshal(data)

	c.Writer.Write(b)
}

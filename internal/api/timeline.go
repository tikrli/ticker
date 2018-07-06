package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	. "git.codecoop.org/systemli/ticker/internal/storage"
	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/util"
)

//GetTimelineHandler returns the public timeline for a ticker.
func GetTimelineHandler(c *gin.Context) {
	domain, err := GetDomain(c)
	if err != nil {
		c.JSON(http.StatusOK, JSONResponse{
			Data:   map[string]interface{}{"messages": nil},
			Status: ResponseError,
			Error: map[string]interface{}{
				"code":    ErrorCodeDefault,
				"message": `Could not find a ticker.`,
			},
		})
		return
	}

	ticker, err := FindTicker(domain)
	if err != nil {
		c.JSON(http.StatusOK, JSONResponse{
			Data:   map[string]interface{}{"messages": nil},
			Status: ResponseError,
			Error: map[string]interface{}{
				"code":    ErrorCodeDefault,
				"message": `Could not find a ticker.`,
			},
		})
		return
	}

	pagination := NewPagination(c)
	messages, err := FindByTicker(ticker, pagination)

	c.JSON(http.StatusOK, JSONResponse{
		Data:   map[string]interface{}{"messages": NewMessagesResponse(messages)},
		Status: ResponseSuccess,
		Error:  nil,
	})
	return
}

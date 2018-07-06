package api

import (
	"net/http"
	"github.com/gin-gonic/gin"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)

//Settings holds options for frontend settings
type Settings struct {
	RefreshInterval int `json:"refresh_interval,omitempty"`
}

//GetInitHandler returns the basic settings for the ticker.
func GetInitHandler(c *gin.Context) {
	domain, err := GetDomain(c)

	settings := Settings{
		RefreshInterval: 10,
	}

	ticker, err := FindTicker(domain)
	if err != nil {
		c.JSON(http.StatusOK, JSONResponse{
			Data:   map[string]interface{}{"ticker": nil, "settings": settings},
			Status: ResponseSuccess,
			Error:  nil,
		})
		return
	}

	c.JSON(http.StatusOK, JSONResponse{
		//TODO: Build NewTickerPublicResponse to hide unnecessary information
		Data:   map[string]interface{}{"ticker": NewTickerResponse(ticker), "settings": settings},
		Status: ResponseSuccess,
		Error:  nil,
	})
	return
}

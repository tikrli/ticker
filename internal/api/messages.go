package api

import (
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"git.codecoop.org/systemli/ticker/internal/bridge"
	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)

//GetMessagesHandler returns all Messages with paging
func GetMessagesHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, ErrorTickerNotFound))
		return
	}

	var messages []*Message
	//TODO: Pagination
	err = DB.Find("Ticker", tickerID, &messages, storm.Reverse())
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusOK, NewJSONSuccessResponse("messages", []string{}))
			return
		}

		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("messages", NewMessagesResponse(messages)))
}

//GetMessageHandler returns a Message for the given id
func GetMessageHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, ErrorTickerNotFound))
		return
	}

	var message Message
	messageID, err := strconv.Atoi(c.Param("messageID"))
	err = DB.One("ID", messageID, &message)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("message", NewMessageResponse(&message)))
}

//PostMessageHandler creates and returns a new Message
func PostMessageHandler(c *gin.Context) {
	var body struct {
		Text string `json:"text" binding:"required"`
	}
	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}

	message := NewMessage()
	message.Text = body.Text
	message.Ticker = tickerID

	if ticker.Twitter.Active {
		id, err := bridge.Twitter.Update(ticker, *message)
		if err == nil {
			message.TweetID = id
		} else {
			log.Error(err)
		}
	}

	err = DB.Save(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("message", NewMessageResponse(message)))
}

//DeleteTickerHandler deletes a existing Ticker
func DeleteMessageHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}

	messageID, err := strconv.Atoi(c.Param("messageID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	var message Message

	err = DB.One("ID", messageID, &message)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if message.TweetID != "" {
		err = bridge.Twitter.Delete(ticker, message.TweetID)
		if err != nil {
			log.Error(err)
		}
	}

	err = DB.DeleteStruct(&message)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   nil,
		"status": ResponseSuccess,
		"error":  nil,
	})
}

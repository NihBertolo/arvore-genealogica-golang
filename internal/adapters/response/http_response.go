package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONSuccessResult struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data"`
}

type JSONBadReqResult struct {
	Code    int         `json:"code" example:"400"`
	Message string      `json:"message" example:"Wrong Parameters"`
	Data    interface{} `json:"data"`
}

type JSONIntServerResult struct {
	Code    int         `json:"code" example:"500"`
	Message string      `json:"message" example:"Internal Error"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, JSONSuccessResult{
		Code:    http.StatusOK,
		Data:    data,
		Message: "Success",
	})
}

func FailResponse(c *gin.Context, respCode int, message string) {
	if respCode == http.StatusInternalServerError {
		c.JSON(respCode, JSONIntServerResult{
			Code:    respCode,
			Data:    nil,
			Message: message,
		})
		return
	}

	c.JSON(respCode, JSONBadReqResult{
		Code:    respCode,
		Data:    nil,
		Message: message,
	})
}

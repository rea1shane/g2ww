package ngalert

import (
	"g2ww/common"
	"g2ww/grafana/general"
	"github.com/gin-gonic/gin"
)

var counter = common.Counter{
	SentSuccessCount: 0,
	SentFailureCount: 0,
}

func GetSendCount(c *gin.Context) {
	general.GetSendCount(c, counter)
}

func SendMsg(c *gin.Context) {

}

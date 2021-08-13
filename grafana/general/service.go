package general

import (
	"fmt"
	"g2ww/common"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetSendCount 记录发送次数
func GetSendCount(c *gin.Context, counter common.Counter) {
	common.PrintCutOffRule()

	_, _ = c.Writer.WriteString("G2WW Server is running! \nParsed & forwarded " + strconv.Itoa(counter.SentSuccessCount) + " messages to WeChat Work! \nParsed & forwarded Failure " + strconv.Itoa(counter.SentFailureCount) + " messages!")

	fmt.Println("Sent Success Count:", counter.SentSuccessCount)
	fmt.Println("Sent Failure Count:", counter.SentFailureCount)
	fmt.Println()

	return
}

func DealError(c *gin.Context, err error, errMsg string) {
	fmt.Println(errMsg)
	fmt.Println(err.Error())
	_, _ = c.Writer.WriteString(errMsg)
}

func GrafanaWebhookUnmarshalJsonError(c *gin.Context, err error) int {
	errMsg := `[ERROR] JSON Unmarshal failure`
	DealError(c, err, errMsg)
	return common.GrafanaWebhookUnmarshalJsonError
}

func ClientCallAPIError(c *gin.Context, err error) int {
	errMsg := `[ERROR] Client call API failure`
	DealError(c, err, errMsg)
	return common.ClientCallAPIError
}

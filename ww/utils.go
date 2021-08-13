package ww

import (
	"bytes"
	"encoding/json"
	"fmt"
	"g2ww/common"
	"net/http"
)

// CheckWechatWorkResponse 检测企业微信是否接受成功
// 企业微信有 20条/min 的发送速率限制
func CheckWechatWorkResponse(r *http.Response) int {
	w := WechatWorkResponse{}
	buffer := new(bytes.Buffer)
	_, _ = buffer.ReadFrom(r.Body)
	_ = json.Unmarshal(buffer.Bytes(), &w)
	if w.ErrCode != 0 {
		fmt.Println(`[ERROR] Sending to WeChat-Work bot failure`)
		fmt.Printf("ErrorCode: [%v]", w.ErrCode)
		fmt.Println()
		fmt.Printf("ErrorMsg : [%v]", w.ErrMsg)
		fmt.Println()
		fmt.Println()
		return common.WechatWorkCallAPIError
	} else if w.ErrMsg == "ok. Warning: wrong json format." {
		fmt.Println(`[Warning] Wrong json format`)
		fmt.Println()
		return common.WechatWorkCallAPIWrongJsonFormatWarning
	} else if w.ErrMsg == "" {
		fmt.Println(`[Warning] Parse response body failure`)
		fmt.Println()
		return common.WechatWorkParseResponseBodyFailure
	} else {
		return common.OK
	}
}

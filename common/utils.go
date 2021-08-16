package common

import (
	"fmt"
)

func CheckStatus(s StatusCode, c *Counter) {
	if s == OK {
		c.SentSuccessCount++
	} else {
		c.SentFailureCount++
	}
	fmt.Printf("Status Code : %d", s)
	fmt.Println()
	fmt.Printf("Status Msg  : %s", s)
	fmt.Println()
}

func PrintCutOffRule() {
	fmt.Println()
	fmt.Println("----------------------------------------------------------------------")
	fmt.Println()
}

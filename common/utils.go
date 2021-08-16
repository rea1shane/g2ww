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
	fmt.Printf("status   : %d", s)
	fmt.Println()
	fmt.Printf("message  : %s", s)
	fmt.Println()
}

func PrintCutOffRule() {
	fmt.Println()
	fmt.Println("----------------------------------------------------------------------")
	fmt.Println()
}

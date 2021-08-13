package common

import "fmt"

func CheckStatus(status int, c *Counter) {
	if status == OK {
		c.SentSuccessCount++
	} else {
		c.SentFailureCount++
		if status == InternalError {
			fmt.Println(`[Error] Internal Error`)
		}
	}
}

func PrintCutOffRule() {
	fmt.Println()
	fmt.Println("----------------------------------------------------------------------")
	fmt.Println()
}

package txlog

import (
	"fmt"
	"strconv"
)

type ConsoleTxTransfer struct {
}

func NewConsoleTransfer() *ConsoleTxTransfer {
	return &ConsoleTxTransfer{}
}

func (c *ConsoleTxTransfer) Send(data *TxLogHistory) {
	fmt.Println("SEND :: History ->")

	fmt.Println("==============================================")
	for ii := 0; ii < len(data.list); ii++ {
		fmt.Println("INDEX : ", strconv.Itoa(ii+1), "DATA : ", data.list[ii].act)
	}
	fmt.Println("==============================================")
}

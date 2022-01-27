package txlog

type ITxTransfer interface {
	Send(data *TxLogHistory)
}

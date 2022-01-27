package txlog

import "testing"

func TestManager(t *testing.T) {
	trans := NewConsoleTransfer()
	txmgr := NewTxManager(trans)

	txmgr.Push("1", NewTxLog("ACT1"))
	txmgr.Push("1", NewTxLog("ACT2"))
	txmgr.Push("1", NewTxLog("ACT3"))
	txmgr.Finish("1", NewTxLog("ACT1"))
}

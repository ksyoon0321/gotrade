package txlog

import (
	"sync"
	"time"
)

type TxLog struct {
	ts  time.Time
	act string
}

func NewTxLog(act string) TxLog {
	return TxLog{
		ts:  time.Now(),
		act: act,
	}
}

type TxLogHistory struct {
	id   string
	list []TxLog
}

func NewTxLogHistory(id string) *TxLogHistory {
	return &TxLogHistory{
		id:   id,
		list: make([]TxLog, 0),
	}
}

type TxManager struct {
	mutex sync.RWMutex
	list  map[string]*TxLogHistory
	trans ITxTransfer
}

func NewTxManager(tran ITxTransfer) *TxManager {
	return &TxManager{
		trans: tran,
		list:  make(map[string]*TxLogHistory),
	}
}

func (t *TxManager) Push(id string, l TxLog) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, ok := t.list[id]; !ok {
		t.list[id] = NewTxLogHistory(t.makeTxId(id))
	}

	item := t.list[id]
	item.list = append(item.list, l)
}

func (t *TxManager) Finish(id string, l TxLog) {
	t.Push(id, l)
	//

	t.trans.Send(t.list[id])
	//
	delete(t.list, id)
}

func (t *TxManager) makeTxId(id string) string {
	return id + time.Now().Format("yyyyMMddHHmmss")
}

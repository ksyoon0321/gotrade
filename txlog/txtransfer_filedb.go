package txlog

import "github.com/ksyoon0321/gotrade/database"

type FileDBTxTransfer struct {
	filedb *database.FileDatabase
}

func NewFileDBTxTransfer(path string) *FileDBTxTransfer {
	return &FileDBTxTransfer{
		filedb: database.NewFileDatabase(path),
	}
}

func (db *FileDBTxTransfer) Send(data *TxLogHistory) {
	arrPayload := make([]*database.InsertPayLoad, len(data.list))
	for ii := 0; ii < len(data.list); ii++ {
		arrPayload[ii] = database.NewInsertPayLoad(data.id, db.filedb.ParseContent(data.list[ii].act))
	}

	db.filedb.InsertArray(arrPayload)
}

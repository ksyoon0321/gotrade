package database

import (
	"testing"
)

func TestFileDB(t *testing.T) {
	db := NewFileDatabase("/Users/yoonkwangsuk/Downloads/gofiledb")

	tid := "20220127hhnnss_KRW-CRO"
	tmap := make(map[string]interface{})
	tmap["cmd"] = "act_cmd"
	tmap["11"] = "22"
	tmap["33"] = "44"
	ret, err := db.Insert(NewInsertPayLoad(tid, tmap))

	if ret != 1 {
		t.Errorf("insert should be 1 not %d %s", ret, err)
	}

	arrW := make([]WhereCond, 0)
	w := NewWhereCond("11", "aa", COND_EQUAL)
	arrW = append(arrW, *w)

	tmap2 := make(map[string]interface{})
	tmap2["11"] = "aa"
	ret2, err2 := db.Update(NewUpdatePayLoad(tid, arrW, tmap2))

	if ret2 != 1 {
		t.Errorf("update should be 1 not %d %s", ret2, err2)
	}

	searchdata := db.Select(NewSelectPayLoad(tid, nil, nil))
	t.Errorf("searchdt = %s", searchdata)

	// ret3, err3 := db.Delete((*DeletePayLoad)(NewUpdatePayLoad(tid, arrW, nil)))
	// if ret3 != 1 {
	// 	t.Errorf("update should be 1 not %d %s", ret3, err3)
	// }

	tid4 := "20220127000000_KRW-CRO2"
	arrPayload := make([]*InsertPayLoad, 3)
	arrPayload[0] = NewInsertPayLoad(tid4, tmap)
	arrPayload[1] = NewInsertPayLoad(tid4, tmap)
	arrPayload[2] = NewInsertPayLoad(tid4, tmap)
	ret4, err4 := db.InsertArray(arrPayload)
	if ret4 != 3 {
		t.Errorf("it should 3 not %d %s", ret4, err4)
	}

}

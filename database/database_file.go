package database

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ksyoon0321/gotrade/tree"
	"github.com/ksyoon0321/gotrade/util"
)

type FormatFileSave struct {
	content string
}

type FileDatabase struct {
	rootdir string
	bcache  *tree.BTree
}

func NewFileDatabase(root string) *FileDatabase {
	if !util.IsExistsFileOrDir(root) {
		os.MkdirAll(root, os.ModePerm)
	}

	//temp file loc
	if !strings.HasSuffix(root, "/") {
		root = root + "/"
	}

	tmpdir := root + "tmp"
	if !util.IsExistsFileOrDir(tmpdir) {
		os.MkdirAll(tmpdir, os.ModePerm)
	}

	datadir := root + "data"
	if !util.IsExistsFileOrDir(datadir) {
		os.MkdirAll(datadir, os.ModePerm)
	}

	return &FileDatabase{
		rootdir: root,
		bcache:  tree.NewBTree(),
	}
}

//CRUD는 오류발생하지 않고 0 또는 가짜배열 반환
func (f *FileDatabase) Insert(payload *InsertPayLoad) (int, error) {
	if !f.checkDataDir(f.getFullDirByTid(payload.tid)) {
		os.MkdirAll(f.getFullDirByTid(payload.tid), os.ModePerm)
	}

	sdata := f.mapToContent(payload.data)
	if err := ioutil.WriteFile(f.getFullDirByTid(payload.tid)+"/"+payload.tid, sdata, 0644); err != nil {
		return 0, err
	}

	f.bcache.Put(payload.tid, sdata)
	return 1, nil
}

func (f *FileDatabase) InsertArray(payloads []*InsertPayLoad) (int, error) {
	if len(payloads) <= 0 {
		return 0, errors.New("ZERO LENGTH")
	}

	payload0 := payloads[0]
	if !f.checkDataDir(f.getFullDirByTid(payload0.tid)) {
		os.MkdirAll(f.getFullDirByTid(payload0.tid), os.ModePerm)
	}

	savedata := ""
	for _, item := range payloads {
		savedata += string(f.mapToContent(item.data)) + "\n"
	}

	if err := ioutil.WriteFile(f.getFullDirByTid(payload0.tid)+"/"+payload0.tid, []byte(savedata), 0644); err != nil {
		return 0, err
	}

	f.bcache.Put(payload0.tid, []byte(savedata))
	return len(payloads), nil
}

func (f *FileDatabase) Update(payload *UpdatePayLoad) (int, error) {
	if !f.checkDataDir(f.getFullDirByTid(payload.tid)) {
		os.MkdirAll(f.getFullDirByTid(payload.tid), os.ModePerm)
	}

	content := f.find(payload.tid)
	if content == "" {
		return 0, errors.New("NOT FOUND DATA FILE")
	}

	dbmap := f.ParseContent(content)

	for key, vl := range payload.data {
		dbmap[key] = vl
	}

	sdata := f.mapToContent(dbmap)
	if err := ioutil.WriteFile(f.getFullDirByTid(payload.tid)+"/"+payload.tid, sdata, 0644); err != nil {
		return 0, err
	}

	f.bcache.Put(payload.tid, sdata)
	return 1, nil
}

func (f *FileDatabase) Delete(payload *DeletePayLoad) (int, error) {
	if !f.checkDataDir(f.getFullDirByTid(payload.tid)) {
		os.MkdirAll(f.getFullDirByTid(payload.tid), os.ModePerm)
	}

	//delete bcache

	fullPath := f.getFullDirByTid(payload.tid) + "/" + payload.tid
	err := os.Remove(fullPath)
	if err != nil {
		return 0, err
	}

	//bcache에서 nil반환되면 없다고 간주되므로..
	f.bcache.Put(payload.tid, nil)
	return 1, nil
}

func (f *FileDatabase) Select(payload *SelectPayLoad) []interface{} {
	if !f.checkDataDir(f.getFullDirByTid(payload.tid)) {
		os.MkdirAll(f.getFullDirByTid(payload.tid), os.ModePerm)
	}

	arr := make([]interface{}, 0)
	if payload.tid == "" {

	} else {
		//tid검색
		arr = append(arr, f.ParseContent(f.find(payload.tid)))

	}

	return arr
}

func (f *FileDatabase) checkDataDir(path string) bool {
	return util.IsExistsFileOrDir(path)
}

func (f *FileDatabase) getFullDirByTid(tid string) string {
	//tid = yyyymmddhhnnss
	dir := time.Now().Format("yyyyMMdd")
	if len(tid) >= 8 {
		dir = tid[:7]
	}

	return f.rootdir + "data/" + dir
}
func (f *FileDatabase) find(tid string) string {
	//bcache first
	//bcache.find
	cachedata := f.bcache.Get(tid)
	if cachedata != nil {
		return string(cachedata.([]byte))
	} else {
		fullPath := f.getFullDirByTid(tid) + "/" + tid
		if util.IsExistsFileOrDir(fullPath) {
			content, err := ioutil.ReadFile(fullPath)
			if err != nil {
				return ""
			}

			return string(content)
		}
	}
	return ""
}

func (f *FileDatabase) ParseContent(cont string) map[string]interface{} {
	contmap := make(map[string]interface{})

	arr := strings.Split(cont, ":")
	for _, str := range arr {
		if strings.Contains(str, "=") {
			keyval := strings.Split(str, "=")
			contmap[keyval[0]] = keyval[1]
		} else {
			contmap["CMD"] = str
		}
	}

	return contmap
}

func (f *FileDatabase) mapToContent(contmap map[string]interface{}) []byte {
	var cont string

	for key, vl := range contmap {
		if strings.ToUpper(key) == "CMD" {
			cont = vl.(string) + cont
		} else {
			cont += ":" + strings.ToUpper(key) + "=" + vl.(string)
		}
	}
	return []byte(cont)
}

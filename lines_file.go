package log33

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type File33 struct {
	dirName  string
	fileName string
	fos      *os.File
	suffix   string
	filePath string
	maxLines int
	mu       sync.RWMutex
}
type k1 struct {
	Key string `json:"key"`
}

func NewFile(dirName, fileName, suffix string, maxLines int) (*File33, error) {
	if fileName == "" {
		return nil, errors.New("file name not empty")
	}
	filePath := fmt.Sprintf("%s/%s", dirName, fileName)
	//if fileName != "" && suffix == "" {
	//	suffix = "txt"
	//}
	if suffix != ""{
		filePath = filePath + "."
	}
	filePath = filePath + suffix

	mf := &File33{
		fileName: fileName,
		dirName:  dirName,
		suffix:   suffix,
		maxLines: maxLines,
		filePath: filePath,
	}
	err := mf.initFileDir()
	if err != nil {
		return nil, err
	}
	return mf, err
}
func (f33 *File33) initFileDir() error {
	fn, err := os.Stat(f33.fileName)
	if os.IsNotExist(err) {
		err = os.MkdirAll(f33.dirName, 0766)
		if err != nil {
			return err
		}
	}

	if fn == nil && f33.fileName != "" {
		err = f33.createFile()
		if err != nil {
			return err
		}
	}
	return nil
}
func (f33 *File33) resetFos() {
	_ = f33.fos.Close()
}
func (f33 *File33) createFile() error {
	fos, err := os.OpenFile(f33.filePath, os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	f33.fos = fos
	return nil
}
func (f33 *File33) getAppend() error {
	f33.resetFos()
	fos, err := os.OpenFile(f33.filePath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	f33.fos = fos
	return nil
}
func (f33 *File33) getReadWrite() error {
	f33.resetFos()
	fos, err := os.OpenFile(f33.filePath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	f33.fos = fos
	return nil
}
func (f33 *File33) getRead() error {
	f33.resetFos()
	fos, err := os.OpenFile(f33.filePath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	f33.fos = fos
	return nil
}

func (f33 *File33) WriteLinesAppends(data []string) error {
	f33.mu.Lock()
	defer f33.mu.Unlock()
	oldList, err := f33.ReadLines()
	if err != nil {
		return err
	}
	var cn int
	var sumCn = len(oldList) + len(data)
	if f33.maxLines > 0 && sumCn > f33.maxLines {
		cn = sumCn - f33.maxLines
	}
	oldList = append(oldList, append(data)...)
	var insertList = oldList[cn:]
	err = f33.getReadWrite()
	if err != nil {
		return err
	}
	err = f33.fos.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f33.fos.Seek(0, 0)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f33.fos)
	for _, datum := range insertList {
		_, _ = fmt.Fprintln(w, datum)
	}
	_ = w.Flush()
	return nil
}
func (f33 *File33) ReadLines() ([]string, error) {
	var data = make([]string, 0)
	err := f33.getRead()
	if err != nil {
		return data, err
	}
	rd := bufio.NewReader(f33.fos)
	for {
		lineStr, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			if lineStr == "" {
				break
			}
		}
		data = append(data, strings.Replace(lineStr, "\n", "", -1))
	}
	return data, nil
}
func (f33 *File33) ReadOne(k string) (string, error) {
	var data string
	err := f33.getRead()
	if err != nil {
		return data, err
	}
	rd := bufio.NewReader(f33.fos)
	for {
		lineStr, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			if lineStr == "" {
				break
			}
		}
		var k2 = &k1{}
		lineStr = strings.Replace(lineStr, "\n", "", -1)
		err = json.Unmarshal([]byte(lineStr), &k2)
		if err != nil {
			return "", err
		}
		if k2.Key == k {
			data = strings.Replace(lineStr, "\n", "", -1)
			break
		}
	}
	return data, nil
}
func (f33 *File33) UpdateOne(k, v string) error {
	f33.mu.Lock()
	defer f33.mu.Unlock()
	oldList, err := f33.ReadLines()
	if err != nil {
		return err
	}
	var newList = make([]string, 0)

	for _, s := range oldList {
		var k2 = &k1{}
		s = strings.Replace(s, "\n", "", -1)
		err = json.Unmarshal([]byte(s), &k2)
		if err != nil {
			return err
		}
		if k2.Key == k {
			newList = append(newList, v)
		} else {
			newList = append(newList, s)
		}
	}
	err = f33.getReadWrite()
	_ = f33.fos.Truncate(0)
	_, _ = f33.fos.Seek(0, 0)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f33.fos)
	for _, datum := range newList {
		_, _ = fmt.Fprintln(w, datum)
	}
	_ = w.Flush()
	return nil
}
func (f33 *File33) DelOne(k string) error {
	f33.mu.Lock()
	defer f33.mu.Unlock()
	oldList, err := f33.ReadLines()
	if err != nil {
		return err
	}

	var befor = make([]string, 0)
	var after = make([]string, 0)

	var t1 bool
	for _, s := range oldList {
		var k2 = &k1{}
		s = strings.Replace(s, "\n", "", -1)
		err = json.Unmarshal([]byte(s), &k2)
		if err != nil {
			return err
		}
		if k2.Key == k {
			t1 = true
			continue
		}
		if t1 == false {
			befor = append(befor, s)
		}
		if t1 == true {
			after = append(after, s)
		}
	}
	insertList := append(befor, append(after)...)

	err = f33.getReadWrite()
	_ = f33.fos.Truncate(0)
	_, _ = f33.fos.Seek(0, 0)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f33.fos)
	for _, datum := range insertList {
		_, _ = fmt.Fprintln(w, datum)
	}
	_ = w.Flush()
	return nil
}
func (f33 *File33) DelFile()error {
	f33.mu.Lock()
	defer f33.mu.Unlock()
	f33.resetFos()
	return os.Remove(f33.filePath)
}

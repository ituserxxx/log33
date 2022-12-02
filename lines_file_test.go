package log33

import (
	"encoding/json"
	"fmt"
	"testing"
)

var f3 *File33

func init() {
	f, err := NewFile("./file_tmp", "a", "log", 0)
	if err != nil {
		panic(err.Error())
	}
	f3 = f
}
func TestFile33_WriteLinesAppends(t *testing.T) {
	atr := []map[string]string{
		{"key": "1111", "c": "1", "d": " 2"},
		{"key": "222", "c": "1", "d": "2"},
		{"key": "333", "e": "edg", "f": "fpx"},
	}
	var ls []string
	for _, m := range atr {
		v,_:= json.Marshal(m)
		ls = append(ls ,string(v))
	}
	err := f3.WriteLinesAppends(ls)
	if err != nil {
		panic(err.Error())
	}
}
func TestFile33_ReadLines(t *testing.T) {
	l, err := f3.ReadLines()
	if err != nil {
		panic(err.Error())
	}
	for _, s := range l {
		println(s)
	}
}
func TestFile33_ReadOne(t *testing.T) {
	l, err := f3.ReadOne("1111")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("%#v", l)
	}
}
func TestFile33_UpdateOne(t *testing.T) {
	up,_ := json.Marshal(map[string]string{"key": "1111", "c": "999", "d": " 999"})
	err := f3.UpdateOne("1111",string(up))
	if err != nil {
		panic(err.Error())
	}
	TestFile33_ReadLines(t)
}
func TestFile33_DelOne(t *testing.T) {
	err := f3.DelOne("1111")
	if err != nil {
		panic(err.Error())
	}
	TestFile33_ReadLines(t)
}
func TestFile33_DelFile(t *testing.T) {
	err := f3.DelFile()
	if err != nil {
		println(err.Error())
	}
}
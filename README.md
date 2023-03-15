This is a DIY log library using to golang support custom directory 、file name 、suffix create log and support limit
maximum number of lines

mode of use the following:

```go
// new file
//f, err := NewFile("/data/file_tmp", "a", "log", 0)
func Da(){
    f, err := NewFile("./file_tmp", "a", "log", 0)
    if err != nil {
        panic(err.Error())
    }
    // [ notice ]:  line content  must  have a key
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
    err := f.WriteLinesAppends()
	if err != nil {
        panic(err.Error())
    }
    // check out test for more information
}
```
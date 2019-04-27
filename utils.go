package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
)

func ReadFile(strFileName string) []byte {
	file, err := os.Open(strFileName)
	defer file.Close()
	if err != nil {
		log.Println("open file failed! ", err)
		return make([]byte, 0)
	}
	fi, _ := file.Stat()
	buffer := make([]byte, fi.Size())
	if _, err = file.Read(buffer); err != nil {
		log.Println("read file failed! ", err)
		return make([]byte, 0)
	}
	return buffer
}

func IntCondAssert(ListLen, cmpLen int, char string) bool {
	switch char {
	case "<":
		return ListLen < cmpLen
	case ">":
		return ListLen > cmpLen
	case "==":
		return ListLen == cmpLen
	case ">=":
		return ListLen >= cmpLen
	case "<=":
		return ListLen <= cmpLen
	}
	return false
}

func StringCondAssert(src, sub, char string) bool {
	if sub == "null" {
		sub = ""
	}
	switch char {
	case "==":
		return src == sub
	case "!=":
		return src != sub
	case "in":
		if idx := strings.Index(src, sub); idx != -1 {
			return true
		}
	case "not-in":
		if idx := strings.Index(src, sub); idx == -1 {
			return true
		}
	}
	return false
}

func GetFileList(path string) []string {
	var ret []string
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		ret = append(ret, path)
		return nil
	})
	if err != nil {
		log.Printf("walk err:%v\n", err)
	}
	return ret
}

func GetFileSize(infile string) int64 {
	file, err := os.Open(infile)
	if err != nil && os.IsNotExist(err) {
		return 0
	}
	defer file.Close()
	fi, err := file.Stat()
	return int64(fi.Size())
}

func DelLikeFile(path, like string) []string {
	var ret []string
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			log.Println(path)
			return nil
		}
		return nil
	})
	if err != nil {
		log.Printf("walk err:%v\n", err)
	}
	return ret
}

func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		log.Println(err.Error())
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}

func XMLCharSwitch(buf string) string {
	buf = strings.Replace(buf, "&", "&amp;", -1)
	buf = strings.Replace(buf, "<", "&lt;", -1)
	buf = strings.Replace(buf, ">", "&gt;", -1)
	return buf
}

func IndexN(s, key string, Nth int) int {
	if Nth <= 0 {
		return -1
	}
	idx, str := 0, s
	for i := 0; i < Nth; i++ {
		if index := strings.Index(str, key); index == -1 {
			return -1
		} else {
			idx += index
			if (index + 1) >= len(str) {
				return -1
			}
			str = str[index+1:]
			if i+1 < Nth {
				idx += 1
			}
		}
	}
	return idx
}

func IndexNth(s, key string, Nth int) int {
	re, _ := regexp.Compile(key)
	if IndexS := re.FindAllStringIndex(s, Nth); len(IndexS) > 0 {
		return IndexS[len(IndexS)-1][0]
	}
	return -1
}

func IndexNthL(s, key string, Nth int) (Len, idx int) {
	re, _ := regexp.Compile(key)
	if IndexS := re.FindAllStringIndex(s, Nth); len(IndexS) > 0 {
		return len(IndexS), IndexS[len(IndexS)-1][0]
	}
	return -1, -1
}

func LastIndexN(s, key string, Nth int) int {
	if Nth <= 0 {
		return -1
	}
	idx, str := 0, s
	for i := 0; i < Nth; i++ {
		if index := strings.LastIndex(str, key); index == -1 {
			return -1
		} else {
			idx = index
			if index-1 < 0 {
				return -1
			}
			str = str[:index-1]
		}
	}
	return idx
}

func SlashLinux(r rune) rune {
	if r == '\\' {
		return '/'
	}
	return r
}

func SlashWindows(r rune) rune {
	if r == '/' {
		return '\\'
	}
	return r
}

func WinPath(buf string) string {
	if idx := strings.Index(buf, "\\"); idx != -1 {
		return strings.Replace(buf, "\\\\", "\\", -1)
	}
	return buf
}

func SwitchLinuxPath(file string) string {
	return strings.Map(SlashLinux, file)
}

func SwitchWindowsPath(file string) string {
	return strings.Map(SlashWindows, file)
}

func TrimEnterChar(src string) string {
	re, _ := regexp.Compile("\\s{2,}")
	return re.ReplaceAllString(src, "\n")
}

func ClearMapSS(m map[string]string) {
	for k, _ := range m {
		delete(m, k)
	}
}

func ClearMapII(m map[int]int) {
	for k, _ := range m {
		delete(m, k)
	}
}

func ClearMapIS(m map[int]string) {
	for k, _ := range m {
		delete(m, k)
	}
}

func ClearMapSI(m map[string]int) {
	for k, _ := range m {
		delete(m, k)
	}
}

type signaler func(s os.Signal, arg interface{})
type signalSet struct {
	M map[os.Signal]signaler
}

func newSignalSet() *signalSet {
	ss := new(signalSet)
	ss.M = make(map[os.Signal]signaler)
	return ss
}

func (set *signalSet) register(s os.Signal, handler signaler) {
	if _, found := set.M[s]; !found {
		set.M[s] = handler
	}
}

func (set *signalSet) handle(sig os.Signal, arg interface{}) (err error) {
	if _, found := set.M[sig]; found {
		set.M[sig](sig, arg)
		return nil
	} else {
		return fmt.Errorf("No handler available for signal %v", sig)
	}
	panic("won't reach here")
}

func RegisterSignal(s []os.Signal, msg chan string) {
	sigS := newSignalSet()
	handler := func(s os.Signal, arg interface{}) {
		log.Printf("handle signal: %v\n", s)
		msg <- "done"
	}
	for _, v := range s {
		sigS.register(v, handler)
	}
	for {
		c := make(chan os.Signal)
		var sigs []os.Signal
		for sig := range sigS.M {
			sigs = append(sigs, sig)
		}
		signal.Notify(c)
		sig := <-c
		sigS.handle(sig, nil)
	}
}

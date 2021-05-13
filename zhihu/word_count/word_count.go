package word_count

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

type WordCount map[string]interface{}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	return p[j].Value < p[i].Value
}

// 《》<> 分割句子
func SplitByMoreStr(r rune) bool {
	splitSymbol := []rune("《》<>")
	for _, v := range splitSymbol {
		if r == v {
			return true
		}
	}
	return false
}

func (wc WordCount) SplitAndStatistics(s string) {
	dist1 := strings.FieldsFunc(s, SplitByMoreStr)
	for _, v := range dist1 {
		flag := 0
		v = strings.Replace(v, " ", "", -1)
		for key := range wc {
			if strings.Index(v, key) != -1 {
				wc[key] = wc[key].(int) + 1
				flag = 1
			}
		}
		if flag == 0 {
			if wc[v] == nil {
				wc[v] = 1
			} else {
				wc[v] = wc[v].(int) + 1
			}
		}
	}
}

func (wc WordCount) ReadFile(f *os.File) {
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		wc.SplitAndStatistics(line)
	}
}

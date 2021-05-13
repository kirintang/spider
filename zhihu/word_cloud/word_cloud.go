package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-echarts/go-echarts/charts"
	wordCount "github.com/kirintang/spider/zhihu/word_count"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	questionId := os.Args[1]
	nwc := charts.NewWordCloud()

	nwc.SetGlobalOptions(charts.TitleOpts{Title: "知乎问题:"})
	wc := make(wordCount.WordCount)
	f, err := os.Open("./answer/" + questionId + ".txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	wc.ReadFile(f)
	nwc.Add("wordcloud", wc, charts.WordCloudOpts{SizeRange: []float32{14, 250}})
	nwc.Render(w)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Missing question_id arguments")
		os.Exit(1)
	}
	questionId := os.Args[1]
	if !Exists("./answer/" + questionId + ".txt") {
		log.Fatalf("file answer/%s.text not exists!\n", questionId)
		os.Exit(1)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}

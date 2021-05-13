package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Missing question_id argument")
		os.Exit(1)
	}
	questionID := os.Args[1]
	workingDir, _ := os.Getwd()
	file, err := os.OpenFile(workingDir+"/download/answer/"+questionID+".txt", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	total := 20
	i := 0
	c := colly.NewCollector(func(collector *colly.Collector) {
		extensions.RandomUserAgent(collector)
	})
	c.OnRequest(func(request *colly.Request) {
		fmt.Printf("FETCH -> %s\n", request.URL.String())
	})
	c.OnResponse(func(resp *colly.Response) {
		var f interface{}
		json.Unmarshal(resp.Body, &f)
		paging := f.(map[string]interface{})["paging"]
		total = int(paging.(map[string]interface{})["totals"].(float64))
		data := f.(map[string]interface{})["data"]
		for k, v := range data.([]interface{}) {
			content := v.(map[string]interface{})["content"]
			reader := strings.NewReader(content.(string))
			doc, _ := goquery.NewDocumentFromReader(reader)
			file.Write([]byte(fmt.Sprintf("%d:%s\n", i+k, doc.Find("p").Text())))
		}
	})
	for ; i <= total; i += 20 {
		url := fmt.Sprintf("https://www.zhihu.com/api/v4/questions/%s/answers?include=data[*].is_normal,admin_closed_comment,reward_info,is_collapsed,annotation_action,annotation_detail,collapse_reason,is_sticky,collapsed_by,suggest_edit,comment_count,can_comment,content,editable_content,voteup_count,reshipment_settings,comment_permission,created_time,updated_time,review_info,relevant_info,question,excerpt,relationship.is_authorized,is_author,voting,is_thanked,is_nothelp,is_labeled,is_recognized,paid_info,paid_info_content;data[*].mark_infos[*].url;data[*].author.follower_count,badge[*].topics&offset=%d&limit=%d&sort_by=updated", questionID, i, 20)
		c.Visit(url)
	}
}

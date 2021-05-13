package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/SpaceX-io/go-lib/rand"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func main() {
	t := time.Now()
	if len(os.Args) != 2 {
		log.Println("Missing query argument")
		os.Exit(1)
	}
	query := os.Args[1]
	workingDir, _ := os.Getwd()
	outPutDir := fmt.Sprintf(workingDir+"/download/imgs/%s", query)
	c := colly.NewCollector(func(collector *colly.Collector) {
		extensions.RandomUserAgent(collector)
		collector.Async = true
	})
	imageC := c.Clone()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "BIGipServerPools_Web_ssl=2806819008.47873.0000; Hm_lvt_c01558ab05fd344e898880e9fc1b65c4=1620712471; accessId=578c8dc0-6fab-11e8-ab7a-fda8d0606763; qimo_seosource_578c8dc0-6fab-11e8-ab7a-fda8d0606763=%E7%BB%94%E6%AC%8F%E5%94%B4; qimo_seokeywords_578c8dc0-6fab-11e8-ab7a-fda8d0606763=; pageViewNum=3; Hm_lpvt_c01558ab05fd344e898880e9fc1b65c4=1620716089")
		r.Headers.Add("Connection", "keep-alive")
		r.Headers.Add("Referer", "https://www.quanjing.com/search.aspx?q="+url.QueryEscape(query))
		r.Headers.Add("sec-ch-ua", "'Google Chrome';v='89', 'Chromium';v='89', ';Not A Brand';v='99'")
		r.Headers.Add("sec-ch-ua-mobile", "?0")
		r.Headers.Add("Sec-Fetch-Mode", "cors")
		r.Headers.Add("Sec-Fetch-Site", "same-origin")
		r.Headers.Add("Sec-Fetch-Dest", "empty")
		r.Headers.Add("Host", "www.quanjing.com")
		r.Headers.Add("Accept", "text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01")
		r.Headers.Add("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
		r.Headers.Add("X-Requested-With", "XMLHttpRequest")
	})

	c.OnResponse(func(r *colly.Response) {
		var f interface{}
		if err := json.Unmarshal(r.Body[13:len(r.Body)-1], &f); err != nil {
			panic(err)
		}
		imgList := f.(map[string]interface{})["imglist"]
		for k, img := range imgList.([]interface{}) {
			url := img.(map[string]interface{})["imgurl"].(string)
			url = url + "#" + img.(map[string]interface{})["caption"].(string)
			fmt.Printf("find -> %d:%s\n", k, url)
			imageC.Visit(url)
		}
	})
	c.OnError(func(response *colly.Response, err error) {
		fmt.Println(err)
	})
	imageC.OnResponse(func(r *colly.Response) {
		os.MkdirAll(outPutDir, os.ModePerm)
		fileName := rand.String(20) + ".jpg"
		fmt.Printf("download -> %s \n", fileName)
		r.Save(outPutDir + "/" + fileName)
	})
	pageSize := 200
	pageNum := 10
	for i := 0; i < pageNum; i++ {
		url := fmt.Sprintf("https://www.quanjing.com/Handler/SearchUrl.ashx?t=6554&callback=searchresult&q="+query+"&stype=1&pagesize=%d&pagenum=%d&imageType=2&imageColor=&brand=&imageSType=&fr=1&sortFlag=1&imageUType=&btype=&authid=&_=1620716088380", pageSize, i)
		_ = c.Visit(url)
	}
	c.Wait()
	imageC.Wait()
	fmt.Printf("done! cost: %s\n", time.Since(t))
}

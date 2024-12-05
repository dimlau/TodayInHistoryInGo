package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type History struct {
	Type    string `json:"type"`
	Year    string `json:"year"`
	Contant string `json:"contant"`
}

func generateDates() []string {
	months := []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"}
	var dates []string
	for _, month := range months {
		daysInMonth := 31
		if month == "2月" {
			daysInMonth = 29
		} else if month == "4月" || month == "6月" || month == "9月" || month == "11月" {
			daysInMonth = 30
		}
		for day := 1; day <= daysInMonth; day++ {
			dates = append(dates, fmt.Sprintf("%s%d日", month, day))
		}
	}
	return dates
}

func saveToFile(data []History, date string) error {
	dir := "./data"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("创建目录失败: %v", err)
		}
	}

	filePath := fmt.Sprintf("%s/%s.json", dir, date)
	filePtr, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer filePtr.Close()

	encoder := json.NewEncoder(filePtr)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("编码 JSON 失败: %v", err)
	}
	return nil
}

func main() {
	datem := generateDates()
	//datem := []string{"11月5日"}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // 限制并发数为 5

	for _, d := range datem {
		wg.Add(1)
		go func(date string) {
			defer wg.Done()
			sem <- struct{}{} // 占用一个并发槽

			log.Printf("开始抓取: %s\n", date)
			J := setdaysfromwiki(date)
			if err := saveToFile(J, date); err != nil {
				log.Printf("保存失败: %s, 错误: %v\n", date, err)
			}

			<-sem // 释放一个并发槽
		}(d)
	}

	wg.Wait()
}

func setdaysfromwiki(date string) (JF []History) {
	URL := "https://zh.wikipedia.org/zh-cn/"
	C := colly.NewCollector(
		colly.CacheDir("./cache"),
	)
	C.WithTransport(&http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	var all int
	C.OnHTML("#toc-大事記-sublist", func(k *colly.HTMLElement) {
		k.ForEach("li", func(i int, h *colly.HTMLElement) {
			all += 1
		})
	})
	C.OnHTML("#toc-大事紀-sublist", func(k *colly.HTMLElement) {
		k.ForEach("li", func(i int, h *colly.HTMLElement) {
			all += 1
		})
	})
	C.OnHTML("#toc-大事记-sublist", func(k *colly.HTMLElement) {
		k.ForEach("li", func(i int, h *colly.HTMLElement) {
			all += 1
		})
	})
	C.OnHTML("#mw-content-text", func(e *colly.HTMLElement) {
		e.ForEach("sup", func(_ int, e1 *colly.HTMLElement) {
			e1.DOM.Remove()
		})
		e.ForEach("style", func(_ int, e1 *colly.HTMLElement) {
			e1.DOM.Remove()
		})
		e.ForEach("div.mw-parser-output > ul", func(i int, h *colly.HTMLElement) {
			daytype := "0"
			switch {
			case i < all:
				daytype = "1"
			case i == all:
				daytype = "2"
			case i == all+1:
				daytype = "3"
			}
			h.DOM.Children().Each(func(i int, h1 *goquery.Selection) {
				MYTXT := h1.Text()
				re, _ := regexp.Compile(`在比利时举行的一级方程式大奖赛中`)
				MYTXT = re.ReplaceAllStringFunc(MYTXT, func(s string) string {
					return "1982年：在比利时举行的一级方程式大奖赛中"
				})
				reff, _ := regexp.Compile(`1937年12月16-17日，`)
				MYTXT = reff.ReplaceAllStringFunc(MYTXT, func(s string) string {
					return "1937年：12月16-17日，"
				})
				re7, _ := regexp.Compile(`黄丽芳，香港女配音员`)
				MYTXT = re7.ReplaceAllStringFunc(MYTXT, func(s string) string {
					return "1979年：黄丽芳，香港女配音员"
				})
				re8, _ := regexp.Compile(`年菲利普·勒`)
				MYTXT = re8.ReplaceAllStringFunc(MYTXT, func(s string) string {
					return "年：菲利普·勒"
				})
				re1, _ := regexp.Compile(`[名|年]\s?[：|︰|——|﹕｜；|:| :]\s?`)
				MYTXT = re1.ReplaceAllStringFunc(MYTXT, func(_ string) string {
					return "年："
				})
				re2, _ := regexp.Compile(`\d+年，`)
				MYTXT = re2.ReplaceAllStringFunc(MYTXT, func(s string) string {
					return strings.TrimRight(s, "，") + "："
				})
				re4, _ := regexp.Compile(`\d+年\s?（[\S\s]+?）：`)
				MYTXT = re4.ReplaceAllStringFunc(MYTXT, func(s string) string {
					return s + "年："
				})
				re5, _ := regexp.Compile(`\d+(]|：)`)
				MYTXT = re5.ReplaceAllStringFunc(MYTXT, func(s string) string {
					tt := strings.TrimRight(s, "]")
					tt = strings.TrimRight(tt, "：")
					return tt + "年："
				})
				re6, _ := regexp.Compile(`生年不\S?：`)
				MYTXT = re6.ReplaceAllStringFunc(MYTXT, func(s string) string {
					return s + "年："
				})
				re9, _ := regexp.Compile(`2003年伯纳德·卡茨，`)
				MYTXT = re9.ReplaceAllStringFunc(MYTXT, func(s string) string {
					return "2003年：伯纳德·卡茨，"
				})
				reX, _ := regexp.Compile(`1990年随黄日华签约亚视`)
				MYTXT = reX.ReplaceAllStringFunc(MYTXT, func(s string) string {
					return "1990年：随黄日华签约亚视"
				})
				reX1, _ := regexp.Compile(`2023年]：赵有亮，`)
				MYTXT = reX1.ReplaceAllStringFunc(MYTXT, func(s string) string {
					return "2023年：赵有亮，"
				})
				reX2, _ := regexp.Compile(`1994年申根公约生效。`)
				MYTXT = reX2.ReplaceAllStringFunc(MYTXT, func(s string) string {
					return "1994年：申根公约生效。"
				})
				reX3, _ := regexp.Compile(`907年，完全控制朝廷`)
				MYTXT = reX3.ReplaceAllStringFunc(MYTXT, func(s string) string {
					return "907年：完全控制朝廷"
				})
				var newTD History
				if h1.Children().Is("ul") {
					h1.Children().Children().Each(func(_ int, li *goquery.Selection) {
						if daytype == "1" || daytype == "2" || daytype == "3" {
							log.Println("条目格式化", MYTXT)
							newTD = History{
								Type:    daytype,
								Year:    strings.Split(MYTXT, "：")[0],
								Contant: li.Text(),
							}
							JF = append(JF, newTD)
						}
					})
				} else if h1.Children().Nodes != nil {
					if daytype == "1" || daytype == "2" || daytype == "3" {
						log.Println("条目格式化", MYTXT)
						if len(strings.Split(MYTXT, "年：")) <= 1 {
							log.Println("条目格式化失败", MYTXT)
						}
						newTD = History{
							Type:    daytype,
							Year:    strings.Split(MYTXT, "：")[0],
							Contant: strings.Split(MYTXT, "年：")[1],
						}
						JF = append(JF, newTD)
					}
				}
			})
		})
	})
	C.Visit(URL + date)
	C.Wait()
	return
}

package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type History struct {
	Type    string `json:"type"`
	Year    string `json:"year"`
	Contant string `json:"contant"`
}

func main() {
	datem := []string{"1月1日", "1月2日", "1月3日", "1月4日", "1月5日", "1月6日", "1月7日", "1月8日", "1月9日", "1月10日", "1月11日", "1月12日", "1月13日", "1月14日", "1月15日", "1月16日", "1月17日", "1月18日", "1月19日", "1月20日", "1月21日", "1月22日", "1月23日", "1月24日", "1月25日", "1月26日", "1月27日", "1月28日", "1月29日", "1月30日", "1月31日", "2月1日", "2月2日", "2月3日", "2月4日", "2月5日", "2月6日", "2月7日", "2月8日", "2月9日", "2月10日", "2月11日", "2月12日", "2月13日", "2月14日", "2月15日", "2月16日", "2月17日", "2月18日", "2月19日", "2月20日", "2月21日", "2月22日", "2月23日", "2月24日", "2月25日", "2月26日", "2月27日", "2月28日", "2月29日", "3月1日", "3月2日", "3月3日", "3月4日", "3月5日", "3月6日", "3月7日", "3月8日", "3月9日", "3月10日", "3月11日", "3月12日", "3月13日", "3月14日", "3月15日", "3月16日", "3月17日", "3月18日", "3月19日", "3月20日", "3月21日", "3月22日", "3月23日", "3月24日", "3月25日", "3月26日", "3月27日", "3月28日", "3月29日", "3月30日", "3月31日", "4月1日", "4月2日", "4月3日", "4月4日", "4月5日", "4月6日", "4月7日", "4月8日", "4月9日", "4月10日", "4月11日", "4月12日", "4月13日", "4月14日", "4月15日", "4月16日", "4月17日", "4月18日", "4月19日", "4月20日", "4月21日", "4月22日", "4月23日", "4月24日", "4月25日", "4月26日", "4月27日", "4月28日", "4月29日", "4月30日", "5月1日", "5月2日", "5月3日", "5月4日", "5月5日", "5月6日", "5月7日", "5月8日", "5月9日", "5月10日", "5月11日", "5月12日", "5月13日", "5月14日", "5月15日", "5月16日", "5月17日", "5月18日", "5月19日", "5月20日", "5月21日", "5月22日", "5月23日", "5月24日", "5月25日", "5月26日", "5月27日", "5月28日", "5月29日", "5月30日", "5月31日", "6月1日", "6月2日", "6月3日", "6月4日", "6月5日", "6月6日", "6月7日", "6月8日", "6月9日", "6月10日", "6月11日", "6月12日", "6月13日", "6月14日", "6月15日", "6月16日", "6月17日", "6月18日", "6月19日", "6月20日", "6月21日", "6月22日", "6月23日", "6月24日", "6月25日", "6月26日", "6月27日", "6月28日", "6月29日", "6月30日", "7月1日", "7月2日", "7月3日", "7月4日", "7月5日", "7月6日", "7月7日", "7月8日", "7月9日", "7月10日", "7月11日", "7月12日", "7月13日", "7月14日", "7月15日", "7月16日", "7月17日", "7月18日", "7月19日", "7月20日", "7月21日", "7月22日", "7月23日", "7月24日", "7月25日", "7月26日", "7月27日", "7月28日", "7月29日", "7月30日", "7月31日", "8月1日", "8月2日", "8月3日", "8月4日", "8月5日", "8月6日", "8月7日", "8月8日", "8月9日", "8月10日", "8月11日", "8月12日", "8月13日", "8月14日", "8月15日", "8月16日", "8月17日", "8月18日", "8月19日", "8月20日", "8月21日", "8月22日", "8月23日", "8月24日", "8月25日", "8月26日", "8月27日", "8月28日", "8月29日", "8月30日", "8月31日", "9月1日", "9月2日", "9月3日", "9月4日", "9月5日", "9月6日", "9月7日", "9月8日", "9月9日", "9月10日", "9月11日", "9月12日", "9月13日", "9月14日", "9月15日", "9月16日", "9月17日", "9月18日", "9月19日", "9月20日", "9月21日", "9月22日", "9月23日", "9月24日", "9月25日", "9月26日", "9月27日", "9月28日", "9月29日", "9月30日", "10月1日", "10月2日", "10月3日", "10月4日", "10月5日", "10月6日", "10月7日", "10月8日", "10月9日", "10月10日", "10月11日", "10月12日", "10月13日", "10月14日", "10月15日", "10月16日", "10月17日", "10月18日", "10月19日", "10月20日", "10月21日", "10月22日", "10月23日", "10月24日", "10月25日", "10月26日", "10月27日", "10月28日", "10月29日", "10月30日", "10月31日", "11月1日", "11月2日", "11月3日", "11月4日", "11月5日", "11月6日", "11月7日", "11月8日", "11月9日", "11月10日", "11月11日", "11月12日", "11月13日", "11月14日", "11月15日", "11月16日", "11月17日", "11月18日", "11月19日", "11月20日", "11月21日", "11月22日", "11月23日", "11月24日", "11月25日", "11月26日", "11月27日", "11月28日", "11月29日", "11月30日", "12月1日", "12月2日", "12月3日", "12月4日", "12月5日", "12月6日", "12月7日", "12月8日", "12月9日", "12月10日", "12月11日", "12月12日", "12月13日", "12月14日", "12月15日", "12月16日", "12月17日", "12月18日", "12月19日", "12月20日", "12月21日", "12月22日", "12月23日", "12月24日", "12月25日", "12月26日", "12月27日", "12月28日", "12月29日", "12月30日", "12月31日"}
	//datem := []string{"11月5日"}
	for _, d := range datem {
		log.Println("开始抓取", d)
		J := setdaysfromwiki(d)
		// 创建文件
		filePtr, err := os.Create("./data/" + d + ".json")
		if err != nil {
			log.Println("文件创建失败", err.Error())
			return
		}
		// 创建Json编码器
		encoder := json.NewEncoder(filePtr)
		err = encoder.Encode(J)
		if err != nil {
			log.Println("编码错误", err.Error())
		}
		defer filePtr.Close()
	}
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
				var newTD History
				if h1.Children().Is("ul") {
					h1.Children().Children().Each(func(_ int, li *goquery.Selection) {
						log.Println("条目格式化", MYTXT)
						if daytype != "0" {
							newTD = History{
								Type:    daytype,
								Year:    strings.Split(MYTXT, "：")[0],
								Contant: li.Text(),
							}
							JF = append(JF, newTD)
						}
					})
				} else if h1.Children().Nodes != nil {
					log.Println("条目格式化", MYTXT)
					if daytype != "0" {
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

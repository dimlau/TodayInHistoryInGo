package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/tidwall/buntdb"
)

type (
	History struct {
		Type    string `json:"type"`
		Contant string `json:"contant"`
	}

	CWer struct {
		B   *colly.Collector
		URL string
	}

	CW interface {
		getalldays()
		gettodayinwiki(date string)
	}
)

func main() {

	dbbotAC, err := buntdb.Open("./tih.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbbotAC.Close()

	C := new(CWer)
	C.URL = "https://zh.wikipedia.org/zh-cn/"
	C.B = colly.NewCollector()
	C.B.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	datem := C.getalldays()
	fmt.Println(datem)
	for _, date := range datem {
		fmt.Println("开始写入", date)
		dbbotAC.Update(func(tx *buntdb.Tx) error {
			_, err := tx.Get(date)
			if err == buntdb.ErrNotFound {
				JF := C.gettodayinwiki(date)
				value, _ := json.Marshal(JF)
				_, _, err = tx.Set(date, string(value), nil)
			}
			return err
		})
	}

	dbbotAC.Shrink()
	//Ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	//fmt.Println(JF[Ra.Intn(len(JF))].Contant)
}

func (C CWer) getalldays() []string {
	datem := []string{"1月1日"}

	C.URL = C.URL + "1月1日"

	C.B.OnHTML(".navbox", func(e *colly.HTMLElement) {
		INX := 99999
		e.ForEach("li a", func(i int, h *colly.HTMLElement) {

			if h.Attr("class") == "mw-selflink selflink" {
				INX = i
			}
			if i > INX {

				nex, _ := h.DOM.Attr("title")
				datem = append(datem, nex)
			}
		})

	})

	C.B.Visit(C.URL)
	return datem
}

func (C CWer) gettodayinwiki(date string) []*History {
	JF := []*History{}
	all := 0

	C.URL = C.URL + date

	C.B.OnHTML("#toc ul ul", func(e *colly.HTMLElement) {
		all = e.DOM.Children().Length()
	})
	C.B.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {
		e.ForEach("ul", func(i int, h *colly.HTMLElement) {
			h.ForEach("sup", func(i int, k *colly.HTMLElement) {
				k.DOM.Remove()
			})
			if i > 1 && (i <= all+1) {
				h.ForEach("li", func(_ int, k *colly.HTMLElement) {
					JF = append(JF, &History{Type: "1", Contant: k.Text})
				})
			}
			if i == all+2 {
				h.ForEach("li", func(_ int, k *colly.HTMLElement) {
					JF = append(JF, &History{Type: "2", Contant: k.Text})
				})
			}
			if i == all+3 {
				h.ForEach("li", func(_ int, k *colly.HTMLElement) {
					JF = append(JF, &History{Type: "3", Contant: k.Text})
				})
			}
		})
	})

	C.B.Visit(C.URL)
	return JF
}

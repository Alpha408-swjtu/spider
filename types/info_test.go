package types

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestInfo(t *testing.T) {
	url := "https://movie.douban.com/top250?start=75&filter="
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 Edg/130.0.0.0")
	resp, _ := client.Do(req)

	defer resp.Body.Close()
	docDetial, _ := goquery.NewDocumentFromReader(resp.Body)
	s := docDetial.Find("#content > div > div.article > ol > li:nth-child(14) > div > div.info > div.bd > p:nth-child(1)").Text()
	fmt.Println(s)
}

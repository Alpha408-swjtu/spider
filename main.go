package main

import (
	"fmt"
	"net/http"
	"spider/types"
	"spider/utils"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const (
	MOVIES       = 250
	EACHP        = 25
	UrlPrefix    = "https://movie.douban.com/top250?start="
	SelectPrefix = "#content > div > div.article > ol > li:nth-child("
	NameBack     = ") > div > div.info > div.hd > a > span:nth-child(1)"
	InfoBack     = ") > div > div.info > div.bd > p:nth-child(1)"
)

var (
	MovieCh    = make(chan types.Movie, 300)
	client     *http.Client
	clientOnce sync.Once
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetName(page int) {
	if page%25 != 0 && page != 0 {
		panic("参数非法")
	}
	if client == nil {
		clientOnce.Do(func() {
			client = &http.Client{}
		})
	}
	url := UrlPrefix + strconv.Itoa(page) + "&filter="
	req, err := http.NewRequest("GET", url, nil)
	HandleErr(err)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 Edg/130.0.0.0")
	resp, err := client.Do(req)
	HandleErr(err)
	defer resp.Body.Close()
	docDetial, err := goquery.NewDocumentFromReader(resp.Body)
	HandleErr(err)

	wg := sync.WaitGroup{}
	wg.Add(EACHP)
	for i := 1; i <= EACHP; i++ {
		go func(int) {
			defer wg.Done()
			nameString := SelectPrefix + strconv.Itoa(i) + NameBack
			InfoString := SelectPrefix + strconv.Itoa(i) + InfoBack
			name := docDetial.Find(nameString).Text()
			info := docDetial.Find(InfoString).Text()
			dir, actor, year := utils.SpliteInfo(info)
			movie := types.Movie{
				Name:     name,
				Year:     year,
				Director: dir,
				Actor:    actor,
			}
			MovieCh <- movie
		}(i)
	}
	wg.Wait()
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(MOVIES / EACHP)
	for i := 0; i < MOVIES; i += EACHP {
		go func(int) {
			defer wg.Done()
			GetName(i)
		}(i)
	}
	wg.Wait()
	close(MovieCh)
	nameS := make([]types.Movie, 0, 500)
	tempCh := make(chan int)
	go func() {
		for {
			name, ok := <-MovieCh
			nameS = append(nameS, name)
			if !ok {
				break
			}
		}
		tempCh <- 1
	}()
	<-tempCh
	for _, movie := range nameS {
		fmt.Println(movie)
	}
}

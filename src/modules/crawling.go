package src

import (
	"context"
	model "crawling.com/models"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"strconv"
	"strings"
	"time"
)

func Crawling() string {
	fmt.Println("crawling")
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	pageCount := 500

	attr_ := make([]map[string]string, 0)

	//var example string
	err := chromedp.Run(ctx,
		chromedp.Navigate(loginUrl),
		chromedp.AttributesAll("#frmLogin input", &attr_, chromedp.ByQueryAll),
		//chromedp.Text(`div`, &example),
	)
	if err != nil {
		log.Fatal(err)
	}
	var list1 []string
	for _, val := range attr_ {
		list1 = append(list1, val["placeholder"])
	}

	var res string
	err = chromedp.Run(ctx, doLogin(loginUrl, &res))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("login suc")
	var infoList []model.GoodsInfo

	var totalCount string
	chromedp.Run(ctx,
		chromedp.Navigate(goodsUrl),
		chromedp.Text(`//*[@id="frmSearchGoods"]/div[4]/div[1]/strong[1]`, &totalCount),
	)

	t, _ := strconv.Atoi(strings.Replace(totalCount, ",", "", -1))
	totalPage := t / pageCount
	for pn := 1; pn <= totalPage; pn++ {
		// 먼저 이동
		// page 이동
		if pn != 1 {
			chromedp.Run(ctx,
				chromedp.Navigate(strings.Replace(goodsUrl, "?page=1&&", "?page="+strconv.Itoa(pn)+"&&", 1)),
			)
		}

		for i := 1; i <= 500; i++ {
			var info model.GoodsInfo
			err = chromedp.Run(ctx, getGoodsInfo(i, &info))
			if err != nil {
				log.Fatal(err)
				break
			}
			infoList = append(infoList, info)
		}
		log.Printf("( %d / %d ) ", pn, totalPage)
	}

	err = chromedp.Run(ctx,
		chromedp.Text(`//*[@id="menu_goods"]`, &res),
	)

	log.Printf("clear")
	//log.Printf(strVar)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("got: `%s`", strings.TrimSpace(res))

	//if err != nil {
	//   panic(err)
	//}

	attr := make([]map[string]string, 0)
	attr2 := make([]map[string]string, 0)

	err = chromedp.Run(ctx,
		chromedp.AttributesAll("#container #header nav div[class=container] ul li a", &attr, chromedp.ByQueryAll),
		chromedp.AttributesAll("#container-wrap #container #header nav ul li a", &attr, chromedp.ByQueryAll),
		chromedp.AttributesAll("#primary ytd-rich-grid-renderer div#contents ytd-rich-grid-row div#contents ytd-rich-item-renderer #video-title-link", &attr, chromedp.ByQueryAll),
		//chromedp.AttributesAll("#primary ytd-rich-grid-renderer div#contents ytd-rich-grid-row div#contents ytd-rich-item-renderer #thumbnail yt-image img", &attr2, chromedp.ByQueryAll),
	)

	var linklist []string
	var titleList []string
	var valueList []string
	var imgUrlList []string
	for _, val := range attr {
		linklist = append(linklist, val["href"])
		titleList = append(titleList, val["text"])
		valueList = append(valueList, val["value"])
		//titleList = append(titleList, val["aria-label"])
	}

	for _, val := range attr2 {
		imgUrlList = append(imgUrlList, val["src"])
	}
	fmt.Println(len(linklist))
	fmt.Println(len(titleList))
	fmt.Println(len(imgUrlList))

	return "err"

}

type WebDriver struct {
	ctx    context.Context
	cancel context.CancelFunc
	Str    string // 데이터 반환
}

func checkErr(err error) {

}

func getGoodsInfo(count int, info *model.GoodsInfo) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Evaluate(`document.querySelector("#frmList > div.table-responsive > table > tbody > tr:nth-child(`+strconv.Itoa(count)+`) > td.width-2xs.center > a").href`, &info.Img),
	}
}
func doLogin(loginUrl string, res *string) chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.Navigate(loginUrl),
		chromedp.Sleep(4 * time.Second),
		chromedp.Text(`html`, res),
	}
}

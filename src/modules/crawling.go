package modules

import (
	"context"
	model "crawling.com/models"
	"github.com/chromedp/chromedp"
	"log"
	"strconv"
	"strings"
	"time"
)

func Crawling(e *model.Environment) []model.GoodsInfo {

	log.Println("Crawling start")
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	pageCount := 500

	loginUrl := e.BaseUrl + e.Login.Url
	goodsUrl := e.BaseUrl + e.Goods.Url

	var res string
	id := e.Login.Id
	pwd := e.Login.Pwd

	err := chromedp.Run(ctx, doLogin(loginUrl, id, pwd, &res))
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

	// 먼저 이동
	t, _ := strconv.Atoi(strings.Replace(totalCount, ",", "", -1))
	totalPage := t / pageCount
	totalPage = 1
	for pn := 1; pn <= totalPage; pn++ {
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
			//log.Printf("( %d / %d ) : %s, %s, %s, %s, %s ", (pn-1)*pageCount+i, t, info.IdNum, info.Brand, info.Supplier, info.GoodsCode, info.Img)
			infoList = append(infoList, info)
		}
		log.Printf("( %d / %d ) ", pn, totalPage)
	}

	return infoList

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
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[2]", &info.IdNum),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[3]", &info.Supplier),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[4]", &info.GoodsCode),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[5]", &info.Brand),
		chromedp.Evaluate(`document.querySelector("#frmList > div.table-responsive > table > tbody > tr:nth-child(1) > td.width-2xs.center > a > img").getAttribute('src')`, &info.Img),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[9]/div[1]/a", &info.Name),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[10]", &info.RetailPrice),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[11]", &info.SellingPrice),
	}
}
func doLogin(loginUrl string, id string, pwd string, res *string) chromedp.Tasks {
	log.Printf("login try")
	return chromedp.Tasks{
		//chromedp.Evaluate(`window.scrollTo(0,document.querySelector("body ytd-app div#content").clientHeight); document.querySelector("body ytd-app div#content").clientHeight;`, &newHeight),
		chromedp.Navigate(loginUrl),
		//chromedp.WaitVisible(`body > footer`),
		//chromedp.WaitVisible(`#frmLogin #login`, chromedp.ByID),
		//chromedp.SendKeys(`#frmLogin #login`, id, chromedp.ByID),
		chromedp.SendKeys(`//*[@id="login"]`, id),
		chromedp.SendKeys(`//*[@id="frmLogin"]/table/tbody/tr/td/div/div[1]/div[1]/div[2]/input`, pwd),
		chromedp.Click(`//*[@id="frmLogin"]/table/tbody/tr/td/div/div[1]/div[2]/input`),
		//chromedp.Click(`//*[@id="loginFrm"]/button`),
		chromedp.Sleep(4 * time.Second),
		//chromedp.Text(`/html/body/div[2]/div/div[1]/div/span`, res),
		chromedp.Text(`html`, res),
	}
}

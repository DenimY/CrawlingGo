package modules

import (
	"context"
	model "crawling.com/models"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var LoginUrl string
var GoodsUrl string
var PageCount = 500
var Done int

func Crawling(e *model.Environment) []model.GoodsInfo {

	log.Println("Crawling start")
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	PageCount = 500

	LoginUrl = e.BaseUrl + e.Login.Url
	GoodsUrl = e.BaseUrl + e.Goods.Url

	id := e.Login.Id
	pwd := e.Login.Pwd

	log.Printf("login start")
	err := chromedp.Run(ctx, doLogin(LoginUrl, id, pwd))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("login suc")
	var infoList []model.GoodsInfo

	// 먼저 이동
	var totalCount string
	chromedp.Run(ctx,
		chromedp.Navigate(GoodsUrl),
		chromedp.WaitVisible(`//*[@id="frmList"]/div[2]/table`),
		chromedp.Text(`//*[@id="frmSearchGoods"]/div[4]/div[1]/strong[1]`, &totalCount),
	)

	t, _ := strconv.Atoi(strings.Replace(totalCount, ",", "", -1))

	/*slicePage := (t / 500)
	go func(msg int) {
		defer wait.Done() //끝나면 .Done() 호출
		fmt.Println(msg)
	}(slicePage)

	go func() {
		defer wait.Done() //끝나면 .Done() 호출
		fmt.Println("Hello")
	}()*/

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Printf(string(runtime.NumCPU()))

	var wait sync.WaitGroup

	totalPage := int(math.Ceil(float64(t) / float64(PageCount)))
	//totalPage = 1 //test code
	excelFileName := time.Now().Format("20060102150405")
	sliceCnt := 1
	goSlice := int(math.Floor(float64(totalPage) / float64(sliceCnt)))

	fmt.Println(goSlice)
	//startPn := 1
	//endPn := 0
	wait.Add(sliceCnt)

	for j := 0; j < sliceCnt; j++ {
		go func(n int) {
			defer wait.Done()

			/*log.Printf("Routine-%d \n", n)
			ctxn, cancel := chromedp.NewContext(
				context.Background(),
			)
			defer cancel()

			PageCount = 500

			LoginUrl = e.BaseUrl + e.Login.Url
			GoodsUrl = e.BaseUrl + e.Goods.Url

			log.Printf("Routine-%d login try\n", n)
			err := chromedp.Run(ctxn, doLogin(LoginUrl, id, pwd))
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Routine-%d login suc\n", n)*/

			//

			startPn := n * (goSlice)
			if startPn < 1 {
				startPn = 1
			}
			endPn := startPn + (goSlice - 1)
			//fmt.Println(strconv.Itoa(startPn) + "-" + strconv.Itoa(endPn))
			//fmt.Println(n)

			log.Printf("Routine-%d crawlingGo\n", n)
			crawlingGo(ctx, startPn, endPn, n, fmt.Sprintf(excelFileName+"_", strconv.Itoa(startPn), "-", strconv.Itoa(endPn)))

		}(j)
	}

	wait.Wait()

	log.Println("clear go routine")

	/*routineCnt = 1
	pn := 1 // start page
	for ; pn <= totalPage; pn++ {
		le := routineCnt * (totalPage / 12)
		// page 이동
		if pn != 1 {
			chromedp.Run(ctx,
				chromedp.Navigate(strings.Replace(GoodsUrl, "?page=1&&", "?page="+strconv.Itoa(pn)+"&&", 1)),
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
			log.Printf("( %d / %d )", (pn-1)*PageCount+i, t)
			infoList = append(infoList, info)
		}
		log.Printf("( %d / %d ) ", pn, totalPage)

		if le == (pn) || pn == totalPage {
			GoodsOutput(infoList, fmt.Sprintf(excelFileName+"_", saveCnt))
			// 초기화
			infoList = []model.GoodsInfo{}
			routineCnt++
		}
	}*/

	return infoList

}

func crawlingGo(ctx context.Context, startPage, endPage, routineCnt int, excelFileName string) error {
	pn := startPage
	saveCnt := 1
	//t, _ := strconv.Atoi(strings.Replace(totalCount, ",", "", -1))
	var infoList []model.GoodsInfo
	for ; pn <= endPage; pn++ {
		// page 이동
		url := strings.Replace(GoodsUrl, "?page=1&&", "?page="+strconv.Itoa(pn)+"&&", 1)
		err := chromedp.Run(ctx,
			chromedp.WaitVisible(`//*[@id="frmList"]/div[2]/table`),
			chromedp.Navigate(strings.Replace(GoodsUrl, "?page=1&&", "?page="+strconv.Itoa(pn)+"&&", 1)),
		)

		if err != nil {
			log.Printf("navigate err = %v\n: %v\n", url, err)
			log.Printf("navigate err = %v: %v\n", startPage, err)
			return nil
		}

		for i := 1; i <= 500; i++ {
			var info model.GoodsInfo
			err := chromedp.Run(ctx, getGoodsInfo(i, &info))
			if err != nil {
				log.Fatal(err)
				break
			}
			//log.Printf("( %d / %d ) : %s, %s, %s, %s, %s ", (pn-1)*pageCount+i, t, info.IdNum, info.Brand, info.Supplier, info.GoodsCode, info.Img)
			log.Printf("Routine-%d ( %d / %d )", routineCnt, (pn-1)*PageCount+i, endPage)
			infoList = append(infoList, info)
		}
		log.Printf("Routine-%d : ( %d / %d ) ", routineCnt, pn, endPage)

		//le := saveCnt * (endPage / 50)
		//if le == (pn) || pn == endPage {
		GoodsOutput(infoList, fmt.Sprintf(excelFileName+"_", saveCnt))
		// 초기화
		infoList = []model.GoodsInfo{}
		saveCnt++
		//}

		//GoodsOutput(infoList, fmt.Sprintf(excelFileName+"_", routineCnt))

	}

	return nil
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
		chromedp.WaitVisible(`//*[@id="frmList"]/div[2]/table`),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[2]", &info.IdNum),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[3]", &info.Supplier),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[4]", &info.GoodsCode),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[5]", &info.Brand),
		chromedp.Evaluate(`document.querySelector("#frmList > div.table-responsive > table > tbody > tr:nth-child(`+strconv.Itoa(count)+`) > td.width-2xs.center > a > img").getAttribute('src')`, &info.Img),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[9]/div[1]/a", &info.Name),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[10]", &info.RetailPrice),
		chromedp.Text(`//*[@id="frmList"]/div[2]/table/tbody/tr[`+strconv.Itoa(count)+"]/td[11]", &info.SellingPrice),
	}
}
func doLogin(loginUrl string, id string, pwd string) chromedp.Tasks {
	return chromedp.Tasks{
		//chromedp.Evaluate(`window.scrollTo(0,document.querySelector("body ytd-app div#content").clientHeight); document.querySelector("body ytd-app div#content").clientHeight;`, &newHeight),
		chromedp.Navigate(loginUrl),
		//chromedp.WaitVisible(`body`),
		//chromedp.WaitVisible(`#frmLogin #login`, chromedp.ByID),
		//chromedp.SendKeys(`#frmLogin #login`, id, chromedp.ByID),
		chromedp.SendKeys(`//*[@id="login"]`, id),
		chromedp.SendKeys(`//*[@id="frmLogin"]/table/tbody/tr/td/div/div[1]/div[1]/div[2]/input`, pwd),
		chromedp.Click(`//*[@id="frmLogin"]/table/tbody/tr/td/div/div[1]/div[2]/input`),
		//chromedp.Click(`//*[@id="loginFrm"]/button`),
		chromedp.Sleep(5 * time.Second),
		//chromedp.Text(`/html/body/div[2]/div/div[1]/div/span`, res),
	}
}

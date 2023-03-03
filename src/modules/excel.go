package modules

import (
	"bytes"
	model "crawling.com/models"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func excel() {

}
func GoodsOutput(dataList []model.GoodsInfo, fileName string) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	sheetName := "Sheet1"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("start start excel")
	cellCount := 0
	for _, data := range dataList {
		cellCount++
		// 12개로 분할하여 저장
		//f.SetColWidth("Sheet1", "E", "H", 20)
		f.SetRowHeight(sheetName, cellCount, 50)
		f.SetCellValue(sheetName, fmt.Sprint("A", cellCount), data.IdNum)
		f.SetCellValue(sheetName, fmt.Sprint("B", cellCount), data.Supplier)
		f.SetCellValue(sheetName, fmt.Sprint("C", cellCount), data.GoodsCode)
		f.SetCellValue(sheetName, fmt.Sprint("D", cellCount), data.Brand)
		addPicture(f, sheetName, data.GoodsCode, fmt.Sprint("E", cellCount), data.Img)
		f.SetCellValue(sheetName, fmt.Sprint("F", cellCount), data.Img)
		f.SetCellValue(sheetName, fmt.Sprint("G", cellCount), data.Name)
		f.SetCellValue(sheetName, fmt.Sprint("H", cellCount), data.RetailPrice)
		f.SetCellValue(sheetName, fmt.Sprint("U", cellCount), data.SellingPrice)
	}
	//통합 문서에 대 한 기본 워크시트를 설정 합니다

	e := saveFile(f, fileName)
	if e != nil {
		log.Printf("excel save failure : %s\n", e)
	} else {
		log.Printf("save excel : %v", fileName)
	}

}

func GoodsOutputSlice(dataList []model.GoodsInfo) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	sheetName := "Sheet1"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("start start excel")
	count := 1
	cellCount := 0
	fullLen := len(dataList)
	for index, data := range dataList {
		// 12개로 분할하여 저장
		le := count * (fullLen / 12)
		cellCount++
		//f.SetColWidth("Sheet1", "E", "H", 20)
		f.SetRowHeight(sheetName, cellCount, 50)
		f.SetCellValue(sheetName, fmt.Sprint("A", cellCount), data.IdNum)
		f.SetCellValue(sheetName, fmt.Sprint("B", cellCount), data.Supplier)
		f.SetCellValue(sheetName, fmt.Sprint("C", cellCount), data.GoodsCode)
		f.SetCellValue(sheetName, fmt.Sprint("D", cellCount), data.Brand)
		addPicture(f, sheetName, data.GoodsCode, fmt.Sprint("E", cellCount), data.Img)
		f.SetCellValue(sheetName, fmt.Sprint("F", cellCount), data.Img)
		f.SetCellValue(sheetName, fmt.Sprint("G", cellCount), data.Name)
		f.SetCellValue(sheetName, fmt.Sprint("H", cellCount), data.RetailPrice)
		f.SetCellValue(sheetName, fmt.Sprint("U", cellCount), data.SellingPrice)

		if le == (index+1) || index+1 == fullLen {
			cellCount = 0
			f.SetActiveSheet(index)
			e := saveFile(f, fmt.Sprint(time.Now().Format("20060102150405")+"_", count))
			if e != nil {
				log.Printf("excel save failure : %s\n", e)
			}
			f = excelize.NewFile()
			defer func() {
				if err := f.Close(); err != nil {
					log.Println(err)
				}
			}()
			_, err = f.NewSheet(sheetName)
			if err != nil {
				fmt.Println(err)
				return
			}
			count++
		}
	}
	//통합 문서에 대 한 기본 워크시트를 설정 합니다

}
func uploadFile() {

}

func updateFile() {

}

func saveFile(file *excelize.File, fileName string) error {
	return file.SaveAs(fileName + ".xlsx")
}

func buildFileName() {
	fullUrlFile := "https://i.imgur.com/m1UIjW1.jpg"
	fileUrl, e := url.Parse(fullUrlFile)
	if e != nil {
		panic(e)
	}

	path := fileUrl.Path
	segments := strings.Split(path, "/")

	fileName := segments[len(segments)-1]
	println(fileName)
}

/*
	func testpick() error {
		//fullUrlFile := "https://i.imgur.com/m1UIjW1.jpg"
		fullUrlFile := "https://pjh5023.cdn-nhncommerce.com/data/goods/23/02/09/1000326716/t50_1000133769_list_022.jpg"
		fileName := "m1UIjW1.jpg"

		r, e := http.Get(fullUrlFile)
		if e != nil {
			panic(e)
		}
		defer r.Body.Close()
		//buildFileName()
		// Create distination
		f, e := os.Create(fileName) // "m1UIjW1.jpg"
		if e != nil {
			panic(e)
		}
		defer f.Close()
		// Fill distination with content
		n, e := f.ReadFrom(r.Body)
		if e != nil {
			panic(e)
		}
		fmt.Println("File size: ", n)

		return nil
	}
*/
func addPicture(f *excelize.File, sheet, imgName, cell, imgUrl string) error {

	extension := imgUrl[strings.LastIndex(imgUrl, "."):len(imgUrl)]
	r, err := http.Get(imgUrl)
	if err != nil {
		log.Printf("http.Get err: %v ", err)
		panic(err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	buf := &bytes.Buffer{}
	//buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(r.Body)
	if err != nil {
		fmt.Println("File size1: ", n)
		fmt.Println("buf readForm err: ", err)
	}

	if err := f.AddPictureFromBytes(sheet, cell, imgName, strings.ToLower(extension), buf.Bytes(),
		&excelize.GraphicOptions{
			AutoFit: true,
			//ScaleX:  0.1,
			//ScaleY:  0.1,
			OffsetX: 15,
			OffsetY: 10,
			//Hyperlink:     imgUrl,
			//HyperlinkType: "External",
			Positioning: "oneCell",
		}); err != nil {
		log.Printf("AddPictureFromBytes err:  %v ", err)
	}

	return nil
}

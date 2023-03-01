package src

import (
	model "crawling.com/models"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
)

func excel() {

}
func goodsOutput(dataList []model.GoodsInfo) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}

	for index, data := range dataList {
		cellNum := strconv.Itoa(index)
		f.SetCellValue(sheetName, "A"+cellNum, data.IdNum)
		f.SetCellValue(sheetName, "B"+cellNum, data.IdNum)
		f.SetCellValue(sheetName, "C"+cellNum, data.IdNum)
		f.SetCellValue(sheetName, "D"+cellNum, data.IdNum)
		println(index, data)
	}
	//통합 문서에 대 한 기본 워크시트를 설정 합니다
	f.SetActiveSheet(index)
	e := saveFile(f, "newFile")

	if e != nil {
		log.Printf("excel save failure : %s\n", e)
	}
}
func uploadFile() {

}

func updateFile() {

}

func saveFile(file *excelize.File, fileName string) error {
	return file.SaveAs(fileName + ".xlsx")
}

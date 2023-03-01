package main

import (
	model "crawling.com/models"
	m "crawling.com/modules"
	"log"
	"strconv"
)

func main() {

	log.Println("start Main()")
	//m.Crawling()

	// test goods excel download
	var dataList []model.GoodsInfo
	i := 1
	for i < 100 {
		n := strconv.Itoa(i)
		dataList = append(dataList, model.GoodsInfo{
			IdNum:     "id" + n,
			Supplier:  "Supplier" + n,
			GoodsCode: "GoodsCode" + n,
			Brand:     "Brand" + n,
			Img:       "Img" + n,
		})
		i++
	}

	m.GoodsOutput(dataList)

}

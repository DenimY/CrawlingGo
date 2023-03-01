package main

import (
	m "crawling.com/modules"
	"log"
)

func main() {

	env, err := m.LoadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("start Main()")
	dataList := m.Crawling(env)

	// test goods excel download
	/*var dataList []model.GoodsInfo
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
	}*/

	m.GoodsOutput(dataList)

}

package models

type Environment struct {
	BaseUrl string `yaml:"baseUrl"`
	Login   Login  `yaml:"login"`
	Goods   Goods  `yaml:"goods"`
}

type Login struct {
	Url string `yaml:"url"`
	Id  string `yaml:"id"`
	Pwd string `yaml:"pwd"`
}

type Goods struct {
	Url string `yaml:"url"`
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

var (
	types = []vehicleType{
		{Name: "Carros", Url: "https://tabelacarros.com/marcas/carros", Brands: []brand{}},
		{Name: "Motos", Url: "https://tabelacarros.com/tabela-fipe-motos", Brands: []brand{}},
		{Name: "Caminhao", Url: "https://tabelacarros.com/marcas/caminhoes", Brands: []brand{}},
	}
)

type vehicleType struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Brands []brand
}

type brand struct {
	Name   string   `json:"name"`
	Url    string   `json:"url"`
	Models []string `json:"models"`
}

func main() {

	fmt.Println("Starting...")

	for typeIndex, vehicleType := range types {
		brands := fetchBrands(vehicleType.Url)
		types[typeIndex].Brands = brands

		for i, brand := range brands {
			models := fetchModels(brand.Url)
			types[typeIndex].Brands[i].Models = models
		}
	}

	// for i, brand := range brands {
	// 	models := fetchModels(brand.Url)
	// 	brands[i].Models = models
	// }

	saveBrandsAsJSON(types, "webscrap.json")

}

func fetchBrands(url string) []brand {

	brands := []brand{}

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML(".botao_fake1", func(e *colly.HTMLElement) {

		sanitized, err := sanitizeString(e.Text)

		if err != nil {
			log.Printf("Model %s is not valid. \n", e.Text)
		} else {
			brands = append(brands, brand{
				Name: sanitized,
				Url:  e.Attr("href"),
			})
		}

	})

	c.Visit(url)

	return brands
}

func fetchModels(url string) []string {
	models := []string{}

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML(".modelo_base2", func(e *colly.HTMLElement) {

		sanitized, err := sanitizeString(e.Text)

		if err != nil {
			log.Printf("Model %s is not valid. \n", e.Text)
		}

		models = append(models, sanitized)

	})

	c.Visit(url)

	return models
}

func sanitizeString(s string) (string, error) {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.Trim(s, " ")
	if s == "" {
		return "", errors.New("string is empty")
	}

	return s, nil
}

func saveBrandsAsJSON(brands []vehicleType, filename string) error {
	data, err := json.MarshalIndent(brands, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("JSON Data: %s\n", data)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

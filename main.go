package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Typer interface {
	GetResult()
}

type JsonType struct {
	Product
	path string
}

type CsvType struct {
	Product
	path string
}

type Product struct {
	Product string  `json:"product"`
	Price   float64 `json:"price"`
	Rating  int     `json:"rating"`
}

func GetType(file string) Typer {
	if strings.ToLower(filepath.Ext(file)) == ".json" {
		return &JsonType{path: file}
	} else if strings.ToLower(filepath.Ext(file)) == ".csv" {
		return &CsvType{path: file}
	} else {
		log.Fatalf("this format is not provided")
		return nil
	}

}

func (c *CsvType) GetResult() {
	file, err := os.Open("db.csv")
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(file)
	var rating, highest CsvType
	for {
		eachRecord, err := reader.Read()
		if err != nil || err == io.EOF {
			break
		}
		var productCSV CsvType
		productCSV.Product.Product = eachRecord[0]
		productCSV.Price, _ = strconv.ParseFloat(eachRecord[1], 64)
		productCSV.Rating, _ = strconv.Atoi(eachRecord[2])
		if productCSV.Price >= highest.Price {
			highest = productCSV
		}
		if productCSV.Rating >= rating.Rating {
			rating = productCSV
		}
	}
	fmt.Printf("item in [%s] with highest price is [product:%s price:%f rating:%d] \n", c.path, highest.Product.Product, highest.Product.Price, highest.Product.Rating)
	fmt.Printf("item in [%s] with highest rating is [product:%s price:%f rating:%d] \n", c.path, rating.Product.Product, rating.Product.Price, rating.Product.Rating)
}
func (j *JsonType) GetResult() {
	var highest JsonType
	var rating JsonType

	fileName := j.path
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error to read [file=%v]: %v", fileName, err.Error())
	}

	if err != nil {
		log.Fatalf("Could not obtain stat, handle error: %v", err.Error())
	}

	r := bufio.NewReader(f)
	d := json.NewDecoder(r)

	i := 0

	d.Token()
	for d.More() {
		elm := &JsonType{}
		d.Decode(elm)
		if elm.Price >= highest.Price {
			highest = *elm
		}
		if elm.Rating >= rating.Rating {
			rating = *elm
		}
		i++
	}
	d.Token()
	fmt.Printf("item in [%s] with highest price is [product:%s price:%f rating:%d] \n", j.path, highest.Product.Product, highest.Product.Price, highest.Product.Rating)
	fmt.Printf("item in [%s] with highest rating is [product:%s price:%f rating:%d] \n", j.path, rating.Product.Product, rating.Product.Price, rating.Product.Rating)
}

func main() {
	var path string
	flag.StringVar(&path, "p", "./", "path to file")
	flag.Parse()

	format := GetType(path)
	format.GetResult()
}

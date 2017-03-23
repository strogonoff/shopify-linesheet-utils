package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type ProductSet struct {
	handle string
	name   string
}

type ProductType string

type ProductColorVariant struct {
	sku         string
	color       string
	picturePath string
}

type Product struct {
	ProductColorVariant

	set ProductSet

	productType ProductType

	otherColors []ProductColorVariant

	wholesalePrice string
	retailPrice    string
}

func startsValidSet(record []string) bool {
	return strings.Contains(record[1], " - ") == false
}

func startsNewSet(record []string) bool {
	if record[1] == "" {
		return false
	}
	return true
}

func main() {
	const IN_FILENAME = "shopifydata.csv"
	const OUT_FILENAME = "indesigndata.csv"

	fmt.Printf("Looking for CSV file called %s...\n", IN_FILENAME)

	fileContents, err := ioutil.ReadFile(IN_FILENAME)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Reading CSV...")

	csvReader := csv.NewReader(strings.NewReader(string(fileContents)))

	var recordCount int = -1 // Begin at -1 to compensate for header line

	var products []Product

	var validSet bool = false
	var p Product
	var ps ProductSet
	var pt ProductType
	var pc ProductColorVariant

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		recordCount = recordCount + 1

		if recordCount < 1 {
			continue
		}

		newSet := startsNewSet(record)

		if newSet {
			validSet = startsValidSet(record)
		}

		if !validSet {
			continue
		}

		newProduct := newSet || pt != ProductType(record[8])

		if newProduct && recordCount > 1 {
			products = append(products, p)
		}

		if newSet {
			ps = ProductSet{handle: record[0], name: record[1]}
		}

		pc = ProductColorVariant{
			sku:   record[13],
			color: record[10],
		}

		if newProduct {
			pt = ProductType(record[8])
			p = Product{
				ProductColorVariant: pc,
				set:         ps,
				productType: pt,
				wholesalePrice: record[19],
			}
		} else {
			p.otherColors = append(p.otherColors, pc)
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	products = append(products, p)

	fmt.Printf("Processed %d records\n", recordCount)

	fmt.Printf("Found %d products\n", len(products))

	fmt.Printf("Writing CSV into %s\n", OUT_FILENAME)

	f, err := os.Create(OUT_FILENAME)
	defer f.Close()

	csvWriter := csv.NewWriter(f)

	for _, product := range products {
		record := []string{
			product.set.name,
			string(product.productType),
			product.wholesalePrice,
			product.sku,
			product.color,
		}

		if len(product.otherColors) > 0 {
			for _, c := range product.otherColors {
				record = append(
					record, 
					c.sku,
					c.color)
			}
		}

		if err := csvWriter.Write(record); err != nil {
			log.Fatalln("Error writing product to CSV:", err)
		}
	}
}

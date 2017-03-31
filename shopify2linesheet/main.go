package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) != 4 {
		log.Fatal("Exactly four arguments are expected: in file, out file, discount factor and asset root path")
	}

	var IN_FILENAME = args[0]
	var OUT_FILENAME = args[1]
	var ASSET_ROOT_PATH = args[3]

	_discount := args[2]

	WHOLESALE_DISCOUNT_FACTOR, err := strconv.ParseFloat(_discount, 32)

	if err != nil {
		log.Fatal("Unable to convert specified discount factor to floating point type")
	}

	log.Println(fmt.Sprintf("Reading Shopify records from %s...", IN_FILENAME))

	fileContents, err := ioutil.ReadFile(IN_FILENAME)

	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(strings.NewReader(string(fileContents)))

	// Below: big not-so-elegant block of instructions.
	// Iterates over CSV entries and constructs product sets.
	// Keeps some state in variables.

	var shopifyEntryCount int = -1 // begin at -1 to compensate for CSV header

	var pSets []ProductSet

	var curRecord ShopifyRecord
	var oldRecord ShopifyRecord

	var curSet ProductSet
	var processingSet bool = false

	var curProduct Product

	dlQueue := make(map[string]string)

	for {
		values, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		shopifyEntryCount += 1

		_headerRow := shopifyEntryCount < 1

		if _headerRow {
			continue
		}

		oldRecord = curRecord
		curRecord = CSVEntry(values).ShopifyRecord()

		_startsEntries := shopifyEntryCount <= 1

		// New set

		_startsSet := curRecord.handle != oldRecord.handle

		if _startsSet {
			processingSet = isValidSetTitle(curRecord.title)

			if processingSet {
				if !_startsEntries {
					curSet.products = append(curSet.products, curProduct)
					pSets = append(pSets, curSet)
				}

				curSet = curRecord.ProductSet()

				if curRecord.imageSrc != "" {
					path := filepath.Join(
						ASSET_ROOT_PATH,
						SuggestFilename(curRecord.imageSrc))
					dlQueue[curRecord.imageSrc] = path
					curSet.picturePath = path
				}

				//log.Println("started set", curSet.name)
			} else {
				//log.Println("skipped set", curRecord.title)
			}
		}

		if !processingSet {
			continue
		}

		// New product

		_startsProduct := _startsSet || curRecord.oProductType != oldRecord.oProductType

		if _startsProduct {
			if !_startsSet {
				curSet.products = append(curSet.products, curProduct)
			}
			curProduct = curRecord.Product(WHOLESALE_DISCOUNT_FACTOR)
			//log.Println("-- added product", curProduct.name)
		}

		// New variant

		curVariant := curRecord.ProductVariant()

		if curRecord.vImage != "" {
			path := filepath.Join(
				ASSET_ROOT_PATH,
				SuggestFilename(curRecord.vImage))
			dlQueue[curRecord.vImage] = path
			curVariant.picturePath = path
		}

		curProduct.variants = append(curProduct.variants, curVariant)
		//log.Println("-- added variant", curVariant.sku)
	}

	// The above loop leaves things hanging
	curSet.products = append(curSet.products, curProduct)
	pSets = append(pSets, curSet)

	log.Println(fmt.Sprintf("Read %d records.", shopifyEntryCount))
	log.Println(fmt.Sprintf("Found %d product sets.", len(pSets)))
	log.Println("Converting sets into linesheet entries...")

	var lsEntries []LSEntry

	for _, s := range pSets {
		e := s.LSEntry()
		lsEntries = append(lsEntries, e)
	}

	log.Println(fmt.Sprintf("Got %d linesheet entries.", len(lsEntries)))
	log.Println(fmt.Sprintf("Writing linesheet records to %s...", OUT_FILENAME))

	f, err := os.Create(OUT_FILENAME)
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	csvWriter.Write(LSEntryCSVFields())

	for _, e := range lsEntries {
		csvEntry := e.CSVEntry()
		record := []string(csvEntry)

		if err := csvWriter.Write(record); err != nil {
			log.Fatalln("Error writing record to CSV: ", err)
		}
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		log.Fatal(err)
	}

	log.Println("Done.")
	log.Println(fmt.Sprintf("Downloading %d image assets...", len(dlQueue)))

	DownloadQueue(dlQueue, 10)

	log.Println("Assets downloaded, or were already existing.")
}

/* Converting between string arrays from CSV reader and native data types */

type CSVEntry []string

func (csv CSVEntry) ShopifyRecord() ShopifyRecord {
	return ShopifyRecord{
		csv[0],
		csv[1],
		csv[2],

		csv[8],
		csv[10],

		csv[13],
		csv[19],

		csv[24],
		csv[25],

		csv[43],
	}
}

func (e LSEntry) CSVEntry() CSVEntry {
	csv := CSVEntry{}

	for _, fname := range LSEntryCSVFields() {
		csv = append(csv, e[fname])
	}

	return csv
}

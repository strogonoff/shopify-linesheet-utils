package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type ShopifyRecord struct {
	// Fields in the order they appear in Shopify-exported CSV

	handle string
	title  string
	body   string

	// options
	oProductType  string
	oProductColor string

	vSku   string
	vPrice string

	imageSrc string
	imagePos string

	vImage string
}

func isValidSetTitle(title string) bool {
	return strings.Contains(title, " - ") == false
}

func (sr ShopifyRecord) ProductSet() ProductSet {
	s := ProductSet{
		handle: sr.handle,
		name:   sr.title,
	}
	return s
}

func (sr ShopifyRecord) Product(wholesaleDiscountFactor float64) Product {
	var (
		retailPrice, wholesalePrice float64
	)

	retailPrice, err := strconv.ParseFloat(sr.vPrice, 32)
	if err != nil {
		log.Fatal("Error converting price to floating point", sr)
	}

	_exactPrice := retailPrice * wholesaleDiscountFactor

	wholesalePrice = RoundPlus(_exactPrice, 2)

	p := Product{
		name:           sr.oProductType,
		wholesalePrice: fmt.Sprintf("%.2f", wholesalePrice),
	}
	return p
}

func (sr ShopifyRecord) ProductVariant() ProductVariant {
	v := ProductVariant{
		sku:   sr.vSku,
		color: sr.oProductColor,
	}

	return v
}

/* Abstracting sets/products/variants */

type ProductSet struct {
	handle      string
	name        string
	products    []Product
	picturePath string
}

func (s ProductSet) maxVariantCount() int {
	maxCount := 0
	for _, p := range s.products {
		count := len(p.variants)
		if count > maxCount {
			maxCount = count
		}
	}
	return maxCount
}

type ProductVariant struct {
	sku         string
	color       string
	picturePath string
}

type Product struct {
	name           string
	variants       []ProductVariant
	wholesalePrice string
}

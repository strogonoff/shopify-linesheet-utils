package main

import (
	"fmt"
	"log"
	"strings"
	"errors"
)

// Line sheet layouts

type Layout struct {
	numColumns int
	useBigImage bool
}

func (s ProductSet) Layout() (Layout, error) {
	var nProducts = len(s.products)
	var nVariants = s.maxVariantCount()

	switch {
	case nProducts == 1 && nVariants == 1:
		return Layout{1, true}, nil
	case nProducts > 1 && nVariants == 1:
		return Layout{2, true}, nil
	case nProducts > 1 && nProducts <= 4 && nVariants <= 6:
		return Layout{1, false}, nil
	}

	return Layout{1, false}, errors.New("Product/variant count unsupported")
}

// InDesign line sheet data entry.
// Field names should correspond to data merge placeholders in corresponding
// InDesign layout.
type LSEntry map[string]string

// Logic below is specific to a very specific product line sheet.

// Returns line sheet entry created from given product set
func (s ProductSet) LSEntry() LSEntry {

	layout, err := s.Layout()

	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to determine layout for product set %s: %s",
			s.handle, err))
	}

	e := LSEntry{}

	e["setName"] = s.name

	e["@setPhoto"] = idPath(s.picturePath)

	var (
		row   int = 1
		pCell int = 1
	)

	firstVariantImage := s.products[0].variants[0].picturePath
	if layout.useBigImage == true && firstVariantImage != "" {
		e["@r1_p1_BigPhoto"] = idPath(firstVariantImage)
	}

	for pIdx, p := range s.products {
		if len(p.variants) > 6 {
			log.Fatal("Layout supports at most 6 variants per product: ", s.handle)
		}

		e[fmt.Sprintf("r%d_p%d_Product", row, pCell)] = fmt.Sprintf("%s	$%s",
			p.name,
			p.wholesalePrice)

		//e[fmt.Sprintf("r%d_p%d_Product", row, pCell)] = p.name
		//e[fmt.Sprintf("r%d_p%d_Price", row, pCell)] = fmt.Sprintf("$%s", p.wholesalePrice)

		nextProductStartsNewRow := startNewRow(pIdx, layout.numColumns)

		var vCell int
		if row != 1 && layout.numColumns == 2 && nextProductStartsNewRow {
			vCell = 2
		} else {
			vCell = 1
		}

		for _, v := range p.variants {
			var sku, color, photo string

			sku = fmt.Sprintf(" %s ", v.sku)

			if v.color != "" {
				color = v.color
			} else {
				color = "-"
			}

			if v.picturePath != "" {
				photo = idPath(v.picturePath)
			} else {
				photo = ""
			}

			e[fmt.Sprintf("r%d_pv%d_Sku", row, vCell)] = sku
			e[fmt.Sprintf("r%d_pv%d_Color", row, vCell)] = color
			e[fmt.Sprintf("@r%d_pv%d_Photo", row, vCell)] = photo

			vCell += 1
		}

		// If needed, inc row and product cell after processing product.
		// If using two columns, start rows after products with even indexes.
		// If also using big image, skip second row,
		// and start rows after products with odd indexes.
		if layout.useBigImage {
			if pIdx < 1 {
				row += 2
				pCell = 1
			} else if nextProductStartsNewRow {
				row += 1
				pCell = 1
			} else {
				pCell += 1
			}
		} else {
			if nextProductStartsNewRow {
				row += 1
				pCell = 1
			} else {
				pCell += 1
			}
		}
	}

	return e
}

func startNewRow(idx int, colNum int) bool {
	return idx%colNum == 0
}

// Constructs an InDesign link path from OS path
// This might not be, uhm, cross-platform
func idPath(path string) string {
	return fmt.Sprintf("Macintosh HD%s", strings.Replace(path, "/", ":", -1))
}

func LSEntryCSVFields() []string {
	return []string{
		"setName",
		"@setPhoto",
		"r1_p1_Product",
		"r1_p1_Price",
		"@r1_p1_BigPhoto",
		"r1_p2_Product",
		"r1_p2_Price",
		"r1_pv1_Color",
		"r1_pv1_Sku",
		"@r1_pv1_Photo",
		"r1_pv2_Color",
		"r1_pv2_Sku",
		"@r1_pv2_Photo",
		"r1_pv3_Color",
		"r1_pv3_Sku",
		"@r1_pv3_Photo",
		"r1_pv4_Color",
		"r1_pv4_Sku",
		"@r1_pv4_Photo",
		"r1_pv5_Color",
		"r1_pv5_Sku",
		"@r1_pv5_Photo",
		"r1_pv6_Color",
		"r1_pv6_Sku",
		"@r1_pv6_Photo",
		"r2_p1_Product",
		"r2_p1_Price",
		"r2_p2_Product",
		"r2_p2_Price",
		"r2_pv1_Color",
		"r2_pv1_Sku",
		"@r2_pv1_Photo",
		"r2_pv2_Color",
		"r2_pv2_Sku",
		"@r2_pv2_Photo",
		"r2_pv3_Color",
		"r2_pv3_Sku",
		"@r2_pv3_Photo",
		"r2_pv4_Color",
		"r2_pv4_Sku",
		"@r2_pv4_Photo",
		"r2_pv5_Color",
		"r2_pv5_Sku",
		"@r2_pv5_Photo",
		"r2_pv6_Color",
		"r2_pv6_Sku",
		"@r2_pv6_Photo",
		"r3_p1_Product",
		"r3_p1_Price",
		"r3_p2_Product",
		"r3_p2_Price",
		"r3_pv1_Color",
		"r3_pv1_Sku",
		"@r3_pv1_Photo",
		"r3_pv2_Color",
		"r3_pv2_Sku",
		"@r3_pv2_Photo",
		"r3_pv3_Color",
		"r3_pv3_Sku",
		"@r3_pv3_Photo",
		"r3_pv4_Color",
		"r3_pv4_Sku",
		"@r3_pv4_Photo",
		"r3_pv5_Color",
		"r3_pv5_Sku",
		"@r3_pv5_Photo",
		"r3_pv6_Color",
		"r3_pv6_Sku",
		"@r3_pv6_Photo",
		"r4_p1_Product",
		"r4_p1_Price",
		"r4_p2_Product",
		"r4_p2_Price",
		"r4_pv1_Color",
		"r4_pv1_Sku",
		"@r4_pv1_Photo",
		"r4_pv2_Color",
		"r4_pv2_Sku",
		"@r4_pv2_Photo",
		"r4_pv3_Color",
		"r4_pv3_Sku",
		"@r4_pv3_Photo",
		"r4_pv4_Color",
		"r4_pv4_Sku",
		"@r4_pv4_Photo",
		"r4_pv5_Color",
		"r4_pv5_Sku",
		"@r4_pv5_Photo",
		"r4_pv6_Color",
		"r4_pv6_Sku",
		"@r4_pv6_Photo",


		// These are not going to be used but removing these fields
		// from CSV crashes InDesign's data merge even if all
		// corresponding placeholders are removed also.
		"r5_p1_Product",
		"r5_p1_Price",
		"r5_p2_Product",
		"r5_p2_Price",
		"r5_pv1_Color",
		"r5_pv1_Sku",
		"@r5_pv1_Photo",
		"r5_pv2_Color",
		"r5_pv2_Sku",
		"@r5_pv2_Photo",
		"r5_pv3_Color",
		"r5_pv3_Sku",
		"@r5_pv3_Photo",
		"r5_pv4_Color",
		"r5_pv4_Sku",
		"@r5_pv4_Photo",
		"r5_pv5_Color",
		"r5_pv5_Sku",
		"@r5_pv5_Photo",
		"r5_pv6_Color",
		"r5_pv6_Sku",
		"@r5_pv6_Photo",
	}
}

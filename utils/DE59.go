// Package utils provides various utility functions used across the ISO 8583 library.
package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// DE59 represents the structure of ISO 8583 Field 59, which typically contains
// product-specific data with nested sub-fields.
type DE59 struct {
	// Products is a map where the key is the product code (e.g., "001")
	// and the value is another map representing sub-products.
	// The inner map's key is the sub-product code (e.g., "001") and its value is the sub-product data.
	Products map[string]map[string]string
}

// UnpackDe59 unpacks a raw string representation of ISO 8583 Field 59
// into the DE59 struct's Products map.
// It parses the product codes, quantities of sub-products, and their respective data.
// It returns an error if parsing fails due to invalid format or insufficient data.
func (de59 *DE59) UnpackDe59(de59Raw string) (err error) {
	position := 0
	products := make(map[string]map[string]string)
	var subProducts map[string]string

	if len(de59Raw) > 13 {
		for position < len(de59Raw) {

			codProduct := de59Raw[position : position+3]
			qtySubProducts, err := strconv.Atoi(de59Raw[position+3 : position+7])
			if err != nil {
				return err
			}

			position += 7

			subProducts = make(map[string]string)
			for i := 0; i < qtySubProducts; i++ {
				if len(de59Raw) >= position+3 {
					lenSubProd := de59Raw[position : position+3]
					lenSubProduct, err := strconv.Atoi(lenSubProd)
					if err != nil {
						return err
					}

					position += 3

					if len(de59Raw) >= position+3 {
						codSubProduct := de59Raw[position : position+3]
						position += 3

						if len(de59Raw) >= position+lenSubProduct {
							subProducts[codSubProduct] = de59Raw[position : position+lenSubProduct]
						} else {
							return errors.New("error when parsing field 59")
						}
						position += lenSubProduct
					} else {
						return errors.New("error when parsing field 59")
					}
				} else {
					return errors.New("error when parsing field 59")
				}
			}
			products[codProduct] = subProducts
		}
	}

	de59.Products = products

	return err
}

// RegenerateDe59 updates a specific sub-field within the DE59 structure
// and then repacks the entire DE59 field into its raw string format.
// It takes the raw DE59 string, the field path (e.g., "P01.S01"), and the new value.
// It returns the updated raw DE59 string and an error if unpacking or packing fails.
func (de59 *DE59) RegenerateDe59(de59Raw string, field string, value string) (string, error) {
	err := de59.UnpackDe59(de59Raw)
	if err != nil {
		return de59Raw, err
	}

	de59Field := strings.Split(field, ".")
	if _, ok := de59.Products[de59Field[1]]; ok {
		de59.Products[de59Field[1][1:]][de59Field[2][2:]] = value
	} else {
		product := make(map[string]string)
		product[de59Field[2][2:]] = value
		de59.Products[de59Field[1][1:]] = product
	}

	de59Raw = de59.PackDe59()

	return de59Raw, nil
}

// PackDe59 packs the DE59 struct's Products map into its raw string representation.
// It iterates through products and sub-products, formatting them according to the DE59 specification.
// It returns the raw DE59 string.
func (de59 *DE59) PackDe59() (de59Raw string) {
	for key, product := range de59.Products {
		de59Raw += fmt.Sprintf("%s%04d", key, len(product))
		for key, subProduct := range product {
			de59Raw += fmt.Sprintf("%03d%s%s", len(subProduct), key, subProduct)
		}
	}
	return de59Raw
}

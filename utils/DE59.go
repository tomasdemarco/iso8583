package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type DE59 struct {
	Products map[string]map[string]string
}

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

func (de59 *DE59) PackDe59() (de59Raw string) {
	for key, product := range de59.Products {
		de59Raw += fmt.Sprintf("%s%04d", key, len(product))
		for key, subProduct := range product {
			de59Raw += fmt.Sprintf("%03d%s%s", len(subProduct), key, subProduct)
		}
	}
	return de59Raw
}

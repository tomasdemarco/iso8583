package subfield

//type Subfields map[string]string
//
//func (s *Subfields) Pack(fieldPkg pkgField.Field) (string, error) {
//
//	subfieldsKeys := make([]string, 0, len(*s))
//	for k := range *s {
//		subfieldsKeys = append(subfieldsKeys, k)
//	}
//	sort.Strings(subfieldsKeys)
//
//	var fieldResult string
//
//	for _, subfield := range subfieldsKeys {
//		value := (*s)[subfield]
//		length := len(value)
//
//		if fieldPkg.SubfieldsData.SubFields[subfield].Encoding == encoding.Binary {
//			length = length / 2
//		}
//
//		var subfieldResult string
//
//		if fieldPkg.SubfieldsData.Format == pkgSubField.TLV {
//			subfieldResult = subfield
//		}
//
//		fieldPrefix, err := prefix.Pack(fieldPkg.SubfieldsData.SubFields[subfield].Prefix, length)
//		if err != nil {
//			return "", err
//		}
//
//		subfieldResult += fmt.Sprintf("%x", fieldPrefix)
//
//		padLeft, padRight := padding.Pack(fieldPkg.SubfieldsData.SubFields[subfield].Padding, fieldPkg.SubfieldsData.SubFields[subfield].Length, len(value))
//
//		if fieldPkg.SubfieldsData.SubFields[subfield].Encoding != encoding.None {
//
//			fieldEncode, err := encoding.Pack(fieldPkg.SubfieldsData.SubFields[subfield].Encoding, padLeft+value+padRight)
//			if err != nil {
//				return "", err
//			}
//
//			subfieldResult += fmt.Sprintf("%x", fieldEncode)
//		} else {
//			subfieldResult += padLeft + subfield + padRight
//		}
//
//		fieldResult += subfieldResult
//	}
//
//	return fieldResult, nil
//}

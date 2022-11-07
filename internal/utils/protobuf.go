package utils

import (
	"github.com/Jeffail/gabs/v2"
	"google.golang.org/genproto/protobuf/field_mask"
)

// Mask накладывает маску field_mask.FieldMask на входящий json
// и возвращает json который соответствует маске.
func Mask(fm *field_mask.FieldMask, data []byte) ([]byte, error) {
	jsonContainer, err := gabs.ParseJSON(data)
	if err != nil {
		return nil, err
	}
	baseJSON := gabs.New()
	for _, path := range fm.Paths {
		_, err := baseJSON.SetP(jsonContainer.Path(path).Data(), path)
		if err != nil {
			return nil, err
		}
	}

	return baseJSON.EncodeJSON(), nil
}

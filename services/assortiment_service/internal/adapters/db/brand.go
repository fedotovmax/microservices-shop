package db

import (
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
)

type BrandEntityFields string

func (ue BrandEntityFields) String() string {
	return string(ue)
}

const (
	BrandFieldID   BrandEntityFields = "id"
	BrandFieldSlug BrandEntityFields = "slug"
)

var ErrBrandEntityField = errors.New("the passed field does not belong to the brand entity")

func IsBrandEntityField(f BrandEntityFields) error {

	const op = "db.IsBrandEntityField"

	switch f {
	case BrandFieldID, BrandFieldSlug:
		return nil
	}

	return fmt.Errorf("%s: %w", op, ErrBrandEntityField)
}

type UpdateBrandParams struct {
	Input        *inputs.UpdateBrand
	NewSlug      *string
	SearchColumn BrandEntityFields
	SearchValue  string
}

package db

import (
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
)

type CategoryEntityFields string

func (ue CategoryEntityFields) String() string {
	return string(ue)
}

const (
	CategoryFieldID   CategoryEntityFields = "id"
	CategoryFieldSlug CategoryEntityFields = "slug"
)

var ErrCategoryEntityField = errors.New("the passed field does not belong to the category entity")

func IsCategoryEntityField(f CategoryEntityFields) error {

	const op = "db.IsCategoryEntityField"

	switch f {
	case CategoryFieldID, CategoryFieldSlug:
		return nil
	}

	return fmt.Errorf("%s: %w", op, ErrCategoryEntityField)
}

type UpdateCategoryParams struct {
	Input        *inputs.UpdateCategory
	NewSlug      *string
	SearchColumn CategoryEntityFields
	SearchValue  string
}

type FindCategoryByFieldParams struct {
	SearchColumn   CategoryEntityFields
	SearchValue    string
	Locale         string
	Recursive      bool
	WithAllLocales bool
	OnlyActive     bool
}

type FindAllCategoriesParams struct {
	Locale         string
	WithAllLocales bool
	OnlyActive     bool
}

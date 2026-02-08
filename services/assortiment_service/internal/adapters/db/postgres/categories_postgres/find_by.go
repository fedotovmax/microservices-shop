package categoriespostgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain"
	"github.com/fedotovmax/pgxtx"
)

func findByQuery(
	column db.CategoryEntityFields,
	withAllLocales bool,
	onlyActive bool,
) string {

	activeFilter := ""
	if onlyActive {
		activeFilter = "and c.is_active = true"
	}

	q := fmt.Sprintf(`
  select
  c.id as c_id, c.parent_id as c_parent_id, c.slug as c_slug,
  c.logo_url as c_logo_url, c.is_active as c_is_active,
  c.created_at as c_created_at, c.updated_at as c_updated_at, c.deleted_at as c_deleted_at,
  ct.id as ct_id, ct.title as ct_title, ct.description as ct_description,
  ct.language_code as ct_language_code
  from categories as c
  left join category_translations as ct on c.id = ct.category_id
  where c.%s = $1
	%s
	`, column, activeFilter)

	if withAllLocales {
		return q + ";"
	}

	return q + " and ct.language_code = $2;"
}

func findByQueryRecursive(column db.CategoryEntityFields, withAllLocales bool, onlyActive bool) string {

	activeFilter := ""
	if onlyActive {
		activeFilter = "and c.is_active = true"
	}

	q := fmt.Sprintf(`
	with recursive category_tree as (
    select
        c.id,
        c.parent_id,
        c.slug,
        c.logo_url,
        c.is_active,
        c.created_at,
        c.updated_at,
        c.deleted_at
    from categories as c
    where c.%s = $1
      and c.deleted_at is null
			%s
    union all

  select
        child.id,
        child.parent_id,
        child.slug,
        child.logo_url,
        child.is_active,
        child.created_at,
        child.updated_at,
        child.deleted_at
  from categories as child
    inner join category_tree as parent
        on child.parent_id = parent.id
    where child.deleted_at is null
		%s
		)
	select
    c.id           as c_id,
    c.parent_id    as c_parent_id,
    c.slug         as c_slug,
    c.logo_url     as c_logo_url,
    c.is_active    as c_is_active,
    c.created_at   as c_created_at,
    c.updated_at   as c_updated_at,
    c.deleted_at   as c_deleted_at,

    ct.id            as ct_id,
    ct.title         as ct_title,
    ct.description   as ct_description,
    ct.language_code as ct_language_code
		from category_tree as c
		left join category_translations as ct
		on ct.category_id = c.id`, column, activeFilter, activeFilter)

	if withAllLocales {
		return q + ";"
	}

	return q + " and ct.language_code = $2;"
}

func (p *postgres) FindBy(ctx context.Context, params *db.FindCategoryByFieldParams) (
	*domain.Category, error,
) {
	const op = "adapters.db.postgres.categories.FindBy"

	err := db.IsCategoryEntityField(params.SearchColumn)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tx := p.ex.ExtractTx(ctx)

	if params.Recursive {
		return p._findCategoryByRecursive(ctx, tx, params)
	}

	return p._findCategoryBy(ctx, tx, params)

}

func pushTranslationInArr(
	c *domain.Category,
	trID, trTitle, trLang *string,
	descr *string,
) {
	if trID == nil {
		return
	}

	for idx := range c.Translations {
		if c.Translations[idx].ID == *trID {
			return
		}
	}

	c.Translations = append(c.Translations, domain.Translation{
		ID:           *trID,
		Title:        *trTitle,
		LanguageCode: *trLang,
		Description:  descr,
	})
}

func (p *postgres) _findCategoryBy(
	ctx context.Context,
	tx pgxtx.PgxExecutor,
	params *db.FindCategoryByFieldParams,
) (
	*domain.Category, error,
) {
	const op = "adapters.db.postgres.categories._findCategoryBy"

	q := findByQuery(params.SearchColumn, params.WithAllLocales, params.OnlyActive)

	var args []any

	args = append(args, params.SearchValue)

	if !params.WithAllLocales {
		args = append(args, params.Locale)
	}

	rows, err := tx.Query(ctx, q, args...)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	defer rows.Close()

	c := &domain.Category{}
	found := false

	for rows.Next() {

		var trID, trTitle, trLang, descr *string

		err := rows.Scan(
			&c.ID,
			&c.ParentID,
			&c.Slug,
			&c.LogoURL,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.DeletedAt,
			&trID,
			&trTitle,
			&descr,
			&trLang,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
		}

		pushTranslationInArr(c, trID, trTitle, trLang, descr)

		if !found {
			found = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	if !found {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrNotFound, err)
	}

	return c, nil

}

func (p *postgres) _findCategoryByRecursive(
	ctx context.Context,
	tx pgxtx.PgxExecutor,
	params *db.FindCategoryByFieldParams,
) (
	*domain.Category, error,
) {
	const op = "adapters.db.postgres.categories._findCategoryByRecursive"

	q := findByQueryRecursive(params.SearchColumn, params.WithAllLocales, params.OnlyActive)

	var args []any
	args = append(args, params.SearchValue)
	if !params.WithAllLocales {
		args = append(args, params.Locale)
	}

	rows, err := tx.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}
	defer rows.Close()

	catMap := make(map[string]*domain.Category)
	childsMap := make(map[string][]*domain.Category)
	var lookingForKey string

	for rows.Next() {
		var cID string
		var parentID, logoURL, cSlug *string
		var createdAt, updatedAt time.Time
		var deletedAt *time.Time
		var isActive bool
		var trID, trTitle, trLang, descr *string

		err := rows.Scan(
			&cID,
			&parentID,
			&cSlug,
			&logoURL,
			&isActive,
			&createdAt,
			&updatedAt,
			&deletedAt,
			&trID,
			&trTitle,
			&descr,
			&trLang,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
		}

		cat, exists := catMap[cID]
		if !exists {
			cat = &domain.Category{
				ID:        cID,
				Slug:      cSlug,
				ParentID:  parentID,
				LogoURL:   logoURL,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				DeletedAt: deletedAt,
				IsActive:  isActive,
			}
			catMap[cID] = cat

			if parentID != nil {
				childsMap[*parentID] = append(childsMap[*parentID], cat)
			} else {
				if lookingForKey == "" {
					lookingForKey = cID
				}
			}
		}

		pushTranslationInArr(cat, trID, trTitle, trLang, descr)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	if len(catMap) == 0 {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrNotFound)
	}

	// Привязываем Children к каждой категории
	for id, children := range childsMap {
		if cat, ok := catMap[id]; ok {
			cat.Children = children
		}
	}

	// Возвращаем корень
	root, ok := catMap[lookingForKey]

	if !ok {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrNotFound)
	}

	b, _ := json.MarshalIndent(root, "", "  ")

	fmt.Println(string(b))

	return root, nil

}

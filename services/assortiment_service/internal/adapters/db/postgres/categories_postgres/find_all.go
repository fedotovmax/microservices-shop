package categoriespostgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain"
)

func findAllQuery(withAllLocales bool, onlyActive bool) string {

	activeFilter := ""

	if onlyActive {
		activeFilter = "and c.is_active = true"
	}

	localeFilter := ""

	if !withAllLocales {
		localeFilter = "and ct.language_code = $1"
	}

	return fmt.Sprintf(`
with recursive category_tree as (
    select
        c.id,
        c.parent_id,
        c.slug,
        c.logo_url,
        c.is_active,
        c.created_at,
        c.updated_at,
        c.deleted_at,
        0 as level
    from categories as c
    where c.deleted_at is null
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
        child.deleted_at,
        parent.level + 1
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
		on ct.category_id = c.id
    %s
    order by c.level, c.parent_id nulls first, c.created_at;`, activeFilter, activeFilter, localeFilter)
}

func (p *postgres) FindAll(
	ctx context.Context,
	params *db.FindAllCategoriesParams,
) (
	[]*domain.Category, error,
) {

	const op = "adapters.db.postgres.categories.FindAll"

	tx := p.ex.ExtractTx(ctx)

	q := findAllQuery(params.WithAllLocales, params.OnlyActive)

	var args []any
	if !params.WithAllLocales {
		args = append(args, params.Locale)
	}

	rows, err := tx.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}
	defer rows.Close()

	catMap := make(map[string]*domain.Category)
	childsMap := make(map[string][]*domain.Category)
	translationsMap := make(map[string][]domain.Translation)
	var rootIDs []string // сохраняем порядок корней

	for rows.Next() {
		var cID string
		var parentID, logoURL, cSlug *string
		var createdAt, updatedAt time.Time
		var deletedAt *time.Time
		var isActive bool
		var trID, trTitle, trLang, descr *string

		if err := rows.Scan(
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
		); err != nil {
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
				rootIDs = append(rootIDs, cID)
			}
		}

		if trID != nil {
			trs := translationsMap[cID]

			trExists := false

			for idx := range trs {
				if trs[idx].ID == *trID {
					trExists = true
				}
			}

			if !trExists {
				trs = append(trs, domain.Translation{
					ID:           *trID,
					Title:        *trTitle,
					LanguageCode: *trLang,
					Description:  descr,
				})
				translationsMap[cID] = trs
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	if len(catMap) == 0 {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	for id, children := range childsMap {
		if cat, ok := catMap[id]; ok {
			cat.Children = children
		}
	}

	for id, trs := range translationsMap {
		if cat, ok := catMap[id]; ok {
			cat.Translations = trs
		}
	}

	roots := make([]*domain.Category, 0, len(rootIDs))
	for _, id := range rootIDs {
		if cat, ok := catMap[id]; ok {
			roots = append(roots, cat)
		}
	}

	b, _ := json.MarshalIndent(roots, "", "  ")

	fmt.Println(string(b))

	return roots, nil
}

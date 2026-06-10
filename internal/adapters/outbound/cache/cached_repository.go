package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/notOliveira/onde-tem/internal/adapters/dto"
	"github.com/notOliveira/onde-tem/internal/core/domain"
	"github.com/notOliveira/onde-tem/internal/core/ports"
)

var _ ports.EstablishmentRepository = (*CachedRepository)(nil)

type CachedRepository struct {
	repo  ports.EstablishmentRepository
	cache ports.Cache
	ttl   time.Duration
}

func NewCachedRepository(repository ports.EstablishmentRepository, cache ports.Cache, ttl time.Duration) *CachedRepository {
	return &CachedRepository{
		repo:  repository,
		cache: cache,
		ttl:   ttl,
	}
}

func keyByID(id domain.EstablishmentID) string {
	return "est-" + id.String()
}

func keyBySlug(slug domain.Slug) string {
	return "est-slug:" + slug.String()
}

const listKeyPrefix = "est-list:"

func generateListCacheKey(filter *domain.EstablishmentFilter) string {
	sortedTypes := append([]string{}, filter.Types...)
	sort.Strings(sortedTypes)
	typesStr := strings.Join(sortedTypes, ",")
	raw := fmt.Sprintf("%s|%s|%d|%d", typesStr, filter.Search, filter.Limit, filter.Offset)
	hash := md5.Sum([]byte(raw))
	return listKeyPrefix + hex.EncodeToString(hash[:])
}

func (c *CachedRepository) invalidateListCache(ctx context.Context) {
	_ = c.cache.DeletePattern(ctx, listKeyPrefix+"*")
}

func (c *CachedRepository) GetByID(ctx context.Context, id domain.EstablishmentID) (*domain.Establishment, error) {
	key := keyByID(id)

	cached, err := c.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var cacheDTO dto.EstablishmentCacheDTO
		if err := json.Unmarshal([]byte(cached), &cacheDTO); err == nil {
			est, derr := cacheDTOToEstablishment(&cacheDTO)
			if derr == nil {
				return est, nil
			}
		}
	}

	est, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if cacheDTO, err := marshalEstablishmentCache(est); err == nil {
		_ = c.cache.Set(ctx, key, cacheDTO, c.ttl)
	}

	return est, nil
}

func (c *CachedRepository) GetBySlug(ctx context.Context, slug domain.Slug) (*domain.Establishment, error) {
	key := keyBySlug(slug)

	cached, err := c.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var cacheDTO dto.EstablishmentCacheDTO
		if err := json.Unmarshal([]byte(cached), &cacheDTO); err == nil {
			est, derr := cacheDTOToEstablishment(&cacheDTO)
			if derr == nil {
				return est, nil
			}
		}
	}

	est, err := c.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if cacheDTO, err := marshalEstablishmentCache(est); err == nil {
		_ = c.cache.Set(ctx, key, cacheDTO, c.ttl)
	}

	return est, nil
}

func (c *CachedRepository) List(ctx context.Context, filter *domain.EstablishmentFilter) ([]*domain.Establishment, error) {
	key := generateListCacheKey(filter)

	cached, err := c.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var cacheDTOs []dto.EstablishmentCacheDTO
		if err := json.Unmarshal([]byte(cached), &cacheDTOs); err == nil {
			establishments := make([]*domain.Establishment, 0, len(cacheDTOs))
			allValid := true
			for i := range cacheDTOs {
				est, derr := cacheDTOToEstablishment(&cacheDTOs[i])
				if derr != nil {
					allValid = false
					break
				}
				establishments = append(establishments, est)
			}
			if allValid {
				return establishments, nil
			}
		}
	}

	establishments, err := c.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	cacheDTOs := make([]dto.EstablishmentCacheDTO, 0, len(establishments))
	for _, est := range establishments {
		cacheDTOs = append(cacheDTOs, *dto.ToEstablishmentCacheDTO(est))
	}
	if bytes, err := json.Marshal(cacheDTOs); err == nil {
		_ = c.cache.Set(ctx, key, string(bytes), c.ttl)
	}

	return establishments, nil
}

func (c *CachedRepository) Create(ctx context.Context, e *domain.Establishment) error {
	if err := c.repo.Create(ctx, e); err != nil {
		return err
	}
	_ = c.cache.Delete(ctx, keyByID(e.ID()))
	_ = c.cache.Delete(ctx, keyBySlug(e.Slug()))
	c.invalidateListCache(ctx)
	return nil
}

func (c *CachedRepository) Update(ctx context.Context, e *domain.Establishment) error {
	if err := c.repo.Update(ctx, e); err != nil {
		return err
	}
	_ = c.cache.Delete(ctx, keyByID(e.ID()))
	_ = c.cache.Delete(ctx, keyBySlug(e.Slug()))
	c.invalidateListCache(ctx)
	return nil
}

func (c *CachedRepository) Delete(ctx context.Context, id domain.EstablishmentID) error {
	if err := c.repo.Delete(ctx, id); err != nil {
		return err
	}
	_ = c.cache.Delete(ctx, keyByID(id))
	c.invalidateListCache(ctx)
	return nil
}

func marshalEstablishmentCache(e *domain.Establishment) (string, error) {
	bytes, err := json.Marshal(dto.ToEstablishmentCacheDTO(e))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func cacheDTOToEstablishment(d *dto.EstablishmentCacheDTO) (*domain.Establishment, error) {
	id, err := domain.ParseEstablishmentID(d.ID)
	if err != nil {
		return nil, err
	}

	slug, err := domain.NewSlug(d.Slug)
	if err != nil {
		return nil, err
	}

	types := make([]domain.EstablishmentType, len(d.Types))
	for i, t := range d.Types {
		types[i] = domain.EstablishmentType(t)
	}

	location, err := domain.NewLocation(d.Location.Lat, d.Location.Lon)
	if err != nil {
		return nil, err
	}

	address, err := domain.NewAddress(
		d.Address.Street,
		d.Address.Number,
		d.Address.District,
		d.Address.City,
		d.Address.State,
		d.Address.Country,
		d.Address.ZipCode,
	)
	if err != nil {
		return nil, err
	}

	phones := make([]domain.Phone, 0, len(d.Phones))
	for _, p := range d.Phones {
		phone, err := domain.NewPhone(p.CountryCode, p.Number, p.Label)
		if err != nil {
			return nil, err
		}
		phones = append(phones, phone)
	}

	createdAt, _ := time.Parse(time.RFC3339, d.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, d.UpdatedAt)

	return domain.RehydrateEstablishment(
		id,
		d.Name,
		slug,
		types,
		d.Email,
		d.Website,
		phones,
		location,
		address,
		d.Timezone,
		createdAt,
		updatedAt,
	), nil
}

package domain

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var (
	ErrInvalidSlugFormat = errors.New("invalid slug format")
)

var slugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

// Slug represents a custom URL-friendly identifier for an establishment, typically derived from its name.
type Slug string

// NewSlug validates the input string and returns a Slug if it's valid, or an error if it's not.
// The slug must be non-empty, contain only lowercase letters, numbers, and hyphens, and cannot have leading or trailing hyphens.
func NewSlug(value string) (Slug, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return "", ErrInvalidSlugFormat
	}

	if !slugRegex.MatchString(value) {
		return "", ErrInvalidSlugFormat
	}

	return Slug(value), nil
}

func GenerateSlugFromName(name string) Slug {
	name = strings.ToLower(name)
	name = removeAccents(name)

	name = strings.ReplaceAll(name, " ", "-")

	var builder strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			builder.WriteRune(r)
		}
	}

	return Slug(builder.String())
}

func removeAccents(input string) string {
	t := norm.NFD.String(input)
	var b strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}
	return norm.NFC.String(b.String())
}

func (s Slug) String() string {
	return string(s)
}

func (s Slug) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

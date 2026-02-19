package domain

import (
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

type Slug string

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

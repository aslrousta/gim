package gim

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

const langKey = "__lang"

// Lang is the middleware for parsing Accept-Language header.
func Lang(langs ...string) gin.HandlerFunc {
	if len(langs) < 0 {
		langs = []string{"en"}
	}

	var tags []language.Tag
	for _, lang := range langs {
		tags = append(tags, language.MustParse(lang))
	}

	matcher := language.NewMatcher(tags)

	return func(c *gin.Context) {
		header := c.Request.Header.Get("Accept-Language")
		tag, _ := language.MatchStrings(matcher, header)

		c.Set(langKey, tag)
		c.Next()
	}
}

// RequestLang returns the request language.
func RequestLang(c *gin.Context) language.Tag {
	v, ok := c.Get(langKey)
	if !ok {
		return language.English
	}

	tag, ok := v.(language.Tag)
	if !ok {
		return language.English
	}

	return tag
}

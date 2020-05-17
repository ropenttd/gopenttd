// Package gopenttd provides utilities for communicating with OpenTTD game servers in various ways.
package gopenttd

import (
	"golang.org/x/text/language"
)

type OpenttdLanguage uint8

const (
	LanguageAny OpenttdLanguage = iota
	LanguageEnglish
	LanguageGerman
	LanguageFrench
	LanguageBrazilianPortuguese
	LanguageBulgarian
	LanguageChinese
	LanguageCzech
	LanguageDanish
	LanguageDutch
	LanguageEsperanto
	LanguageFinnish
	LanguageHungarian
	LanguageIcelandic
	LanguageItalian
	LanguageJapanese
	LanguageKorean
	LanguageLithuanian
	LanguageNorwegian
	LanguagePolish
	LanguagePortuguese
	LanguageRomanian
	LanguageRussian
	LanguageSlovak
	LanguageSloveninan
	LanguageSpanish
	LanguageSwedish
	LanguageTurkish
	LanguageUkraninan
	LanguageAfrikaans
	LanguageCroatian
	LanguageCatalan
	LanguageEstonian
	LanguageGalican
	LanguageGreek
	LanguageLatvian
	LanguageCount
)

// String returns the string representation of the Language.
func (lang OpenttdLanguage) String() string {
	names := [...]string{
		"Any",
		"English",
		"German",
		"French",
		"BrazilianPortuguese",
		"Bulgarian",
		"Chinese",
		"Czech",
		"Danish",
		"Dutch",
		"Esperanto",
		"Finnish",
		"Hungarian",
		"Icelandic",
		"Italian",
		"Japanese",
		"Korean",
		"Lithuanian",
		"Norwegian",
		"Polish",
		"Portuguese",
		"Romanian",
		"Russian",
		"Slovak",
		"Slovenian",
		"Spanish",
		"Swedish",
		"Turkish",
		"Ukrainian",
		"Afrikaans",
		"Croatian",
		"Catalan",
		"Estonian",
		"Galican",
		"Greek",
		"Latvian",
	}
	// prevent panics for out of range lookups
	if lang < LanguageAny || lang > LanguageLatvian {
		return "Unknown"
	}
	return names[lang]
}

// Tag returns a language.Tag corresponding to the given OpenttdLanguage.
// Important: If a tag is not available (for example, the Any language, or if the language doesn't exist in the text/language library),
// this will return an empty Tag{}.
func (lang OpenttdLanguage) Tag() language.Tag {
	tags := [...]language.Tag{
		language.Tag{},
		language.English,
		language.German,
		language.French,
		language.BrazilianPortuguese,
		language.Bulgarian,
		language.Chinese,
		language.Czech,
		language.Danish,
		language.Dutch,
		language.Tag{}, // Esperanto
		language.Finnish,
		language.Hungarian,
		language.Icelandic,
		language.Italian,
		language.Japanese,
		language.Korean,
		language.Lithuanian,
		language.Norwegian,
		language.Polish,
		language.Portuguese,
		language.Romanian,
		language.Russian,
		language.Slovak,
		language.Slovenian,
		language.Spanish,
		language.Swedish,
		language.Turkish,
		language.Ukrainian,
		language.Afrikaans,
		language.Croatian,
		language.Catalan,
		language.Estonian,
		language.Tag{}, // Galican
		language.Greek,
		language.Latvian,
	}
	// prevent panics for out of range lookups
	if lang < LanguageAny || lang > LanguageLatvian {
		return language.Tag{}
	}
	return tags[lang]
}

type OpenttdEnvironment uint8

const (
	EnvironmentTemperate OpenttdEnvironment = iota
	EnvironmentArctic
	EnvironmentTropic
	EnvironmentToyland
)

// String returns the string representation of the Environment.
func (env OpenttdEnvironment) String() string {
	names := [...]string{
		"Temperate",
		"Arctic",
		"Tropic",
		"Toyland",
	}
	// prevent panics for out of range lookups
	if env < EnvironmentTemperate || env > EnvironmentToyland {
		return "Unknown"
	}
	return names[env]
}

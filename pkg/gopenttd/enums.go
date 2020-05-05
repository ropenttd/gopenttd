package gopenttd

import (
	"errors"
	"golang.org/x/text/language"
)

var OpenttdLanguages = map[int]language.Tag{
	0:  {}, // Any
	1:  language.English,
	2:  language.German,
	3:  language.French,
	4:  language.BrazilianPortuguese,
	5:  language.Bulgarian,
	6:  language.Chinese,
	7:  language.Czech,
	8:  language.Danish,
	9:  language.Dutch,
	10: {}, // Esperanto
	11: language.Finnish,
	12: language.Hungarian,
	13: language.Icelandic,
	14: language.Italian,
	15: language.Japanese,
	16: language.Korean,
	17: language.Lithuanian,
	18: language.Norwegian,
	19: language.Polish,
	20: language.Portuguese,
	21: language.Romanian,
	22: language.Russian,
	23: language.Slovak,
	24: language.Slovenian,
	25: language.Spanish,
	26: language.Swedish,
	27: language.Turkish,
	28: language.Ukrainian,
	29: language.Afrikaans,
	30: language.Croatian,
	31: language.Catalan,
	32: language.Estonian,
	33: {}, // Galican
	34: language.Greek,
	35: language.Latvian,
	36: {}, // Count (afaict a marker of the last language, effectively invalid)
}

var OpenttdEnvironments = map[int]string{
	0: "Temperate",
	1: "Arctic",
	2: "Tropic",
	3: "Toyland",
}

// GetISOLanguage returns a language.Tag associated with the given language code.
// Note: The language codes "any", "esperanto", and "galican" do not presently have language tags in the Language library - an empty tag will be returned for these!
func GetISOLanguage(langcode int) (lang language.Tag, error error) {
	if val, ok := OpenttdLanguages[langcode]; ok {
		return val, nil
	} else {
		return language.Tag{}, errors.New("not a valid language identifier")
	}
}

// GetLanguage returns the english language name associated with the given language code.
func GetLanguage(langcode int) (lang string, error error) {
	if val, ok := OpenttdLanguages[langcode]; ok {
		// special cases
		switch langcode {
		case 0:
			return "Any", nil
		case 10:
			return "Esperanto", nil
		case 33:
			return "Galican", nil
		case 36:
			// invalid? whatever, return what they want
			return "Count", nil
		default:
			return val.String(), nil
		}
	} else {
		return "", errors.New("not a valid language identifier")
	}
}

// GetEnvironment returns the english environment name associated with the given environment code.
func GetEnvironment(envcode int) (environment string, error error) {
	if val, ok := OpenttdEnvironments[envcode]; ok {
		return val, nil
	} else {
		return "", errors.New("not a valid environment identifier")
	}
}

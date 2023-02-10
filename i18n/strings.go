// Package i18n provides localization and internationalization
// functionality for Ego itself.
package i18n

import (
	"fmt"
	"os"
	"strings"
)

// Language is a string that identifies the current language, such as
// "en" for English or "fr" for French. This is used as a key in the
// internal localization dictionaries.
var Language string

// Register adds additional localizations to the localization database
// at runtime. This allows an application to extend the localizations
// and still use the package i18n functions.
//
// The primary key for the map is the message code, which is the
func Register(localizations map[string]map[string]string) {
	for key, languages := range localizations {
		for language, text := range languages {
			if _, found := messages[key]; !found {
				messages[key] = map[string]string{}
			}

			messages[key][language] = text
		}
	}
}

// T converts a key into the localized string, based on the current language
// definition. If there is no localization in the given language, then "en"
// is also searched to see if it has a value. Finally, if no localization exists,
// the string is returned as-is.
//
// The second optional parameter is a map of string substitutions that should
// be done within the message text before returning the formatted value.
func T(key string, valueMap ...map[string]interface{}) string {
	// If we haven't yet figure out what language, do that now.
	if Language == "" {
		Language = os.Getenv("APP_LANG")
		if Language == "" {
			Language = os.Getenv("LANG")
		}

		if len(Language) > 2 {
			Language = Language[0:2]
		}
	}

	// Find the message using the current language
	text, ok := messages[key][Language]
	if !ok {
		text, ok = messages[key]["en"]
		if !ok {
			text = key
		}
	}

	if len(valueMap) > 0 {
		for tag, value := range valueMap[0] {
			text = strings.ReplaceAll(text, "{{"+tag+"}}", fmt.Sprintf("%v", value))
		}
	}

	return text
}

// L returns a label with the given key.
func L(key string, valueMap ...map[string]interface{}) string {
	return strings.TrimPrefix(T("label."+key, valueMap...), "label.")
}

// M returns a message with the given key.
func M(key string, valueMap ...map[string]interface{}) string {
	return T("msg."+key, valueMap...)
}

// E returns an error with the given key.
func E(key string, valueMap ...map[string]interface{}) string {
	return strings.TrimPrefix(T("error."+key, valueMap...), "error.")
}

// O returns an option description with the given key.
func O(key string, valueMap ...map[string]interface{}) string {
	return strings.TrimPrefix(T("opt."+key, valueMap...), "opt.")
}

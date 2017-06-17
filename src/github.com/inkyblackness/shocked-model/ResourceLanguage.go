package model

// ResourceLanguage specifies the language of a localized resource
type ResourceLanguage uint8

const (
	// ResourceLanguageUnspecific is for non-specific resources.
	ResourceLanguageUnspecific = ResourceLanguage(0)
	// ResourceLanguageStandard is for the default language (English).
	ResourceLanguageStandard = ResourceLanguage(1)
	// ResourceLanguageFrench is for French.
	ResourceLanguageFrench = ResourceLanguage(2)
	// ResourceLanguageGerman is for German.
	ResourceLanguageGerman = ResourceLanguage(3)
)

// ToIndex returns an integer for localized arrays.
func (lang ResourceLanguage) ToIndex() int {
	return int(lang) - 1
}

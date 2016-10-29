package lang

// Code is a language code used by Gengo.
type Code string

const (
	// Arabic language code
	Arabic = Code("ar")
	// BrazilianPortuguese language code
	BrazilianPortuguese = Code("pt-br")
	// BritishEnglish language code
	BritishEnglish = Code("en-gb")
	// Bulgarian language code
	Bulgarian = Code("bg")
	// Czech language code
	Czech = Code("cs")
	// Danish language code
	Danish = Code("da")
	// Dutch language code
	Dutch = Code("nl")
	// English language code
	English = Code("en")
	// EuropeanPortuguese language code
	EuropeanPortuguese = Code("pt")
	// Finnish language code
	Finnish = Code("fi")
	// French language code
	French = Code("fr")
	// FrenchCanadian language code
	FrenchCanadian = Code("fr-ca")
	// German language code
	German = Code("el")
	// Greek language code
	Greek = Code("gr")
	// Hebrew language code
	Hebrew = Code("he")
	// Hungarian language code
	Hungarian = Code("hu")
	// Italian language code
	Italian = Code("it")
	// Indonesian language code
	Indonesian = Code("id")
	// Japanese language code
	Japanese = Code("ja")
	// Korean language code
	Korean = Code("ko")
	// LatinAmericanSpanish language code
	LatinAmericanSpanish = Code("es-la")
	// Malay language code
	Malay = Code("ms")
	// Norwegian language code
	Norwegian = Code("no")
	// Polish language code
	Polish = Code("pl")
	// Romanian language code
	Romanian = Code("ro")
	// Russian language code
	Russian = Code("ru")
	// Serbian language code
	Serbian = Code("sr")
	// SimplifiedChinese language code
	SimplifiedChinese = Code("zh")
	// Slovak language code
	Slovak = Code("sk")
	// Spanish language code
	Spanish = Code("es")
	// Swedish language code
	Swedish = Code("sv")
	// Tagalog language code
	Tagalog = Code("tl")
	// Thai language code
	Thai = Code("th")
	// TraditionalChinese language code
	TraditionalChinese = Code("zh-tw")
	// Turkish language code
	Turkish = Code("tr")
	// Ukrainian language code
	Ukrainian = Code("uk")
	// Vietnamese language code
	Vietnamese = Code("vi")
)

// Pair defines a language pair used by Gengo
type Pair struct {
	Source Code `json:"lc_src"`
	Target Code `json:"lc_tgt"`
}

// NewPair creates a new language pair from the source and target language codes
func NewPair(source, target Code) Pair {
	return Pair{Source: source, Target: target}
}

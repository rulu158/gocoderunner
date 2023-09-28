package languages

type Language int64

const (
	Go Language = iota
)

var LanguageExtensions = map[Language]string{Go: "go"}

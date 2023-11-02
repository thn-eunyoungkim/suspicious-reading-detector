package main

const (
	spanish = "Spanish"
	korean  = "Korean"

	englishHelloPrefix = "Hello, "
	spanishHelloPrefix = "Hola, "
	koreanHelloPrefix  = "안녕, "
)

func Hello(name string, language string) string {
	if name == "" {
		name = "no name"
	}

	return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
	switch language {
	case spanish:
		prefix = spanishHelloPrefix
	case korean:
		prefix = koreanHelloPrefix
	default:
		prefix = englishHelloPrefix
	}
	return
}

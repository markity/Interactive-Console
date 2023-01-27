package interactive

type Config struct {
	Prompt               rune
	BlockInputAfterRun   bool
	BlockInputAfterEnter bool
	TraceAfterRun        bool
}

func GetDefaultConfig() Config {
	return Config{
		Prompt:               '>',
		BlockInputAfterRun:   false,
		BlockInputAfterEnter: false,
		TraceAfterRun:        false,
	}
}

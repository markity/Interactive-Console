package interactive

type Config struct {
	Prompt                 rune
	BlockInputAfterRun     bool
	BlockInputAfterEnter   bool
	TraceAfterRun          bool
	SpecialEventHandleMask int64
}

func GetDefaultConfig() Config {
	return Config{
		Prompt:                 '>',
		BlockInputAfterRun:     false,
		BlockInputAfterEnter:   false,
		TraceAfterRun:          false,
		SpecialEventHandleMask: 0,
	}
}

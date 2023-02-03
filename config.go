package interactive

type Config struct {
	// 命令提示符, 允许使用一个utf8符号
	Prompt rune

	// 命令提示符的颜色
	PromptStyle StyleAttr

	// 是否在运行后阻塞用户输入, 只有使用Win.SetBlockInput(false)后用户才能输入
	BlockInputAfterRun bool

	// 是否在用户输入内容按回车之后阻塞输入, 这样每次输入命令后由命令接收方决定是否允许用户进一步输入
	BlockInputAfterEnter bool

	// 是否在运行后追踪最新的信息
	TraceAfterRun bool

	// 用来说明需要接收哪些事件
	EventHandleMask int64
}

func GetDefaultConfig() Config {
	return Config{
		Prompt:               '>',
		PromptStyle:          GetDefaultSytleAttr(),
		BlockInputAfterRun:   false,
		BlockInputAfterEnter: false,
		TraceAfterRun:        false,
		EventHandleMask:      0,
	}
}

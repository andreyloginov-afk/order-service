package section

type Monitor struct {
	LogLevel    string `split_words:"true" default:"info"`
	Environment string `split_words:"true" default:"development"`
}

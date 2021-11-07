package vendors

type Client interface {
	TranslateText(text string, sl string, tl string) ([]string, error)
}

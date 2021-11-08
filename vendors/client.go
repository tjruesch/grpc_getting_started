package vendors

import "io"

type Client interface {
	TranslateText(text []string, sl string, tl string) ([]string, error)
	TranslateFile(file io.ReadCloser, sl string, tl string) (io.ReadCloser, error)
}

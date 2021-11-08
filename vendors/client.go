package vendors

import (
	"os"
)

type Client interface {
	// TranslateText takes the string of input texts and returns a list of
	// translated output texts
	TranslateText(text []string, sl string, tl string) ([]string, error)
	// TranslateFile takes an input file and returns the translated version
	// (NOT IMPLEMENTED)
	TranslateFile(file *os.File, sl string, tl string) (*os.File, error)
}

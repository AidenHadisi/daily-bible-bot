package bible

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
)

var verseRegex = regexp.MustCompile(`(?P<book>(\d\s?)?\w+)\s(?P<chapter>\d+)\s*:\s*(?P<begin>\d+)(\s?-\s?(?P<end>\d+))?`)
var imgRegex = regexp.MustCompile(`(\w+)=("(.*?)"|\S+)`)

//Parser defines parsed text results.
type Parser struct {
	Book    string
	Chapter string
	Start   string
	End     string
	Img     string
	Size    int
}

func NewParser() *Parser {
	return &Parser{
		Size: 40,
	}
}

//HasImage returns if the text included an image request
func (p *Parser) HasImage() bool {
	return p.Img != ""
}

//GetPath constructs a API path from the ParsedText
func (p *Parser) GetPath() string {
	verse := fmt.Sprintf("%s %s:%s", p.Book, p.Chapter, p.Start)

	if p.End != "" {
		verse = fmt.Sprintf("%s-%s", verse, p.End)
	}

	return verse
}

//Parse parses a tweet text.
func (v *Parser) Parse(text string) error {
	result := verseRegex.FindStringSubmatch(text)
	if result == nil {
		return errors.New("incorrect text provided")
	}

	v.Book = strings.Title(strings.ToLower(result[1]))
	v.Chapter = result[3]
	v.Start = result[4]
	v.End = result[6]

	customParams := make(map[string]string)
	params := imgRegex.FindAllStringSubmatch(text, -1)
	for _, param := range params {
		if param[3] == "" {
			customParams[param[1]] = param[2]
		} else {
			customParams[param[1]] = param[3]
		}
	}

	err := mapstructure.WeakDecode(&customParams, &v)
	if err != nil {
		return err
	}
	return nil
}

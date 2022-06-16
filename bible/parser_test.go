package bible

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		hasImage bool
		wantErr  bool
		expected *Parser
	}{
		{"invalid text", "hello world", false, true, &Parser{Size: 40}},
		{"valid text", "hello mark 1:1-2 world", false, false, &Parser{Book: "Mark", Chapter: "1", Start: "1", End: "2", Size: 40}},
		{"with image", "hello mark 1:1-2 world img=https://imgur.com/test.jpg", true, false, &Parser{Book: "Mark", Chapter: "1", Start: "1", End: "2", Img: "https://imgur.com/test.jpg", Size: 40}},
		{"with image and size", "hello mark 1:1-2 world img=https://imgur.com/test.jpg size=10", true, false, &Parser{Book: "Mark", Chapter: "1", Start: "1", End: "2", Img: "https://imgur.com/test.jpg", Size: 10}},
	}

	for _, test := range tests {
		parser := NewParser()
		err := parser.Parse(test.text)
		assert.Equal(t, test.wantErr, (err != nil))
		assert.Equal(t, test.expected, parser)
		assert.Equal(t, test.hasImage, parser.HasImage())
	}
}

func TestPath(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		data     *Parser
	}{
		{"with end", "Mark 1:1-2", &Parser{Book: "Mark", Chapter: "1", Start: "1", End: "2", Size: 40}},
		{"without end", "Mark 1:1", &Parser{Book: "Mark", Chapter: "1", Start: "1", Size: 40}},
	}

	for _, test := range tests {
		parser := test.data
		assert.Equal(t, test.expected, parser.GetPath())
	}
}

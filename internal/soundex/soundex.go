package soundex

import (
	"unicode"
)

type Entry struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Entries []Entry

// SoundexEncoder Soundex encoder
type SoundexEncoder struct{}

// NewSoundexEncoder create new encoder instance
func NewSoundexEncoder() *SoundexEncoder {
	return &SoundexEncoder{}
}

// Encode thansform the input string to its phonetic representation using Soundex algorithm.
func (se *SoundexEncoder) Encode(input string) string {
	if len(input) == 0 {
		return input
	}

	runes := clean(input)

	result := []rune{'0', '0', '0', '0'}
	if len(runes) == 0 {
		return string(runes)
	}
	lastChar, ok := mapper[runes[0]]
	if ok {
		result[0] = runes[0]
	}

	cnt := 1
	for i := 1; i < len(runes) && cnt < len(result); i++ {
		char := runes[i]

		// ignore this characters if they are not at first position
		if char == 'H' || char == 'W' || char == 'Y' {
			continue
		}

		char = mapper[char]
		if char != '0' && char != lastChar {
			result[cnt] = char
			cnt++
		}

		lastChar = char
	}

	return string(result)
}

func clean(input string) []rune {
	cleaned := []rune{}
	if len(input) == 0 {
		return cleaned
	}

	for _, r := range input {
		if !unicode.IsLetter(r) {
			continue
		}
		cleaned = append(cleaned, unicode.ToUpper(r))
	}

	return cleaned
}

var mapper = map[rune]rune{
	'A': '0',
	'B': '1',
	'C': '2',
	'D': '3',
	'E': '0',
	'F': '1',
	'G': '2',
	'H': '0',
	'I': '0',
	'J': '2',
	'K': '2',
	'L': '4',
	'M': '5',
	'N': '5',
	'O': '0',
	'P': '1',
	'Q': '2',
	'R': '6',
	'S': '2',
	'T': '3',
	'U': '0',
	'V': '1',
	'W': '0',
	'X': '2',
	'Y': '0',
	'Z': '2',
}

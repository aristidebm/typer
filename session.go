package typer

import (
	"encoding/json"
	"io"
	"strings"
	"time"
	"unicode"
)

type KeyEvent struct {
	Key     rune      `json:"key"`
	Date    time.Time `json:"date"`
	Correct bool      `json:"correct"`
}

type Word struct {
	Text     string     `json:"text"`
	Progress string     `json:"progress"`
	Events   []KeyEvent `json:"events"`
}

func (w *Word) IsMissed() bool {
	for _, evt := range w.Events {
		if !evt.Correct {
			return true
		}
	}
	return false
}

type Result struct {
	Missing []string `json:"missing"`
}

type Session struct {
	Words       []Word `json:"words"`
	CurrentWord int    `json:"currentWord"`
	Result      Result `json:"results"`
}

func DecodeSession(r io.Reader) (*Session, error) {
	var s Session
	dec := json.NewDecoder(r)
	if err := dec.Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func NewSession(r io.Reader) (*Session, error) {
	words, err := wordsFrom(r)
	if err != nil {
		return nil, err
	}
	return &Session{Words: words}, nil
}

func (s *Session) HandleKey(key rune) {
	word := &s.Words[s.CurrentWord]
	if len(word.Progress) >= len(word.Text) && unicode.IsSpace(key) {
		// We have reached the end of the word go to the next word
		s.nextWord()
		return
	}
	word.Progress += string(key)
	word.Events = append(word.Events, KeyEvent{
		Key:     key,
		Date:    time.Now().UTC(),
		Correct: strings.HasPrefix(word.Text, word.Progress),
	})
}

func (s *Session) DeleteWord() {
	word := &s.Words[s.CurrentWord]
	if word.Progress != "" {
		word.Progress = ""
		return
	}
	s.prevWord()
	word = &s.Words[s.CurrentWord]
	word.Progress = ""
}

func (s *Session) ComputeResult() {
	for _, word := range s.Words {
		s.Result.Missing = append(s.Result.Missing, word.Text)
	}
}

func (s *Session) Encode(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(s)
}

func (s *Session) nextWord() {
	if s.CurrentWord < len(s.Words)-1 {
		s.CurrentWord += 1
	}
}

func (s *Session) prevWord() {
	if s.CurrentWord > 0 {
		s.CurrentWord -= 1
	}
}

func wordsFrom(r io.Reader) ([]Word, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	text := string(data)

	var words []Word
	for word := range strings.FieldsSeq(text) {
		words = append(words, Word{Text: word})
	}
	return words, nil
}

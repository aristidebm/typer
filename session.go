package typer

import (
	"encoding/json"
	_ "fmt"
	"io"
	_ "regexp"
	"slices"
	"strings"
	"time"
	"unicode"
)

type KeyEvent struct {
	Key         rune      `json:"key"`
	Date        time.Time `json:"date"`
	ExpectedKey rune      `json:"expecteKey"`
}

func (ke KeyEvent) IsMissed() bool {
	return ke.ExpectedKey != ke.Key
}

func (ke *KeyEvent) MarshalJSON() ([]byte, error) {
	aux := struct {
		Key         string    `json:"key"`
		Date        time.Time `json:"date"`
		ExpectedKey string    `json:"expecteKey"`
	}{
		Key:         string(ke.Key),
		Date:        ke.Date,
		ExpectedKey: string(ke.ExpectedKey),
	}
	return json.Marshal(&aux)
}

func (ke *KeyEvent) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Key         string    `json:"key"`
		Date        time.Time `json:"date"`
		ExpectedKey string    `json:"expecteKey"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	ke.Key = []rune(aux.Key)[0]
	ke.ExpectedKey = []rune(aux.ExpectedKey)[0]
	ke.Date = aux.Date
	return nil
}

type Word struct {
	Text     []rune     `json:"text"`
	Progress []rune     `json:"progress"`
	Events   []KeyEvent `json:"events"`
}

func (w *Word) MarshalJSON() ([]byte, error) {
	aux := struct {
		Text     string     `json:"text"`
		Progress string     `json:"progress"`
		Events   []KeyEvent `json:"events"`
	}{
		Text:     string(w.Text),
		Progress: string(w.Progress),
		Events:   w.Events,
	}
	return json.Marshal(&aux)
}

func (w *Word) UnmarshalJSON(data []byte) error {
	aux := struct {
		Text     string     `json:"text"`
		Progress string     `json:"progress"`
		Events   []KeyEvent `json:"events"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	w.Text = []rune(aux.Text)
	w.Progress = []rune(aux.Progress)
	w.Events = aux.Events
	return nil
}

func (w *Word) IsMissed() bool {
	for _, evt := range w.Events {
		if evt.IsMissed() {
			return true
		}
	}
	return false
}

type Result struct {
	Missing []string `json:"missing"`
}

type Session struct {
	Words       []Word    `json:"words"`
	CurrentWord int       `json:"currentWord"`
	Result      Result    `json:"results"`
	StartedAt   time.Time `json:"startedAt"`
	EndedAt     time.Time `json:"endedAt"`
	completed   bool
}

func DecodeSession(r io.Reader) (*Session, error) {
	var s Session
	dec := json.NewDecoder(r)
	if err := dec.Decode(&s); err != nil {
		return nil, err
	}
	if s.StartedAt.IsZero() {
		s.StartedAt = time.Now().UTC()
	}
	return &s, nil
}

func NewSession(r io.Reader) (*Session, error) {
	words, err := wordsFrom(r)
	if err != nil {
		return nil, err
	}
	return &Session{Words: words, StartedAt: time.Now().UTC()}, nil
}

func (s *Session) HandleKey(key rune) {
	word := &s.Words[s.CurrentWord]
	if len(word.Progress) >= len(word.Text) && unicode.IsSpace(key) {
		// We have reached the end of the word go to the next word
		s.nextWord()
		return
	}
	word.Progress = append(word.Progress, key)
	expectedKeyIndex := min(len(word.Text), len(word.Progress)) - 1
	expectedKey := word.Text[expectedKeyIndex]
	if len(word.Progress) > len(word.Text) && !unicode.IsSpace(key) {
		expectedKey = []rune(" ")[0]
	}
	word.Events = append(word.Events, KeyEvent{
		Key:         key,
		Date:        time.Now().UTC(),
		ExpectedKey: expectedKey,
	})
	if s.CurrentWord >= len(s.Words)-1 && len(word.Progress) >= len(word.Text) {
		s.completeSession()
		return
	}
}

func (s *Session) DeleteWord() {
	word := &s.Words[s.CurrentWord]
	if len(word.Progress) != 0 {
		word.Progress = nil
		return
	}
	s.prevWord()
	word = &s.Words[s.CurrentWord]
	word.Progress = nil
}

func (s *Session) DeleteChar() {
	word := &s.Words[s.CurrentWord]
	if len(word.Progress) == 0 {
		s.prevWord()
		return
	}
	begin := max(len(word.Progress)-1, 0)
	// remove the last element
	word.Progress = slices.Delete(word.Progress, begin, begin+1)
}

func (s *Session) ComputeResult() {
	for _, word := range s.Words {
		if word.IsMissed() {
			s.Result.Missing = append(s.Result.Missing, string(word.Text))
		}
	}
}

func (s *Session) Encode(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(s)
}

func (s Session) IsCompleted() bool {
	return s.completed
}

func (s Session) Duration() time.Duration {
	EndedAt := s.EndedAt
	if EndedAt.IsZero() {
		EndedAt = time.Now().UTC()
	}
	return EndedAt.Sub(s.StartedAt)
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

func (s *Session) completeSession() {
	s.completed = true
	s.EndedAt = time.Now()
}

// var _ = regexp.MustCompile(`[\p{L}\p{N}]+`)

func wordsFrom(r io.Reader) ([]Word, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	text := string(data)
	var words []Word
	// sequences := strings.FieldsFuncSeq(text, func(r rune) bool {
	// 		// Split on anything that is NOT a letter or digit
	// 		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	// })
	// for _, word := range wordRe.FindAllString(text, -1) {
	for word := range strings.FieldsSeq(text) {
		words = append(words, Word{Text: []rune(word)})
	}
	return words, nil
}

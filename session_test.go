package typer

import (
	"io"
	"strings"
	"testing"
)

func TestHandleKey(t *testing.T) {
	testCases := []struct {
		Name       string
		Text       io.Reader
		Keypresses []rune
		Session    Session
	}{
		{
			"CorrectKeypressesOnSingleWord",
			strings.NewReader("hello"),
			[]rune{'h', 'e', 'l', 'l', 'o'},
			Session{
				CurrentWord: 0,
				Words: []Word{
					{
						Text:     "hello",
						Progress: "hello",
						Events: []KeyEvent{
							{
								Key:     'h',
								Correct: true,
							},
							{
								Key:     'e',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'o',
								Correct: true,
							},
						},
					},
				},
			},
		},
		{
			"CorrectKeypressesOnMultipleWords",
			strings.NewReader("Hello World"),
			[]rune{'H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r', 'l', 'd'},
			Session{
				CurrentWord: 1,
				Words: []Word{
					{
						Text:     "Hello",
						Progress: "Hello",
						Events: []KeyEvent{
							{
								Key:     'H',
								Correct: true,
							},
							{
								Key:     'e',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'o',
								Correct: true,
							},
						},
					},
					{
						Text:     "World",
						Progress: "World",
						Events: []KeyEvent{
							{
								Key:     'W',
								Correct: true,
							},
							{
								Key:     'o',
								Correct: true,
							},
							{
								Key:     'r',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'd',
								Correct: true,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			// instanciate session
			s, err := NewSession(tt.Text)
			if err != nil {
				t.Errorf("failed creating session: %s", err)
			}

			for _, kp := range tt.Keypresses {
				s.HandleKey(kp)
			}
			s1 := tt.Session
			CompareSessions(t, *s, s1)
		})
	}
}

func TestDeleteWord(t *testing.T) {
	testCases := []struct {
		Name       string
		Text       io.Reader
		Keypresses []rune
		Session    Session
	}{
		{
			"DeleteCurrentWorldInTheMiddleOfTyping",
			strings.NewReader("Hello World"),
			[]rune{'H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r'},
			Session{
				CurrentWord: 1,
				Words: []Word{
					{
						Text:     "Hello",
						Progress: "Hello",
						Events: []KeyEvent{
							{
								Key:     'H',
								Correct: true,
							},
							{
								Key:     'e',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'o',
								Correct: true,
							},
						},
					},
					{
						Text:     "World",
						Progress: "",
						Events: []KeyEvent{
							{
								Key:     'W',
								Correct: true,
							},
							{
								Key:     'o',
								Correct: true,
							},
							{
								Key:     'r',
								Correct: true,
							},
						},
					},
				},
			},
		},
		{
			"DeleteCurrentWorldAtTheEndOfTyping",
			strings.NewReader("Hello World"),
			[]rune{'H', 'e', 'l', 'l', 'o', ' '},
			Session{
				CurrentWord: 0,
				Words: []Word{
					{
						Text:     "Hello",
						Progress: "",
						Events: []KeyEvent{
							{
								Key:     'H',
								Correct: true,
							},
							{
								Key:     'e',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'o',
								Correct: true,
							},
						},
					},
					{
						Text:     "World",
						Progress: "",
						Events:   []KeyEvent{},
					},
				},
			},
		},
		{
			"DeleteWorldTwoTimes",
			strings.NewReader("Hello World"),
			[]rune{'H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r', 'l', 'd'},
			Session{
				CurrentWord: 1,
				Words: []Word{
					{
						Text:     "Hello",
						Progress: "",
						Events: []KeyEvent{
							{
								Key:     'H',
								Correct: true,
							},
							{
								Key:     'e',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'o',
								Correct: true,
							},
						},
					},
					{
						Text:     "World",
						Progress: "",
						Events: []KeyEvent{
							{
								Key:     'W',
								Correct: true,
							},
							{
								Key:     'o',
								Correct: true,
							},
							{
								Key:     'r',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'd',
								Correct: true,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {

			// instanciate session
			s, err := NewSession(tt.Text)
			if err != nil {
				t.Errorf("failed creating session: %s", err)
			}

			for _, kp := range tt.Keypresses {
				s.HandleKey(kp)
			}

			// special case handling
			numberOfDeletion := 1
			if strings.Contains(tt.Name, "Two") {
				numberOfDeletion += 1
				tt.Session.CurrentWord -= 1
			}

			for _ = range numberOfDeletion {
				s.DeleteWord()
			}

			s1 := tt.Session
			CompareSessions(t, *s, s1)
		})
	}
}

func TestMissingKeys(t *testing.T) {
	testCases := []struct {
		Name       string
		Text       io.Reader
		Keypresses []rune
		Session    Session
	}{
		{
			"MissingKeysInTheMiddleOfTyping",
			strings.NewReader("hello"),
			[]rune{'h', 'e', 'x', 'l', 'o'},
			Session{
				CurrentWord: 0,
				Words: []Word{
					{
						Text:     "hello",
						Progress: "hexlo",
						Events: []KeyEvent{
							{
								Key:     'h',
								Correct: true,
							},
							{
								Key:     'e',
								Correct: true,
							},
							{
								Key:     'x',
								Correct: false,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'o',
								Correct: true,
							},
						},
					},
				},
			},
		},
		{
			"MissingKeysInTheBeginningOfTyping",
			strings.NewReader("hello"),
			[]rune{'H', 'e', 'l', 'l', 'o'},
			Session{
				CurrentWord: 0,
				Words: []Word{
					{
						Text:     "hello",
						Progress: "Hello",
						Events: []KeyEvent{
							{
								Key:     'H',
								Correct: false,
							},
							{
								Key:     'e',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'o',
								Correct: true,
							},
						},
					},
				},
			},
		},
		{
			"MissingKeysAtTheEndOfTyping",
			strings.NewReader("hello"),
			[]rune{'h', 'e', 'l', 'l', 'p'},
			Session{
				CurrentWord: 0,
				Words: []Word{
					{
						Text:     "hello",
						Progress: "hellp",
						Events: []KeyEvent{
							{
								Key:     'h',
								Correct: true,
							},
							{
								Key:     'e',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'l',
								Correct: true,
							},
							{
								Key:     'p',
								Correct: false,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			// instanciate session
			s, err := NewSession(tt.Text)
			if err != nil {
				t.Errorf("failed creating session: %s", err)
			}
			for _, kp := range tt.Keypresses {
				s.HandleKey(kp)
			}
			s1 := tt.Session
			CompareSessions(t, *s, s1)
		})
	}
}

func CompareSessions(t *testing.T, actual Session, expected Session) {
	for i := range actual.Words {

		// check current word
		if actual.CurrentWord != expected.CurrentWord {
			t.Errorf("Current Word expected: %d, actual: %d", expected.CurrentWord, actual.CurrentWord)
		}

		// check Text
		if actual.Words[i].Text != expected.Words[i].Text {
			t.Errorf("Text expected: %s, actual: %s", expected.Words[i].Text, actual.Words[i].Text)
		}

		// check Progress
		if actual.Words[i].Progress != expected.Words[i].Progress {
			t.Errorf("Progress expected: %s, actual: %s", expected.Words[i].Progress, actual.Words[i].Progress)
		}

		// check Events
		for j := range actual.Words[i].Events {
			if actual.Words[i].Events[j].Key != expected.Words[i].Events[j].Key {
				t.Errorf("KeyEvent at position %d on word %d expected: %d, actual: %d", j, i, expected.Words[i].Events[j].Key, actual.Words[i].Events[j].Key)
			}
			if actual.Words[i].Events[j].Correct != expected.Words[i].Events[j].Correct {
				t.Errorf("KeyEvent at position %d on word %d expected: %v, actual: %v", j, i, expected.Words[i].Events[j].Correct, actual.Words[i].Events[j].Correct)
			}
		}
	}
}

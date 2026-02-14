package typer

import (
	_ "fmt"
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
						Text:     []rune("hello"),
						Progress: []rune("hello"),
						Events: []KeyEvent{
							{
								Key:         'h',
								ExpectedKey: 'h',
							},
							{
								Key:         'e',
								ExpectedKey: 'e',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
						},
					},
				},
				Result: Result{Missing: []string{}},
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
						Text:     []rune("Hello"),
						Progress: []rune("Hello"),
						Events: []KeyEvent{
							{
								Key:         'H',
								ExpectedKey: 'H',
							},
							{
								Key:         'e',
								ExpectedKey: 'e',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
						},
					},
					{
						Text:     []rune("World"),
						Progress: []rune("World"),
						Events: []KeyEvent{
							{
								Key:         'W',
								ExpectedKey: 'W',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
							{
								Key:         'r',
								ExpectedKey: 'r',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'd',
								ExpectedKey: 'd',
							},
						},
					},
				},
				Result: Result{Missing: []string{}},
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
			s.ComputeResult()
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
						Text:     []rune("Hello"),
						Progress: []rune("Hello"),
						Events: []KeyEvent{
							{
								Key:         'H',
								ExpectedKey: 'H',
							},
							{
								Key:         'e',
								ExpectedKey: 'e',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
						},
					},
					{
						Text:     []rune("World"),
						Progress: []rune(""),
						Events: []KeyEvent{
							{
								Key:         'W',
								ExpectedKey: 'W',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
							{
								Key:         'r',
								ExpectedKey: 'r',
							},
						},
					},
				},
				Result: Result{Missing: []string{}},
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
						Text:     []rune("Hello"),
						Progress: []rune(""),
						Events: []KeyEvent{
							{
								Key:         'H',
								ExpectedKey: 'H',
							},
							{
								Key:         'e',
								ExpectedKey: 'e',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
						},
					},
					{
						Text:     []rune("World"),
						Progress: []rune(""),
						Events:   []KeyEvent{},
					},
				},
				Result: Result{Missing: []string{}},
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
						Text:     []rune("Hello"),
						Progress: []rune(""),
						Events: []KeyEvent{
							{
								Key:         'H',
								ExpectedKey: 'H',
							},
							{
								Key:         'e',
								ExpectedKey: 'e',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
						},
					},
					{
						Text:     []rune("World"),
						Progress: []rune(""),
						Events: []KeyEvent{
							{
								Key:         'W',
								ExpectedKey: 'W',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
							{
								Key:         'r',
								ExpectedKey: 'r',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'd',
								ExpectedKey: 'd',
							},
						},
					},
				},
				Result: Result{Missing: []string{}},
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
			s.ComputeResult()
			s1 := tt.Session
			CompareSessions(t, *s, s1)
		})
	}
}

func TestDeleteChar(t *testing.T) {
	testCases := []struct {
		Name       string
		Text       io.Reader
		Keypresses []rune
		Session    Session
	}{
		{
			"DeleteCharWhenAtTheBeginningOfTyping",
			strings.NewReader("Hello"),
			[]rune{},
			Session{
				CurrentWord: 0,
				Words: []Word{
					{
						Text:     []rune("Hello"),
						Progress: []rune(""),
						Events:   []KeyEvent{},
					},
				},
				Result: Result{Missing: []string{}},
			},
		},
		{
			"DeleteCharWhenAtTheBeginningOfNextWord",
			strings.NewReader("Hello World"),
			[]rune{'H', 'e', 'l', 'l', 'o', ' ', 'W'},
			Session{
				CurrentWord: 1,
				Words: []Word{
					{
						Text:     []rune("Hello"),
						Progress: []rune("Hello"),
						Events: []KeyEvent{
							{
								Key:         'H',
								ExpectedKey: 'H',
							},
							{
								Key:         'e',
								ExpectedKey: 'e',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
						},
					},
					{
						Text:     []rune("World"),
						Progress: []rune(""),
						Events: []KeyEvent{
							{
								Key:         'W',
								ExpectedKey: 'W',
							},
						},
					},
				},
				Result: Result{Missing: []string{}},
			},
		},
		{
			"DeleteCharTwoTimesWhenAtTheBeginningOfNextWord",
			strings.NewReader("Hello World"),
			[]rune{'H', 'e', 'l', 'l', 'o', ' ', 'W'},
			Session{
				CurrentWord: 0,
				Words: []Word{
					{
						Text:     []rune("Hello"),
						Progress: []rune("Hello"),
						Events: []KeyEvent{
							{
								Key:         'H',
								ExpectedKey: 'H',
							},
							{
								Key:         'e',
								ExpectedKey: 'e',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
						},
					},
					{
						Text:     []rune("World"),
						Progress: []rune(""),
						Events: []KeyEvent{
							{
								Key:         'W',
								ExpectedKey: 'W',
							},
						},
					},
				},
				Result: Result{Missing: []string{}},
			},
		},
		{
			"DeleteCharWhenAtTheEndOfCurrentWord",
			strings.NewReader("Hello World"),
			[]rune{'H', 'e', 'l', 'l', 'o', ' '},
			Session{
				CurrentWord: 0,
				Words: []Word{
					{
						Text:     []rune("Hello"),
						Progress: []rune("Hello"),
						Events: []KeyEvent{
							{
								Key:         'H',
								ExpectedKey: 'H',
							},
							{
								Key:         'e',
								ExpectedKey: 'e',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
						},
					},
					{
						Text:     []rune("World"),
						Progress: []rune(""),
						Events:   []KeyEvent{},
					},
				},
				Result: Result{Missing: []string{}},
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
			}
			for _ = range numberOfDeletion {
				s.DeleteChar()
			}
			s.ComputeResult()
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
						Text:     []rune("hello"),
						Progress: []rune("hexlo"),
						Events: []KeyEvent{
							{
								Key:         'h',
								ExpectedKey: 'h',
							},
							{
								Key:         'e',
								ExpectedKey: 'e',
							},
							{
								Key:         'x',
								ExpectedKey: 'l',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
						},
					},
				},
				Result: Result{Missing: []string{"hello"}},
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
						Text:     []rune("hello"),
						Progress: []rune("Hello"),
						Events: []KeyEvent{
							{
								Key:         'H',
								ExpectedKey: 'h',
							},
							{
								Key:         'e',
								ExpectedKey: 'e',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
						},
					},
				},
				Result: Result{Missing: []string{"hello"}},
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
						Text:     []rune("hello"),
						Progress: []rune("hellp"),
						Events: []KeyEvent{
							{
								Key:         'h',
								ExpectedKey: 'h',
							},
							{
								Key:         'e',
								ExpectedKey: 'e',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'p',
								ExpectedKey: 'o',
							},
						},
					},
				},
				Result: Result{Missing: []string{"hello"}},
			},
		},
		{
			"MissingKeysBecauseOfSpace",
			strings.NewReader("hello"),
			// notice that the user does not hit the space
			[]rune{'h', 'e', 'l', 'l', 'o', 'w', 'o', 'r', 'l', 'd'},
			Session{
				CurrentWord: 0,
				Words: []Word{
					{
						Text:     []rune("hello"),
						Progress: []rune("helloworld"),
						Events: []KeyEvent{
							{
								Key:         'h',
								ExpectedKey: 'h',
							},
							{
								Key:         'e',
								ExpectedKey: 'e',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'l',
								ExpectedKey: 'l',
							},
							{
								Key:         'o',
								ExpectedKey: 'o',
							},
							{
								Key:         'w',
								ExpectedKey: ' ',
							},
							{
								Key:         'o',
								ExpectedKey: ' ',
							},
							{
								Key:         'r',
								ExpectedKey: ' ',
							},
							{
								Key:         'l',
								ExpectedKey: ' ',
							},
							{
								Key:         'd',
								ExpectedKey: ' ',
							},
						},
					},
				},
				Result: Result{Missing: []string{"hello"}},
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
			s.ComputeResult()
			s1 := tt.Session
			CompareSessions(t, *s, s1)
		})
	}
}

func CompareSessions(t *testing.T, actual Session, expected Session) {
	for i := range actual.Words {

		// // check current word
		if actual.CurrentWord != expected.CurrentWord {
			t.Errorf("Current Word expected: %d, actual: %d", expected.CurrentWord, actual.CurrentWord)
		}

		if string(actual.Words[i].Text) != string(expected.Words[i].Text) {
			t.Errorf("Text expected: %s, actual: %s", string(expected.Words[i].Text), string(actual.Words[i].Text))
		}

		// check Progress
		if string(actual.Words[i].Progress) != string(expected.Words[i].Progress) {
			t.Errorf("Progress expected: %s, actual: %s", string(expected.Words[i].Progress), string(actual.Words[i].Progress))
		}

		// check Events
		for j := range actual.Words[i].Events {
			if actual.Words[i].Events[j].Key != expected.Words[i].Events[j].Key {
				t.Errorf("KeyEvent at position %d on word %d expected: %d, actual: %d", j, i, expected.Words[i].Events[j].Key, actual.Words[i].Events[j].Key)
			}
			if actual.Words[i].Events[j].ExpectedKey != expected.Words[i].Events[j].ExpectedKey {
				t.Errorf("KeyEvent at position %d on word %d expected: %v, actual: %v", j, i, expected.Words[i].Events[j].ExpectedKey, actual.Words[i].Events[j].ExpectedKey)
			}
		}
	}

	// check results
	if len(actual.Result.Missing) != len(expected.Result.Missing) {
		t.Errorf("Result.Missing expected number: %d, actual number: %d", len(expected.Result.Missing), len(actual.Result.Missing))
	}

	for i := range actual.Result.Missing {
		if actual.Result.Missing[i] != expected.Result.Missing[i] {
			t.Errorf("Result.Missing at position %d expected number: %s, actual number: %s", i, expected.Result.Missing, actual.Result.Missing)
		}
	}
}

// TODO: Add unit tests for Encode and Decode

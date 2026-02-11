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
								ExpectedKey: 'l',
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
			},
		},
		{
			"MissingKeysBecauseOfSpace",
			strings.NewReader("hello world"),
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

// func TestComputeResult(t *testing.T) {
// 	testCases := []struct {
// 		Name       string
// 		Text       io.Reader
// 		Keypresses []rune
// 		Session    Session
// 	}{
// 		{
// 			"MissingKeysInTheMiddleOfTyping",
// 			strings.NewReader("hello world"),
// 			[]rune{'h', 'e', 'x', 'l', 'o', ' ', 'X', 'o', 'r', 'l', 'd'},
// 			Session{
// 				CurrentWord: 0,
// 				Words: []Word{
// 					{
// 						Text:     []rune("hello"),
// 						Progress: []rune("hexlo"),
// 						Events: []KeyEvent{
// 							{
// 								Key:     'h',
// 								ExpectedKey: '',
// 							},
// 							{
// 								Key:     'e',
// 								ExpectedKey: '',
// 							},
// 							{
// 								Key:     'x',
// 								Correct: false,
// 							},
// 							{
// 								Key:     'l',
// 								ExpectedKey: '',
// 							},
// 							{
// 								Key:     'o',
// 								ExpectedKey: '',
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range testCases {
// 		t.Run(tt.Name, func(t *testing.T) {
// 			// instanciate session
// 			s, err := NewSession(tt.Text)
// 			if err != nil {
// 				t.Errorf("failed creating session: %s", err)
// 			}
// 			for _, kp := range tt.Keypresses {
// 				s.HandleKey(kp)
// 			}
// 			s1 := tt.Session
// 			CompareSessions(t, *s, s1)
// 		})
// 	}
// }

func CompareSessions(t *testing.T, actual Session, expected Session) {
	for i := range actual.Words {

		// check current word
		if actual.CurrentWord != expected.CurrentWord {
			t.Errorf("Current Word expected: %d, actual: %d", expected.CurrentWord, actual.CurrentWord)
		}

		// check Text
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
}

package typer

import (
	"io"
)

type App struct {
	sessions       []*Session
	currentSession int
}

func (app *App) HandleKey(key rune) {
	if session := app.getCurrentSession(); session != nil {
		session.HandleKey(key)
	}
}

func (app *App) DeleteWord() {
	if session := app.getCurrentSession(); session != nil {
		session.DeleteWord()
	}
}

func (app *App) DeleteChar() {
	if session := app.getCurrentSession(); session != nil {
		session.DeleteChar()
	}
}

func (app *App) CurrentWordIndex() int {
	if session := app.getCurrentSession(); session != nil {
		return session.CurrentWord
	}
	return -1
}

func (app *App) CreateSession(r io.Reader) error {
	session, err := NewSession(r)
	if err != nil {
		return err
	}
	app.AppendSession(session)
	return nil
}

func (app *App) LoadSession(r io.Reader) error {
	session, err := DecodeSession(r)
	if err != nil {
		return err
	}
	app.AppendSession(session)
	return nil
}

func (app *App) DumpSession(w io.Writer, index int) error {
	if 0 <= index && index < len(app.sessions) {
		return app.sessions[index].Encode(w)
	}
	return nil
}

func (app *App) ChooseSession(index int) {
	if 0 <= index && index < len(app.sessions) {
		app.currentSession = index
	}
}

func (app *App) AppendSession(session *Session) {
	app.sessions = append(app.sessions, session)
}

func (app *App) Words() []Word {
	if session := app.getCurrentSession(); session != nil {
		return session.Words
	}
	return []Word{}
}

func (app *App) getCurrentSession() *Session {
	if 0 <= app.currentSession && app.currentSession < len(app.sessions) {
		return app.sessions[app.currentSession]
	}
	return nil
}

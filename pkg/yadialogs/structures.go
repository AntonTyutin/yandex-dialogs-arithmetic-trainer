package yadialogs

type InputStruct interface {
	Request
	Session
	State
}

type StateValues map[string]any

type Input struct {
	Request Request `json:"request,omitempty"`
	Session Session `json:"session,omitempty"`
	State   State   `json:"state,omitempty"`
	Version string  `json:"version,omitempty"`
}

type Session struct {
	New       bool   `json:"new,omitempty"`
	MessageID int    `json:"message_id,omitempty"`
	SessionID string `json:"session_id,omitempty"`
	SkillID   string `json:"skill_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
}

type State struct {
	Session     StateValues `json:"session,omitempty"`
	User        StateValues `json:"user,omitempty"`
	Application StateValues `json:"application,omitempty"`
}

type Request struct {
	Command string `json:"command,omitempty"`
}

type Output struct {
	Response     Response    `json:"response,omitempty"`
	SessionState StateValues `json:"session_state,omitempty"`
	Version      string      `json:"version,omitempty"`
}

type Response struct {
	Text       string `json:"text,omitempty"`
	Tts        string `json:"tts,omitempty"`
	EndSession bool   `json:"end_session,omitempty"`
}

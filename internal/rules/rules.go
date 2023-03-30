package rules

import (
	"time"
)

type Rule struct {
	Text         string
	creator      string
	createdAt    time.Time
	lastChangeBy string
	lastChangeAt time.Time
}

func NewRule(text, creatorUserID string, createdAt time.Time) Rule {
	return Rule{
		Text:         text,
		creator:      creatorUserID,
		createdAt:    createdAt,
		lastChangeBy: creatorUserID,
		lastChangeAt: createdAt,
	}
}

func (r *Rule) Edit(text, userID string, changedAt time.Time) {
	r.Text = text
	r.lastChangeBy = userID
	r.lastChangeAt = changedAt
}

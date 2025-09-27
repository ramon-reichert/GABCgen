// Package errors defines domain-specific errors used across the application.
package errors

type DomainErr struct {
	Message string
}

func (e DomainErr) Error() string {
	return e.Message
}

var ErrShortPhrase = DomainErr{"the phrase is to short to apply the whole melody"}
var ErrShortParagraph = DomainErr{"each paragraph must have at least three phrases, not counting the conclusion phrase - which can start the last paragraph"}
var ErrNoText = DomainErr{"no incoming text to be parsed"}
var ErrNoLetters = DomainErr{"non-letter char not attached to any letter"}

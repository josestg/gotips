package how_to_redact_sensitive_data

import (
	"encoding"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"
)

// shouldRedacted is list of interfaces that should be implemented by types that need to be redacted.
type shouldRedacted interface {
	Redact() string
	fmt.Stringer
	encoding.TextMarshaler
	slog.LogValuer
}

type Password string

var _ shouldRedacted = Password("")

// Redact returns a redacted version of the password.
func (p Password) Redact() string { return "[REDACTED]" }

// String redacts the password for string representation.
func (p Password) String() string { return p.Redact() }

// MarshalText redacts the password for text representation. Includes in JSON marshaling.
func (p Password) MarshalText() ([]byte, error) { return []byte(p.Redact()), nil }

// LogValue returns a redacted version of the password for logging.
func (p Password) LogValue() slog.Value { return slog.StringValue(p.Redact()) }

type Redactor interface {
	Redact(string) (redacted string, ok bool)
}

type EmailRedactor struct {
	r Redactor
}

func NewEmailRedactor(r Redactor) *EmailRedactor {
	return &EmailRedactor{r: r}
}

func (r *EmailRedactor) Redact(email string) string {
	s, ok := r.r.Redact(email)
	if !ok {
		return email
	}
	return s
}

type SimpleEmailMaskerRedactor struct{}

func NewSimpleEmailMaskerRedactor() Redactor { return &SimpleEmailMaskerRedactor{} }

func (SimpleEmailMaskerRedactor) Redact(email string) (string, bool) {
	// mask email address by replacing the first 3 characters with "xxx" if more than 3 characters before the "@" symbol.
	i := strings.IndexByte(email, '@')
	if i < 4 {
		return email, false
	}
	return "xxx" + email[3:], true
}

var globalEmailRedactor atomic.Pointer[EmailRedactor]

// SetGlobalEmailRedactor sets the global email redactor.
func SetGlobalEmailRedactor(r *EmailRedactor) { globalEmailRedactor.Store(r) }
func GlobalEmailRedactor() *EmailRedactor     { return globalEmailRedactor.Load() }

type Email string

var _ shouldRedacted = Email("")

// Redact returns a redacted version of the email.
func (e Email) Redact() string { return GlobalEmailRedactor().Redact(string(e)) }

// String redacts the email for string representation.
func (e Email) String() string { return e.Redact() }

// MarshalText redacts the email for text representation. Includes in JSON marshaling.
func (e Email) MarshalText() ([]byte, error) { return []byte(e.Redact()), nil }

// LogValue returns a redacted version of the email for logging.
func (e Email) LogValue() slog.Value { return slog.StringValue(e.Redact()) }

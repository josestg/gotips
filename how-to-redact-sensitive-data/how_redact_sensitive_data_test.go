package how_to_redact_sensitive_data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
)

type User struct {
	ID       int      `json:"id"`
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

func TestMain(m *testing.M) {
	emailRedactor := NewEmailRedactor(NewSimpleEmailMaskerRedactor())
	SetGlobalEmailRedactor(emailRedactor)
	code := m.Run()
	os.Exit(code)
}

func TestRedact(t *testing.T) {
	usr := User{
		ID:       1,
		Email:    "alice@example.com",
		Password: Password("password123"),
	}

	str := fmt.Sprintf("%+v", usr)
	if str != "{ID:1 Email:xxxce@example.com Password:[REDACTED]}" {
		t.Errorf("unexpected redacted string: %s", str)
	}

	jsonBytes, err := json.Marshal(usr)
	if err != nil {
		t.Fatal(err)
	}

	jsonStr := string(jsonBytes)
	if jsonStr != `{"id":1,"email":"xxxce@example.com","password":"[REDACTED]"}` {
		t.Errorf("unexpected redacted JSON: %s", jsonStr)
	}

	var buffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buffer, &slog.HandlerOptions{}))
	logger.Info("test", "user", usr)

	logStr := buffer.String()
	if !strings.Contains(logStr, `user="{ID:1 Email:xxxce@example.com Password:[REDACTED]}"`) {
		t.Errorf("unexpected log: %s", logStr)
	}
}

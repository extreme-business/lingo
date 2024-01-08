package register

import (
	"encoding/hex"
	"fmt"
	"log/slog"
)

type LogRegisterHandler struct {
	Logger *slog.Logger
}

func (l *LogRegisterHandler) SendToken(email string, token []byte) error {
	if l.Logger == nil {
		return fmt.Errorf("logger is not set")
	}

	tokenStr := hex.EncodeToString(token)
	l.Logger.Info("SendToken", slog.String("email", email), slog.String("token", tokenStr))

	return nil
}

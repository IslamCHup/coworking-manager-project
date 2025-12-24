package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type SMSSender interface {
	Send(phone string, message string) error
}

type SmsAeroSender struct {
	email  string
	apiKey string
	from   string
	logger *slog.Logger
}

func NewSmsAeroSender(
	email string,
	apiKey string,
	from string,
	logger *slog.Logger,
) *SmsAeroSender {
	return &SmsAeroSender{
		email:  email,
		apiKey: apiKey,
		from:   from,
		logger: logger,
	}
}

func (s *SmsAeroSender) Send(phone, message string) error {
	phone = strings.TrimPrefix(phone, "+")

	payload := map[string]string{
		"number": phone,
		"text":   message,
		"sign":   s.from, // КЛЮЧЕВО
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://gate.smsaero.ru/v2/sms/send",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.SetBasicAuth(s.email, s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sms provider error: %v", result)
	}

	if success, ok := result["success"].(bool); ok && !success {
		return fmt.Errorf("sms provider error: %v", result)
	}

	return nil
}

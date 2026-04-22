package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// EmailService отвечает за отправку email через различных провайдеров
type EmailService struct {
	provider string
	apiKey   string
	fromAddr string
	client   *http.Client
}

// NewEmailService создает новый EmailService
func NewEmailService() *EmailService {
	provider := os.Getenv("EMAIL_PROVIDER")
	if provider == "" {
		provider = "mailtrap" // Default
	}

	return &EmailService{
		provider: provider,
		apiKey:   os.Getenv("EMAIL_API_KEY"),
		fromAddr: os.Getenv("EMAIL_FROM"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendVerificationCode отправляет код верификации на email
func (s *EmailService) SendVerificationCode(toEmail, code string, expiresAt time.Time) error {
	// Проверка конфигурации
	if s.apiKey == "" {
		return fmt.Errorf("EMAIL_API_KEY not configured")
	}
	if s.fromAddr == "" {
		return fmt.Errorf("EMAIL_FROM not configured")
	}

	// Генерируем HTML контент
	htmlContent, err := s.renderVerificationEmail(code, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Отправляем через соответствующий провайдер
	switch s.provider {
	case "sendgrid":
		return s.sendViaSendGrid(toEmail, "Код верификации - Объяснятель ДЗ", htmlContent)
	case "ses":
		return s.sendViaSES(toEmail, "Код верификации - Объяснятель ДЗ", htmlContent)
	case "mailgun":
		return s.sendViaMailgun(toEmail, "Код верификации - Объяснятель ДЗ", htmlContent)
	case "mailtrap":
		return s.sendViaMailtrap(toEmail, "Код верификации - Объяснятель ДЗ", htmlContent)
	default:
		return fmt.Errorf("unsupported email provider: %s", s.provider)
	}
}

// renderVerificationEmail генерирует HTML email с кодом верификации
func (s *EmailService) renderVerificationEmail(code string, expiresAt time.Time) (string, error) {
	tmpl := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Код верификации</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #2D3436;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #F0F4FF;
        }
        .container {
            background-color: white;
            border-radius: 12px;
            padding: 40px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            font-size: 48px;
            margin-bottom: 10px;
        }
        h1 {
            color: #6C5CE7;
            font-size: 24px;
            margin: 0;
        }
        .code-container {
            background: linear-gradient(135deg, #E8E4FF 0%, #DFE4FF 100%);
            border-radius: 12px;
            padding: 30px;
            text-align: center;
            margin: 30px 0;
        }
        .code {
            font-size: 36px;
            font-weight: 700;
            color: #6C5CE7;
            letter-spacing: 8px;
            font-family: 'Courier New', monospace;
        }
        .expires {
            color: #636E72;
            font-size: 14px;
            margin-top: 15px;
        }
        .message {
            color: #2D3436;
            font-size: 16px;
            margin: 20px 0;
        }
        .warning {
            background-color: #FFF3CD;
            border-left: 4px solid #FFC107;
            padding: 15px;
            margin: 20px 0;
            border-radius: 4px;
        }
        .warning-text {
            color: #856404;
            font-size: 14px;
            margin: 0;
        }
        .footer {
            text-align: center;
            color: #636E72;
            font-size: 13px;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #E8E4FF;
        }
        .footer a {
            color: #6C5CE7;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">📚</div>
            <h1>Объяснятель ДЗ</h1>
        </div>

        <p class="message">
            Здравствуйте!
        </p>

        <p class="message">
            Вы получили это письмо, потому что запросили код верификации для приложения <strong>Объяснятель ДЗ</strong>.
        </p>

        <div class="code-container">
            <div class="code">{{.Code}}</div>
            <div class="expires">Код действителен до {{.ExpiresAt}}</div>
        </div>

        <p class="message">
            Введите этот код в приложении, чтобы подтвердить ваш email адрес.
        </p>

        <div class="warning">
            <p class="warning-text">
                ⚠️ <strong>Никому не сообщайте этот код.</strong> Сотрудники Объяснятель ДЗ никогда не попросят вас предоставить код верификации.
            </p>
        </div>

        <p class="message">
            Если вы не запрашивали этот код, просто проигнорируйте это письмо.
        </p>

        <div class="footer">
            <p>
                С уважением,<br>
                Команда <strong>Объяснятель ДЗ</strong>
            </p>
            <p>
                <a href="mailto:support@obiasnyatel-dz.ru">support@obiasnyatel-dz.ru</a><br>
                <a href="https://vk.com/obiasnyatel_dz">vk.com/obiasnyatel_dz</a>
            </p>
        </div>
    </div>
</body>
</html>`

	type EmailData struct {
		Code      string
		ExpiresAt string
	}

	data := EmailData{
		Code:      code,
		ExpiresAt: expiresAt.Format("15:04 02.01.2006"),
	}

	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// sendViaSendGrid отправляет email через SendGrid API
func (s *EmailService) sendViaSendGrid(to, subject, htmlContent string) error {
	url := "https://api.sendgrid.com/v3/mail/send"

	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{"email": to},
				},
			},
		},
		"from": map[string]string{
			"email": s.fromAddr,
			"name":  "Объяснятель ДЗ",
		},
		"subject": subject,
		"content": []map[string]string{
			{
				"type":  "text/html",
				"value": htmlContent,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[EmailService] SendGrid error: %s", string(body))
		return fmt.Errorf("SendGrid returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("[EmailService] Email sent successfully to %s via SendGrid", to)
	return nil
}

// sendViaSES отправляет email через AWS SES
func (s *EmailService) sendViaSES(to, subject, htmlContent string) error {
	// TODO: Реализовать AWS SES при необходимости
	return fmt.Errorf("AWS SES not implemented yet")
}

// sendViaMailgun отправляет email через Mailgun
func (s *EmailService) sendViaMailgun(to, subject, htmlContent string) error {
	// TODO: Реализовать Mailgun при необходимости
	return fmt.Errorf("Mailgun not implemented yet")
}

// sendViaMailtrap отправляет email через Mailtrap API
// Документация: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/YXBpOjM1ODc4-send-email
func (s *EmailService) sendViaMailtrap(to, subject, htmlContent string) error {
	url := "https://send.api.mailtrap.io/api/send"

	payload := map[string]interface{}{
		"from": map[string]string{
			"email": s.fromAddr,
			"name":  "Объяснятель ДЗ",
		},
		"to": []map[string]string{
			{"email": to},
		},
		"subject": subject,
		"html":    htmlContent,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[EmailService] Mailtrap error: %s", string(body))
		return fmt.Errorf("Mailtrap returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("[EmailService] Email sent successfully to %s via Mailtrap", to)
	return nil
}

// IsDevelopmentMode проверяет режим разработки
func (s *EmailService) IsDevelopmentMode() bool {
	env := os.Getenv("ENV")
	return env == "development" || env == "dev" || env == ""
}

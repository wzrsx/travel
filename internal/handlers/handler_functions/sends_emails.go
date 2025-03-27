package handler_functions

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"
	"travel/internal/config"
)

var (
	codeCache    = make(map[string]codeInfo)
	attemptCache = make(map[string]attemptInfo)
	canChange = make(map[string]bool)
)

type codeInfo struct {
	code      string
	expiresAt time.Time
}

type attemptInfo struct {
	attempts int
	lastTry  time.Time
}

func SendEmailMessageWithCode() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			type CredentialsSendEmail struct {
				Email string `json:"email"`
			}

			var creds CredentialsSendEmail
			err := json.NewDecoder(r.Body).Decode(&creds)
			if err != nil {
				http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
				return
			}

			if(codeCache[creds.Email].expiresAt.IsZero() || canChange[creds.Email]){
				canChange[creds.Email] = false

				ip := r.RemoteAddr
				confirmationCode := generateSecureCode()
				codeCache[creds.Email] = codeInfo{
					code:      confirmationCode,
					expiresAt: time.Now().Add(3 * time.Minute),
				}
				attemptCache[ip] = attemptInfo{
					attempts: 1,
					lastTry:  time.Now(),
				}
	
				EmailData, err := config.NewConfigEmail()
				if err != nil {
					log.Printf("Error getting email data from cfg: %v", err)
				}
				// Аутентификация
				auth := smtp.PlainAuth("", EmailData.Sender, EmailData.Password, EmailData.SmtpHost)
	
				// Формируем красивое HTML-письмо
	
				subject := "Ваш код подтверждения"
				body := fmt.Sprintf(`
				<html>
				<body style="font-family: Arial, sans-serif; line-height: 1.6;">
					<div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px;">
						<h2 style="color: #444;">Подтверждение регистрации</h2>
						<p>Для завершения регистрации введите следующий код подтверждения:</p>
						
						<div style="background: #f5f5f5; padding: 15px; text-align: center; 
									margin: 20px 0; font-size: 24px; letter-spacing: 2px;
									border-radius: 5px; font-weight: bold;">
							%s
						</div>
						
						<p>Этот код действителен в течение 3 минут.</p>
						<p style="color: #888; font-size: 12px;">
							Если вы не запрашивали этот код, пожалуйста, проигнорируйте это письмо.
						</p>
					</div>
				</body>
				</html>
				`, confirmationCode)
	
				// Формируем заголовки письма
				headers := make(map[string]string)
				headers["From"] = EmailData.Sender
				headers["To"] = creds.Email
				headers["Subject"] = subject
				headers["MIME-Version"] = "1.0"
				headers["Content-Type"] = "text/html; charset=\"utf-8\""
	
				// Собираем все части письма
				message := ""
				for k, v := range headers {
					message += fmt.Sprintf("%s: %s\r\n", k, v)
				}
				message += "\r\n" + body
	
				// Отправка письма
				err = smtp.SendMail(
					EmailData.SmtpHost+":"+EmailData.SmtpPort,
					auth,
					EmailData.Sender,
					[]string{creds.Email},
					[]byte(message),
				)
				if err != nil {
					log.Printf("Error sending email: %v", err)
				}
			}
			
		}
	}
	return http.HandlerFunc(fn)
}

func generateSecureCode() string {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback если криптографический генератор не сработал
		return fmt.Sprintf("%06d", time.Now().Nanosecond()%1000000)
	}
	return fmt.Sprintf("%06d", uint32(b[0])<<16|uint32(b[1])<<8|uint32(b[2]))[:6]
}

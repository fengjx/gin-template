package email

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"

	"gin-template/internal/app/config"
	appLog "gin-template/internal/app/log"
	"gin-template/pkg/errs"
	"go.uber.org/zap"
)

var ErrProviderDisabled = errors.New("mail provider disabled")

func Send(_ context.Context, to, subject, body string) error {
	cfg := config.Get()
	if !cfg.Mail.Enabled {
		appLog.Info("mail skipped", zap.String("to", to), zap.String("subject", subject))
		return errs.WithStack(ErrProviderDisabled)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Mail.SMTPHost, cfg.Mail.SMTPPort)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))
	auth := smtp.PlainAuth("", cfg.Mail.Username, cfg.Mail.Password, cfg.Mail.SMTPHost)
	if err := smtp.SendMail(addr, auth, cfg.Mail.FromEmail, []string{to}, msg); err != nil {
		return errs.Wrap(err, "发送邮件失败")
	}
	return nil
}

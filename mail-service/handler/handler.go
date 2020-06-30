package handler

import (
	"context"
	"strings"

	log "github.com/micro/go-micro/v2/logger"
	pb "github.com/x-community/mail-service/proto"
	"gopkg.in/gomail.v2"
)

type service struct {
	id  string
	cfg Config
}

func (s *service) SendMail(ctx context.Context, in *pb.SendMailRequest, out *pb.SendMailResponse) error {
	if len(in.Address) <= 0 {
		return s.NewError(errInvalidEmailAddress)
	}
	if len(in.Content) <= 0 {
		return s.NewError(errInvalidEmailContent)
	}
	contentType := in.ContentType
	if len(contentType) == 0 {
		contentType = "text/html"
	}
	mail := gomail.NewMessage()
	mail.SetAddressHeader("From", s.cfg.Mail.Username, s.cfg.Mail.Sender)
	mail.SetAddressHeader("To", in.Address, in.Receiver)
	mail.SetHeader("Subject", in.Subject)
	mail.SetBody(contentType, in.Content)
	d := gomail.NewDialer(s.cfg.Mail.Host, s.cfg.Mail.Port, s.cfg.Mail.Username, s.cfg.Mail.Password)
	if err := d.DialAndSend(mail); err != nil {
		msg := err.Error()
		if strings.HasSuffix(msg, "550 Mailbox not found or access denied") {
			return s.NewError(errInvalidEmailAddress)
		}
		log.Error(msg)
		return s.InternalServerError(msg)
	}
	return nil
}

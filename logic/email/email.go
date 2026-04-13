package email

import (
	"IM_chat/dao"
	"IM_chat/pkg/errcode"
	"IM_chat/pkg/redis"
	"context"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"net/smtp"
	"strconv"
	"time"
)

type UserSendConfirmEmailService struct {
	UserEmail string `json:"email" form:"email" binding:"required"`
}

// SendConfirmMessage 传入目标邮箱以及对应的code码即可发送验证码邮箱
func SendConfirmMessage(targetMailBox string, code string) error {
	from := viper.GetString("email.from")
	smtpHost := viper.GetString("email.smtp_host")
	smtpPort := viper.GetInt("email.smtp_port")
	secretKey := viper.GetString("email.secret_key")
	if from == "" || smtpHost == "" || secretKey == "" {
		return fmt.Errorf("email config missing: check email.from, email.smtp_host, email.secret_key in config.yaml")
	}
	em := email.NewEmail()
	em.From = fmt.Sprintf("IM-Chat <%s>", from)
	em.To = []string{targetMailBox}
	em.Subject = "验证码为: " + code
	emailContentCode := "您IM聊天室的验证码是 " + code + " 。它将在30min后过期。"
	em.Text = []byte(emailContentCode)
	addr := smtpHost + ":" + strconv.Itoa(smtpPort)
	auth := smtp.PlainAuth("", from, secretKey, smtpHost)
	return em.Send(addr, auth)
}
func (service *UserSendConfirmEmailService) SendConfirmEmail() string {
	if !dao.VerifyEmailFormat(service.UserEmail) {
		return errcode.Msg(errcode.InvalidEmail)
	}
	if redis.RDB.Get(context.Background(), "send-email:"+service.UserEmail).Val() != "" {
		return errcode.Msg(errcode.HasSendCode)
	}
	code := dao.GetConfirmCode()
	err := SendConfirmMessage(service.UserEmail, code)
	if err != nil {
		return errcode.Msg(errcode.DontSendCode)
	}
	if err := redis.RDB.Set(context.Background(), "email:"+service.UserEmail, code, time.Minute*30).Err(); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err := redis.RDB.Set(context.Background(), "send-email:"+service.UserEmail, code, time.Minute*1).Err(); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	return errcode.Msg(errcode.SUCCESS)
}

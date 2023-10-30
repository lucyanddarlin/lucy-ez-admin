package email

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"net/smtp"
	"os"
	"strings"
	"sync"

	"github.com/lucyanddarlin/lucy-ez-admin/config"
)

type email struct {
	mu       sync.RWMutex
	template map[string]struct {
		subject string
		html    string
	}
	user     string
	host     string
	password string
	company  string
}

type sender struct {
	tp string
	*email
}

type Sender interface {
	Send(email string, data any) error
	SendAll(emails []string, data any) error
}

type Email interface {
	NewSender(tpName string) Sender
}

// New 初始化 email 实例
func New(conf *config.Email) Email {
	emailIns := email{
		mu:       sync.RWMutex{},
		user:     conf.User,
		host:     conf.Host,
		password: conf.Password,
		company:  conf.Company,
		template: map[string]struct {
			subject string
			html    string
		}{},
	}

	emailIns.mu.Lock()
	defer emailIns.mu.Unlock()

	for _, item := range conf.Template {
		file, err := os.Open(item.Src)
		if err != nil {
			panic("邮箱模板初始化失败" + err.Error())
		}
		val, err := io.ReadAll(file)
		if err != nil {
			panic("邮箱模板读取失败" + err.Error())
		}
		emailIns.template[item.Name] = struct {
			subject string
			html    string
		}{
			subject: item.Subject, html: string(val),
		}
	}

	return &emailIns
}

// NewSender implements Email.
//
//	@Description: 新建一个发送器
//	@receiver e
//	@param tpName 需要发送的模板名
//	@return Sender
func (e *email) NewSender(tpName string) Sender {
	return &sender{
		email: e,
		tp:    tpName,
	}
}

// Send implements Sender.
//
//	@Description: 向指定的邮箱发送邮件
//	@receiver e
//	@param email 需要发送的邮箱
//	@param data 模板参数
//	@return error
func (s *sender) Send(email string, data any) error {
	subject, htmlTemplate, has := s.getTemplate()
	if !has {
		return errors.New("no exist template")
	}
	html, err := s.parseTemplate(htmlTemplate, data)
	if err != nil {
		return err
	}
	hp := strings.Split(s.host, ":")
	auth := smtp.PlainAuth("", s.user, s.password, hp[0])
	ct := "Context-Type: text/html; charset=UTF-8"
	msg := []byte("To: " + email + "\r\nFrom: " + s.user + "\r\nSubject: " + subject + "\r\n" + ct + "\r\n\r\n" + html)
	return smtp.SendMail(s.host, auth, s.user, []string{email}, msg)
}

// SendAll implements Sender.
func (*sender) SendAll(emails []string, data any) error {
	return nil
}

// getTemplate
//
//	@Description: 获取指定模板
//	@receiver s
//	@return string string bool
func (s *sender) getTemplate() (string, string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tp, is := s.email.template[s.tp]
	return tp.subject, tp.html, is
}

// parseTemplate
//
//	@Description: 解析模板变量
//	@receiver s
//	@param tp
//	@param data
//	@return string error
func (s *sender) parseTemplate(tp string, data any) (string, error) {
	n := template.New("")
	t, err := n.Parse(tp)
	if err != nil {
		return "", err
	}
	html := bytes.NewBuffer([]byte(""))
	if err := t.Execute(html, data); err != nil {
		return "", err
	}
	return html.String(), nil
}

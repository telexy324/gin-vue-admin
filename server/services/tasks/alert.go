package tasks

import (
	"bytes"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/mail"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const emailTemplate = `Subject: Task '{{ .Name }}' failed

Task {{ .TaskID }} with template '{{ .Name }}' has failed!
Task log: <a href='{{ .TaskURL }}'>{{ .TaskURL }}</a>`

const telegramTemplate = `{"chat_id": "{{ .ChatID }}","parse_mode":"HTML","text":"<code>{{ .Name }}</code>\n#{{ .TaskID }} <b>{{ .TaskResult }}</b> <code>{{ .TaskVersion }}</code> {{ .TaskDescription }}\nby {{ .Author }}\n{{ .TaskURL }}"}`

// Alert represents an alert that will be templated and sent to the appropriate service
type Alert struct {
	TaskID          string
	Name            string
	TaskURL         string
	ChatID          string
	TaskResult      string
	TaskDescription string
	TaskVersion     string
	Author          string
}

func (t *TaskRunner) sendMailAlert() {
	if !global.GVA_CONFIG.Ansible.EmailAlert || !t.alert {
		return
	}

	mailHost := global.GVA_CONFIG.Ansible.EmailHost + ":" + global.GVA_CONFIG.Ansible.EmailPort

	var mailBuffer bytes.Buffer
	alert := Alert{
		TaskID:  strconv.Itoa(int(t.task.ID)),
		Name:    t.template.Name,
		TaskURL: global.GVA_CONFIG.Ansible.WebHost + "/project/" + strconv.Itoa(t.template.ProjectID),
	}
	tpl := template.New("mail body template")
	tpl, err := tpl.Parse(emailTemplate)
	if err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Any("level", "Error"))
	}

	t.panicOnError(tpl.Execute(&mailBuffer, alert), "Can't generate alert template!")

	for _, user := range t.users {
		err, userObj := systemUserService.FindUserById(user)

		t.Log("Sending email to " + userObj.Email + " from " + global.GVA_CONFIG.Ansible.EmailSender)
		if global.GVA_CONFIG.Ansible.EmailSecure {
			err = mail.SendSecureMail(global.GVA_CONFIG.Ansible.EmailHost, global.GVA_CONFIG.Ansible.EmailPort, global.GVA_CONFIG.Ansible.EmailSender, global.GVA_CONFIG.Ansible.EmailUsername, global.GVA_CONFIG.Ansible.EmailPassword, userObj.Email, mailBuffer)
		} else {
			err = mail.SendMail(mailHost, global.GVA_CONFIG.Ansible.EmailSender, userObj.Email, mailBuffer)
		}
		t.panicOnError(err, "Can't send email!")
	}
}

func (t *TaskRunner) sendTelegramAlert() {
	if !global.GVA_CONFIG.Ansible.TelegramAlert || !t.alert {
		return
	}

	if t.template.SuppressSuccessAlerts && t.task.Status == ansible.TaskSuccessStatus {
		return
	}

	chatID := global.GVA_CONFIG.Ansible.TelegramChat
	if t.alertChat != nil && *t.alertChat != "" {
		chatID = *t.alertChat
	}

	var telegramBuffer bytes.Buffer

	var version string
	if t.task.Version != nil {
		version = *t.task.Version
	} else if t.task.BuildTaskID != nil {
		version = "build " + strconv.Itoa(*t.task.BuildTaskID)
	} else {
		version = ""
	}

	var message string
	if t.task.Message != "" {
		message = "- " + t.task.Message
	}

	var author string
	if t.task.UserID != nil {
		err, user := systemUserService.FindUserById(*t.task.UserID)
		if err != nil {
			panic(err)
		}
		author = user.Username
	}

	alert := Alert{
		TaskID:          strconv.Itoa(int(t.task.ID)),
		Name:            t.template.Name,
		TaskURL:         global.GVA_CONFIG.Ansible.WebHost + "/project/" + strconv.Itoa(t.template.ProjectID) + "/templates/" + strconv.Itoa(int(t.template.ID)) + "?t=" + strconv.Itoa(int(t.task.ID)),
		ChatID:          chatID,
		TaskResult:      strings.ToUpper(string(t.task.Status)),
		TaskVersion:     version,
		TaskDescription: message,
		Author:          author,
	}

	tpl := template.New("telegram body template")

	tpl, err := tpl.Parse(telegramTemplate)
	if err != nil {
		t.Log("Can't parse telegram template!")
		panic(err)
	}

	err = tpl.Execute(&telegramBuffer, alert)
	if err != nil {
		t.Log("Can't generate alert template!")
		panic(err)
	}

	resp, err := http.Post("https://api.telegram.org/bot"+global.GVA_CONFIG.Ansible.TelegramToken+"/sendMessage", "application/json", &telegramBuffer)

	if err != nil {
		t.Log("Can't send telegram alert! Response code not 200!")
	} else if resp.StatusCode != 200 {
		t.Log("Can't send telegram alert! Response code not 200!")
	}
}

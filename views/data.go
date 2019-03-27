package views

import "log"

//Set predefined Alert levels
const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"
	AlertMsgGeneric = "Something went wrong. Please try " +
		"again, and contact us if the problem persists."
)

//Data is the top level data structure that views expect data to come in
type Data struct {
	Alert   *Alert
	Content interface{}
}

//Alert is used to render bootstrap alert messages in templates
type Alert struct {
	Level   string
	Message string
}

type PublicError interface {
	error
	Public() string
}

//SetAlert is a type assertion, checks whether an error is a public error or not public error
func (d *Data) SetAlert(err error) {
	var msg string
	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

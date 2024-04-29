package config

import (
	"bytes"
	"html/template"
	"os"
	"path"

	"github.com/pelletier/go-toml/v2"
)

const (
	MsgTypeInfo  = "info"
	MsgTypeWarn  = "warn"
	MsgTypeError = "error"
)

type MessageRawTemplate struct {
	Info  map[string]string `toml:"info"`
	Warn  map[string]string `toml:"warn"`
	Error map[string]string `toml:"error"`
}

// first key: message type(info/warn/error), second key: message id
var msgs map[string]map[string]*template.Template

func init() {
	msgs = map[string]map[string]*template.Template{}
}

func LoadMessages(dir, lang string) error {
	_, err := os.Stat(path.Join(dir, lang))
	if err != nil {
		return err
	}
	arr, err := os.ReadFile(path.Join(dir, lang, "msg.toml"))
	if err != nil {
		return err
	}
	messagesRow := &MessageRawTemplate{}
	if err = toml.Unmarshal(arr, messagesRow); err != nil {
		return err
	}
	// info msg
	msgs[MsgTypeInfo] = make(map[string]*template.Template, len(messagesRow.Info))
	for msgID, rowTmpl := range messagesRow.Info {
		tmpl, err := template.New(msgID).Parse(rowTmpl)
		if err != nil {
			continue
		}
		msgs[MsgTypeInfo][msgID] = tmpl
	}
	// warn msg
	msgs[MsgTypeWarn] = make(map[string]*template.Template, len(messagesRow.Warn))
	for msgID, rowTmpl := range messagesRow.Warn {
		tmpl, err := template.New(msgID).Parse(rowTmpl)
		if err != nil {
			continue
		}
		msgs[MsgTypeWarn][msgID] = tmpl
	}
	// error msg
	msgs[MsgTypeError] = make(map[string]*template.Template, len(messagesRow.Error))
	for msgID, rowTmpl := range messagesRow.Error {
		tmpl, err := template.New(msgID).Parse(rowTmpl)
		if err != nil {
			continue
		}
		msgs[MsgTypeError][msgID] = tmpl
	}
	return nil
}

func GetMsg(msgType, msgID string, data map[string]any) string {
	if msgTmplMap, ok := msgs[msgType]; ok {
		if msgTmpl, ok := msgTmplMap[msgID]; ok {
			bf := new(bytes.Buffer)
			msgTmpl.Execute(bf, data)
			if bf.Len() != 0 {
				return bf.String()
			}
		}
	}
	return ""
}

package chat

import (
	"bytes"
	"encoding/json"
	"io"
)

// InitChatLogger sets the writer for chat logs. By default, logs are discarded.
func InitChatLogger(writer io.Writer) {
	chatLogger = writer
}

// chatLogger is the writer for chat logs. It is initialized to discard logs by default.
var chatLogger io.Writer = io.Discard

func logChat(message string, object any) {
	if chatLogger == io.Discard || chatLogger == nil {
		return
	}
	var buf bytes.Buffer
	_, err := buf.WriteString(message + "\n")
	if err != nil {
		return
	}
	if object != nil {
		data, err := json.MarshalIndent(object, "", "  ")
		if err != nil {
			_, err = buf.WriteString("error marshaling object: " + err.Error())
			if err != nil {
				return
			}
		} else {
			_, err = buf.Write(data)
			if err != nil {
				return
			}
		}
		_, err = buf.WriteString("\n")
		if err != nil {
			return
		}
	}

	_, err = chatLogger.Write(buf.Bytes())
	if err != nil {
		return
	}
}

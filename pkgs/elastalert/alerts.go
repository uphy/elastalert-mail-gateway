package elastalert

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/uphy/elastalert-mail-gateway/pkgs/jsonutil"
)

type (
	Alerts []*Alert
	Alert  struct {
		Body string          `json:"body"`
		Doc  jsonutil.Object `json:"doc"`
	}
)

const AlertSeparator = "\n----------------------------------------\n"

func (a Alert) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(a.Body)
	buf.WriteString("\n\n")
	for k, v := range a.Doc {
		buf.WriteString(k)
		buf.WriteString(": ")
		b, err := json.MarshalIndent(v, "", "  ")
		var value string
		if err != nil {
			value = fmt.Sprintf("<failed to marshal:%v>", err)
		} else {
			value = string(b)
		}
		buf.WriteString(value)
		buf.WriteString("\n")
	}
	return buf.String()
}

func ParseMailBody(body string) (Alerts, error) {
	entries := strings.Split(body, AlertSeparator)
	v := []*Alert{}
	for _, entry := range entries {
		entry = strings.TrimRight(entry, "\n")
		if len(entry) == 0 {
			continue
		}
		e, err := parseMailBodyEntry(entry)
		if err != nil {
			return nil, err
		}
		v = append(v, e)
	}
	return v, nil
}

func parseMailBodyEntry(body string) (*Alert, error) {
	reader := bufio.NewReader(strings.NewReader(body))
	bodyText := new(bytes.Buffer)
	jsonBody := new(bytes.Buffer)
	waitForBraces := false
	waitForBrackets := false

l:
	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		line := string(b)

		if waitForBraces || waitForBrackets {
			// json object or json array content
			switch line {
			case "]":
				if waitForBrackets {
					waitForBrackets = false
					jsonBody.WriteString("]")
					continue l
				}
			case "}":
				if waitForBraces {
					waitForBraces = false
					jsonBody.WriteString("}")
					continue l
				}
			}
			if waitForBraces || waitForBrackets {
				jsonBody.WriteString(line)
				jsonBody.WriteString("\n")
			}
		} else {
			// json key value or message body
			colonIndex := strings.Index(line, ":")
			if colonIndex > 0 {
				if jsonBody.Len() > 0 {
					if !waitForBraces && !waitForBrackets {
						jsonBody.WriteString(",")
					}
					jsonBody.WriteString("\n")
				} else {
					jsonBody.WriteString("{")
				}

				key := line[0:colonIndex]
				value := line[colonIndex+2:]
				switch value {
				case "[":
					waitForBrackets = true
				case "{":
					waitForBraces = true
				}
				jsonBody.WriteString(fmt.Sprintf(`"%s": %s`, key, value))
			} else {
				if bodyText.Len() > 0 {
					bodyText.WriteString("\n")
				}
				bodyText.WriteString(line)
			}
		}
	}
	var v map[string]interface{}
	if jsonBody.Len() > 0 {
		jsonBody.WriteString("}")
		if err := yaml.Unmarshal(jsonBody.Bytes(), &v); err != nil {
			return nil, err
		}
	} else {
		v = make(map[string]interface{}, 0)
	}
	return &Alert{
		Body: strings.TrimRight(bodyText.String(), "\n"),
		Doc:  v,
	}, nil
}

func (m Alerts) equalsTo(o Alerts) (bool, string) {
	if len(m) != len(o) {
		return false, fmt.Sprintf("entry length not equals:%d,%d", len(m), len(o))
	}
	for i, e1 := range m {
		e2 := o[i]
		equal, desc := e1.equalsTo(e2)
		if !equal {
			return false, fmt.Sprintf("entry %d not equals:%s", i, desc)
		}
	}
	return true, ""
}

func (m *Alert) equalsTo(o *Alert) (bool, string) {
	if m.Body != o.Body {
		return false, fmt.Sprintf("body not equals:\n%s\n---------------------\n%s", m.Body, o.Body)
	}
	doc1, _ := yaml.Marshal(m.Doc)
	doc2, _ := yaml.Marshal(o.Doc)
	doc1s := string(doc1)
	doc2s := string(doc2)
	if doc1s != doc2s {
		return false, fmt.Sprintf("doc not equals:\n%s\n----------------------\n%s", doc1s, doc2s)
	}
	return true, ""
}

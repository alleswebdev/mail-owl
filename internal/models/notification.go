package models

import "encoding/json"

type Notification struct {
	Type        string
	Data        string
	Project     string
	Prefix      string
	Layout      string `json:"layout"`
	NoLayout    bool   `json:"noLayout"`
	Subject     string
	Template    string
	Attachments []Attachment
	Params      map[string]interface{}
}

type Attachment struct {
	Container string
	Filename  string
	Url       string
}

func (n *Notification) Fill() error {
	return json.Unmarshal([]byte(n.Data), &n)
}

const AttachmentTemplate = `
<html>
<body>
<h3>Вложения:</h3>

{{ range $key, $value := . }}
    <a href="{{$value.Url}}" target="_blank"> {{$value.Filename}}</a> <br>
{{end}}

</body>
</html>
`

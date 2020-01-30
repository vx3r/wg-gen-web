package util

import (
	"bytes"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/model"
	"strings"
	"text/template"
)

var (
	clientTpl = `
[Interface]
Address = {{.Client.Address}}
PrivateKey = {{.Client.PrivateKey}}
DNS = {{.Server.Dns}}
[Peer]
PublicKey = {{.Server.PublicKey}}
PresharedKey = {{.Server.PresharedKey}}
AllowedIPs = {{.Client.AllowedIPs}}
Endpoint = {{.Server.Endpoint}}
PersistentKeepalive = {{.Server.PersistentKeepalive}}`

	wgTpl = `
# {{.Server.Name}} / Updated: {{.Server.Updated}} / Created: {{.Server.Created}}
[Interface]
	{{range .ServerAdresses}}
Address = {{.}}
	{{end}}
ListenPort = {{.Server.ListenPort}}
PrivateKey = {{.Server.PrivateKey}}
	{{$server := .Server}}
	{{range .Clients}}
		{{if .Enable}}
# {{.Name}} / {{.Email}} / Updated: {{.Updated}} / Created: {{.Created}}
[Peer]
PublicKey = {{.PublicKey}}
PresharedKey = {{$server.PresharedKey}}
AllowedIPs = {{.Address}}
		{{end}}
	{{end}}`
)

func DumpClient(client *model.Client, server *model.Server) (bytes.Buffer, error) {
	var tplBuff bytes.Buffer

	t, err := template.New("client").Parse(clientTpl)
	if err != nil {
		return tplBuff, err
	}

	return dump(t, struct {
		Client *model.Client
		Server *model.Server
	}{
		Client: client,
		Server: server,
	})
}

func DumpServerWg(clients []*model.Client, server *model.Server) (bytes.Buffer, error) {
	var tplBuff bytes.Buffer

	t, err := template.New("server").Parse(wgTpl)
	if err != nil {
		return tplBuff, err
	}

	return dump(t, struct {
		Clients []*model.Client
		Server *model.Server
		ServerAdresses []string
	}{
		ServerAdresses: strings.Split(server.Address, ","),
		Clients: clients,
		Server: server,
	})
}

func dump(tpl *template.Template , data interface{}) (bytes.Buffer, error) {
	var tplBuff bytes.Buffer

	err := tpl.Execute(&tplBuff, data)
	if err != nil {
		return tplBuff, err
	}

	return tplBuff, nil
}

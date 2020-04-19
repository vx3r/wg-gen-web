package core

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/model"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/storage"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/template"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/util"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// CreateClient client with all necessary data
func CreateClient(client *model.Client) (*model.Client, error) {
	// check if client is valid
	errs := client.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("client validation error")
		}
		return nil, errors.New("failed to validate client")
	}

	u := uuid.NewV4()
	client.Id = u.String()

	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	client.PrivateKey = key.String()
	client.PublicKey = key.PublicKey().String()

	presharedKey, err := wgtypes.GenerateKey()
	if err != nil {
		return nil, err
	}
	client.PresharedKey = presharedKey.String()

	reserverIps, err := GetAllReservedIps()
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)
	for _, network := range client.Address {
		ip, err := util.GetAvailableIp(network, reserverIps)
		if err != nil {
			return nil, err
		}
		if util.IsIPv6(ip) {
			ip = ip + "/128"
		} else {
			ip = ip + "/32"
		}
		ips = append(ips, ip)
	}
	client.Address = ips
	client.Created = time.Now().UTC()
	client.Updated = client.Created

	err = storage.Serialize(client.Id, client)
	if err != nil {
		return nil, err
	}

	v, err := storage.Deserialize(client.Id)
	if err != nil {
		return nil, err
	}
	client = v.(*model.Client)

	// data modified, dump new config
	return client, UpdateServerConfigWg()
}

// ReadClient client by id
func ReadClient(id string) (*model.Client, error) {
	v, err := storage.Deserialize(id)
	if err != nil {
		return nil, err
	}
	client := v.(*model.Client)

	return client, nil
}

// UpdateClient preserve keys
func UpdateClient(Id string, client *model.Client) (*model.Client, error) {
	v, err := storage.Deserialize(Id)
	if err != nil {
		return nil, err
	}
	current := v.(*model.Client)

	if current.Id != client.Id {
		return nil, errors.New("records Id mismatch")
	}

	// check if client is valid
	errs := client.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("client validation error")
		}
		return nil, errors.New("failed to validate client")
	}

	// keep keys
	client.PrivateKey = current.PrivateKey
	client.PublicKey = current.PublicKey
	client.Updated = time.Now().UTC()

	err = storage.Serialize(client.Id, client)
	if err != nil {
		return nil, err
	}

	v, err = storage.Deserialize(Id)
	if err != nil {
		return nil, err
	}
	client = v.(*model.Client)

	// data modified, dump new config
	return client, UpdateServerConfigWg()
}

// DeleteClient from disk
func DeleteClient(id string) error {
	path := filepath.Join(os.Getenv("WG_CONF_DIR"), id)
	err := os.Remove(path)
	if err != nil {
		return err
	}

	// data modified, dump new config
	return UpdateServerConfigWg()
}

// ReadClients all clients
func ReadClients() ([]*model.Client, error) {
	clients := make([]*model.Client, 0)

	files, err := ioutil.ReadDir(filepath.Join(os.Getenv("WG_CONF_DIR")))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		// clients file name is an uuid
		_, err := uuid.FromString(f.Name())
		if err == nil {
			c, err := storage.Deserialize(f.Name())
			if err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"path": f.Name(),
				}).Error("failed to deserialize client")
			} else {
				clients = append(clients, c.(*model.Client))
			}
		}
	}

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].Created.After(clients[j].Created)
	})

	return clients, nil
}

// ReadClientConfig in wg format
func ReadClientConfig(id string) ([]byte, error) {
	client, err := ReadClient(id)
	if err != nil {
		return nil, err
	}

	server, err := ReadServer()
	if err != nil {
		return nil, err
	}

	configDataWg, err := template.DumpClientWg(client, server)
	if err != nil {
		return nil, err
	}

	return configDataWg, nil
}

// EmailClient send email to client
func EmailClient(id string) error {
	client, err := ReadClient(id)
	if err != nil {
		return err
	}

	configData, err := ReadClientConfig(id)
	if err != nil {
		return err
	}

	// conf as .conf file
	tmpfileCfg, err := ioutil.TempFile("", "wireguard-vpn-*.conf")
	if err != nil {
		return err
	}
	if _, err := tmpfileCfg.Write(configData); err != nil {
		return err
	}
	if err := tmpfileCfg.Close(); err != nil {
		return err
	}
	defer os.Remove(tmpfileCfg.Name()) // clean up

	// conf as png image
	png, err := qrcode.Encode(string(configData), qrcode.Medium, 280)
	if err != nil {
		return err
	}
	tmpfilePng, err := ioutil.TempFile("", "qrcode-*.png")
	if err != nil {
		return err
	}
	if _, err := tmpfilePng.Write(png); err != nil {
		return err
	}
	if err := tmpfilePng.Close(); err != nil {
		return err
	}
	defer os.Remove(tmpfilePng.Name()) // clean up

	// get email body
	emailBody, err := template.DumpEmail(client, filepath.Base(tmpfilePng.Name()))
	if err != nil {
		return err
	}

	// port to int
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
	s, err := d.Dial()
	if err != nil {
		return err
	}
	m := gomail.NewMessage()

	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetAddressHeader("To", client.Email, client.Name)
	m.SetHeader("Subject", "WireGuard VPN Configuration")
	m.SetBody("text/html", string(emailBody))
	m.Attach(tmpfileCfg.Name())
	m.Embed(tmpfilePng.Name())

	err = gomail.Send(s, m)
	if err != nil {
		return err
	}

	return nil
}

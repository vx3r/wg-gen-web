package core

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/model"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/storage"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Migrate all changes, current struct fields change
func Migrate() error {
	clients, err := readClients()
	if err != nil {
		return err
	}

	s, err := deserialize("server.json")
	if err != nil {
		return err
	}

	for _, client := range clients {
		switch v := client["allowedIPs"].(type) {
		case []interface{}:
			log.Infof("client %s has been already migrated", client["id"])
			continue
		default:
			log.Infof("unexpected type %T, mus be migrated", v)
		}

		c := &model.Client{}
		c.Id = client["id"].(string)
		c.Name = client["name"].(string)
		c.Email = client["email"].(string)
		c.Enable = client["enable"].(bool)
		c.AllowedIPs = make([]string, 0)
		for _, address := range strings.Split(client["allowedIPs"].(string), ",") {
			if util.IsValidCidr(strings.TrimSpace(address)) {
				c.AllowedIPs = append(c.AllowedIPs, strings.TrimSpace(address))
			}
		}
		c.Address = make([]string, 0)
		for _, address := range strings.Split(client["address"].(string), ",") {
			if util.IsValidCidr(strings.TrimSpace(address)) {
				c.Address = append(c.Address, strings.TrimSpace(address))
			}
		}
		c.PrivateKey = client["privateKey"].(string)
		c.PublicKey = client["publicKey"].(string)
		created, err := time.Parse(time.RFC3339, client["created"].(string))
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Errorf("failed to parse time")
			continue
		}
		c.Created = created
		updated, err := time.Parse(time.RFC3339, client["updated"].(string))
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Errorf("failed to parse time")
			continue
		}
		c.Updated = updated

		err = storage.Serialize(c.Id, c)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Errorf("failed to Serialize client")
		}
	}

	switch v := s["address"].(type) {
	case []interface{}:
		log.Info("server has been already migrated")
		return nil
	default:
		log.Infof("unexpected type %T, mus be migrated", v)
	}

	server := &model.Server{}

	server.Address = make([]string, 0)
	for _, address := range strings.Split(s["address"].(string), ",") {
		if util.IsValidCidr(strings.TrimSpace(address)) {
			server.Address = append(server.Address, strings.TrimSpace(address))
		}
	}
	server.ListenPort = int(s["listenPort"].(float64))
	server.PrivateKey = s["privateKey"].(string)
	server.PublicKey = s["publicKey"].(string)
	server.PresharedKey = s["presharedKey"].(string)
	server.Endpoint = s["endpoint"].(string)
	server.PersistentKeepalive = int(s["persistentKeepalive"].(float64))
	server.Dns = make([]string, 0)
	for _, address := range strings.Split(s["dns"].(string), ",") {
		if util.IsValidIp(strings.TrimSpace(address)) {
			server.Dns = append(server.Dns, strings.TrimSpace(address))
		}
	}
	if val, ok := s["preUp"]; ok {
		server.PreUp = val.(string)
	}
	if val, ok := s["postUp"]; ok {
		server.PostUp = val.(string)
	}
	if val, ok := s["preDown"]; ok {
		server.PreDown = val.(string)
	}
	if val, ok := s["postDown"]; ok {
		server.PostDown = val.(string)
	}
	created, err := time.Parse(time.RFC3339, s["created"].(string))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to parse time")
	}
	server.Created = created
	updated, err := time.Parse(time.RFC3339, s["updated"].(string))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to parse time")
	}
	server.Updated = updated

	err = storage.Serialize("server.json", server)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to Serialize server")
	}

	return nil
}

func readClients() ([]map[string]interface{}, error) {
	clients := make([]map[string]interface{}, 0)

	files, err := ioutil.ReadDir(filepath.Join(os.Getenv("WG_CONF_DIR")))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		// clients file name is an uuid
		_, err := uuid.FromString(f.Name())
		if err == nil {
			c, err := deserialize(f.Name())
			if err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"path": f.Name(),
				}).Error("failed to deserialize client")
			} else {
				clients = append(clients, c)
			}
		}
	}

	return clients, nil
}

func deserialize(id string) (map[string]interface{}, error) {
	path := filepath.Join(os.Getenv("WG_CONF_DIR"), id)

	data, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var d map[string]interface{}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

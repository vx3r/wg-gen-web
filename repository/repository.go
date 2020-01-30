package repository

import (
	"encoding/json"
	"errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/model"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/util"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

/*
 * CreateClient client with all necessary data
 */
func CreateClient(client *model.Client) (*model.Client, error) {
	u := uuid.NewV4()
	client.Id = u.String()

	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	client.PrivateKey = key.String()
	client.PublicKey = key.PublicKey().String()

	// find available IP address from selected networks
	clients, err := ReadClients()
	if err != nil {
		return nil, err
	}

	reserverIps := make([]string, 0)
	for _, client := range clients {
		ips := strings.Split(client.Address, ",")
		for i := range ips {
			if util.IsIPv6(ips[i]) {
				ips[i] = strings.ReplaceAll(strings.TrimSpace(ips[i]), "/128", "")
			} else {
				ips[i] = strings.ReplaceAll(strings.TrimSpace(ips[i]), "/32", "")
			}
		}
		reserverIps = append(reserverIps, ips...)
	}

	networks := strings.Split(client.Address, ",")
	for i := range networks {
		networks[i] = strings.TrimSpace(networks[i])
	}
	ips := make([]string, 0)
	for _, network := range networks {
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
	client.Address = strings.Join(ips, ",")

	client.Created = time.Now().UTC()
	client.Updated = client.Created

	err = serialize(client.Id, client)
	if err != nil {
		return nil, err
	}

	v, err := deserialize(client.Id)
	if err != nil {
		return nil, err
	}
	client = v.(*model.Client)

	return client, nil
}

/*
 * ReadClient client by id
 */
func ReadClient(id string) (*model.Client, error) {
	v, err := deserialize(id)
	if err != nil {
		return nil, err
	}
	client := v.(*model.Client)

	return client, nil
}

/*
 * ReadClientConfig in wg format
 */
func ReadClientConfig(id string) ([]byte, error) {
	client, err := ReadClient(id)
	if err != nil {
		return nil, err
	}

	server, err := ReadServer()
	if err != nil {
		return nil, err
	}

	configDataWg, err := util.DumpClient(client, server)
	if err != nil {
		return nil, err
	}

	return configDataWg.Bytes(), nil
}

/*
 * UpdateClient preserve keys
 */
func UpdateClient(Id string, client *model.Client) (*model.Client, error) {
	v, err := deserialize(Id)
	if err != nil {
		return nil, err
	}
	current := v.(*model.Client)

	if current.Id != client.Id {
		return nil, errors.New("records Id mismatch")
	}
	// keep keys
	client.PrivateKey = current.PrivateKey
	client.PublicKey = current.PublicKey

	client.Updated = time.Now().UTC()

	err = serialize(client.Id, client)
	if err != nil {
		return nil, err
	}

	v, err = deserialize(Id)
	if err != nil {
		return nil, err
	}
	client = v.(*model.Client)

	return client, nil
}

/*
 * DeleteClient from disk
 */
func DeleteClient(id string) error {
	path := filepath.Join(os.Getenv("WG_CONF_DIR"), id)
	err := os.Remove(path)
	if err != nil {
		return err
	}

	// data modified, dump new config
	return generateWgConfig()
}

/*
 * ReadClients all clients
 */
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
			c, err := deserialize(f.Name())
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

/*
 * ReadServer object, create default one
 */
func ReadServer() (*model.Server, error) {
	if !util.FileExists(filepath.Join(os.Getenv("WG_CONF_DIR"), "server.json")) {
		server := &model.Server{}

		key, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		server.PrivateKey = key.String()
		server.PublicKey = key.PublicKey().String()

		presharedKey, err := wgtypes.GenerateKey()
		if err != nil {
			return nil, err
		}
		server.PresharedKey = presharedKey.String()

		server.Name = "Created with default values"
		server.Endpoint = "wireguard.example.com:123"
		server.ListenPort = 51820
		server.Address = "fd9f:6666::10:6:6:1/112, 10.6.6.1/24"
		server.Dns = "fd9f::10:0:0:2, 10.0.0.2"
		server.PersistentKeepalive = 16
		server.Created = time.Now().UTC()
		server.Updated = server.Created

		err = serialize("server.json", server)
		if err != nil {
			return nil, err
		}
	}

	c, err := deserialize("server.json")
	if err != nil {
		return nil, err
	}

	return c.(*model.Server), nil
}

/*
 * UpdateServer keep private values from existing one
 */
func UpdateServer(server *model.Server) (*model.Server, error) {
	current, err := deserialize("server.json")
	if err != nil {
		return nil, err
	}
	server.PrivateKey = current.(*model.Server).PrivateKey
	server.PublicKey = current.(*model.Server).PublicKey
	server.PresharedKey = current.(*model.Server).PresharedKey

	server.Updated = time.Now().UTC()

	err = serialize("server.json", server)
	if err != nil {
		return nil, err
	}

	v, err := deserialize("server.json")
	if err != nil {
		return nil, err
	}
	server = v.(*model.Server)

	return server, nil
}

/*
 * Write object to disk
 */
func serialize(id string, c interface{}) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = util.WriteFile(filepath.Join(os.Getenv("WG_CONF_DIR"), id), b)
	if err != nil {
		return err
	}

	// data modified, dump new config
	return generateWgConfig()
}

/*
 * Read client from disc
 */
func deserializeClient(data []byte) (*model.Client, error) {
	var c *model.Client
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

/*
 * Read server from disc
 */
func deserializeServer(data []byte) (*model.Server, error) {
	var c *model.Server
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
func deserialize(id string) (interface{}, error) {
	path := filepath.Join(os.Getenv("WG_CONF_DIR"), id)

	b, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if id == "server.json" {
		return deserializeServer(b)
	}

	return deserializeClient(b)
}

/*
 * Generate Wireguard interface configuration
 */
func generateWgConfig() error {
	clients, err := ReadClients()
	if err != nil {
		return err
	}

	server, err := ReadServer()
	if err != nil {
		return err
	}

	configDataWg, err := util.DumpServerWg(clients, server)
	if err != nil {
		return err
	}

	err = util.WriteFile(filepath.Join(os.Getenv("WG_CONF_DIR"), os.Getenv("WG_INTERFACE_NAME")), configDataWg.Bytes())
	if err != nil {
		return err
	}

	return nil
}

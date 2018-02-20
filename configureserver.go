package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type jsonConfiguration struct {
	address  string
	port     string
	protocol string
}

type Configuration struct {
	serverConfigs map[string]jsonConfiguration
	misc          map[string]string
}

var (
	PackageConfiguration = Configuration{
		serverConfigs: make(map[string]jsonConfiguration),
		misc:          make(map[string]string)}
)

func GetConfigurationFile(f string) *Configuration {
	PackageConfiguration.loadConfigs(f)
	return &PackageConfiguration
}

func (c *Configuration) loadConfigs(filename string) {
	f, errIO := ioutil.ReadFile(filename)
	if errIO != nil {
		log.Panicf("Can't read configuration file %s. Panicing.", filename)
	}

	var parsedJSON map[string]interface{}
	err := json.Unmarshal(f, &parsedJSON)
	if err != nil {
		log.Panicf("Can't parse json: %s. Panicing", f)
	}

	configs := parsedJSON["serverConfigs"].(map[string]interface{})
	for k, v := range configs {
		vals := v.(map[string]interface{})
		jsonTemp := jsonConfiguration{}
		for i, j := range vals {
			switch i {
			case "address":
				jsonTemp.address = j.(string)
			case "port":
				jsonTemp.port = j.(string)
			case "protocol":
				jsonTemp.protocol = j.(string)
			}
		}
		PackageConfiguration.serverConfigs[k] = jsonTemp
	}

	misc := parsedJSON["misc"].(map[string]interface{})
	for k, v := range misc {
		PackageConfiguration.misc[k] = v.(string)
	}
}

func (c *Configuration) GetListnerDetails(name string) (addr string, protocol string) {
	port := c.serverConfigs[name].port
	addr = fmt.Sprintf(":%s", port)
	protocol = c.serverConfigs[name].protocol
	return
}

func (c *Configuration) GetServerDetails(name string) (addr string, protocol string) {
	ip := c.serverConfigs[name].address
	port := c.serverConfigs[name].port
	addr = fmt.Sprintf("%s:%s", ip, port)
	protocol = c.serverConfigs[name].protocol
	return
}

func (c *Configuration) GetValue(key string) string {
	return c.misc[key]
}

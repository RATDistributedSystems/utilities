package utilities

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type address struct {
	host     string
	port     string
	protocol string
}

type Configuration struct {
	serverConfigs map[string]address
	misc          map[string]string
	docker        bool
}

var (
	conf = Configuration{
		serverConfigs: make(map[string]address),
		misc:          make(map[string]string)}
	potentialAddresses = make(map[string]address)
)

func Load() *Configuration {
	inDocker := os.Getenv("IS_DOCKER")

	if inDocker != "" {
		conf.docker = true
		for _, e := range os.Environ() {
			if position := strings.Index(e, "PARAM_"); position != -1 {
				pair := strings.Split(e[position+6:], "=")
				conf.misc[pair[0]] = pair[1]

				addrConfig := strings.Split(pair[0], "_")
				addPortAddrProto(addrConfig[0], addrConfig[1], pair[1])
			}
		}

	} else {
		conf.docker = false
		f, errIO := os.Open("variable.env")
		if errIO != nil {
			log.Panicf("Can't read configuration file .env. Please place file in current directory")
		}

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if position := strings.Index(line, "PARAM_"); position != -1 {
				params := strings.Split(line[position+6:], "=")
				conf.misc[params[0]] = params[1]

				addrConfig := strings.Split(params[0], "_")
				addPortAddrProto(addrConfig[0], addrConfig[1], params[1])
			}
		}

	}

	for k, v := range potentialAddresses {
		if v.host != "" && v.port != "" && v.protocol != "" {
			conf.serverConfigs[k] = v
		}
	}

	return &conf
}

func addPortAddrProto(servername string, portOrAddress string, param string) {
	if portOrAddress == "ADDR" {
		addr, exists := potentialAddresses[servername]
		if !exists {
			addr = address{}
		}

		addr.host = param
		potentialAddresses[servername] = addr
	}

	// Check if it could be a port name
	if portOrAddress == "PORT" {
		addr, exists := potentialAddresses[servername]
		if !exists {
			addr = address{}
		}

		addr.port = param
		potentialAddresses[servername] = addr
	}

	// Check if it could be a protocol
	if portOrAddress == "PROTOCOL" {
		addr, exists := potentialAddresses[servername]
		if !exists {
			addr = address{}
		}

		addr.protocol = param
		potentialAddresses[servername] = addr
	}
}

func (c Configuration) GetListnerDetails(name string) (addr string, protocol string) {
	r := c.serverConfigs[strings.ToUpper(name)]
	addr = fmt.Sprintf(":%s", r.port)
	protocol = r.protocol
	return
}

func (c Configuration) GetServerDetails(name string) (addr string, protocol string) {
	r := c.serverConfigs[strings.ToUpper(name)]
	addr = fmt.Sprintf("%s:%s", r.host, r.port)
	protocol = r.protocol
	return
}

func (c Configuration) GetValue(key string) string {
	return c.misc[strings.ToUpper(key)]
}

func (c Configuration) Pause() {
	if c.docker {
		sleepTimeInSeconds, err := strconv.Atoi(c.misc["DOCKER_SLEEP"])
		if err != nil {
			log.Printf("Invalid sleep time: %s. Sleeping for a 60 seconds", c.misc["DOCKER_SLEEP"])
			sleepTimeInSeconds = 60
		}
		log.Println("Sleeping for %s seconds...", sleepTimeInSeconds)
		time.Sleep(time.Second * time.Duration(sleepTimeInSeconds))
	}
}

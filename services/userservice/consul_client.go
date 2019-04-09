package userservice

import (
	"os"
	"strings"
	"strconv"
	"log"
	"fmt"

	"github.com/mattermost/viper"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/mattermost/mattermost-server/config"
	"encoding/json"
)

type ApiResponse struct {
	code    int
	message string
	data    interface{}
}

func NewApiResponse() (ApiResponse) {
	return ApiResponse{
		code: 0,
		message: "OK",
	}
}

func NewErrorApiResponse(code int, message string) (ApiResponse) {
	return ApiResponse{
		code: code,
		message: message,
	}
}

func NewApiResponseWithData(data interface{}) (ApiResponse) {
	return ApiResponse{
		code: 0,
		message: "OK",
		data: data,
	}
}

var consulConfig *consulapi.Config

func init() {
	configDSN := viper.GetString("config")
	configStore, err := config.NewStore(configDSN, false)
	if err == nil {
		var configMap map[string]string
		if parseErr := json.Unmarshal([]byte(configStore.Get().ConsulConfigs), &configMap); parseErr != nil {
			log.Fatalln(parseErr)
		}
		consulConfig.Address = configMap["Address"]
		consulConfig.Datacenter = configMap["Datacenter"]
		consulConfig.Scheme = configMap["Scheme"]
	}
}

func registerService(serviceName string) {
	consul, err := consulapi.NewClient(consulConfig)
	if err != nil {
		log.Fatalln(err)
	}

	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = serviceName
	registration.Name = serviceName
	address := hostname()
	registration.Address = address
	p, err := strconv.Atoi(port()[1:len(port())])
	if err != nil {
		log.Fatalln(err)
	}
	registration.Port = p
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", address, p)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	consul.Agent().ServiceRegister(registration)
}

func lookupService(serviceName string) (string, error) {
	consul, err := consulapi.NewClient(consulConfig)
	if err != nil {
		return "", err
	}
	services, err := consul.Agent().Services()
	if err != nil {
		return "", err
	}
	srvc := services[serviceName]
	address := srvc.Address
	port := srvc.Port
	return fmt.Sprintf("http://%s:%v", address, port), nil
}

func port() string {
	p := os.Getenv("SERVICE_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8080"
	}
	return fmt.Sprintf(":%s", p)
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}

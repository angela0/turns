package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pion/turn/v2"
)

type Config struct {
	API    string `json:"api"`
	Port   uint16 `json:"port"`
	Public string `json:"public"`
	Realm  string `json:"realm"`

	Auth map[string]string `json:"auth"`
}

var (
	usersMap  = make(map[string][]byte)
	usersLock sync.RWMutex

	config Config
)

func handleUser(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		usersLock.RLock()
		defer usersLock.RUnlock()

		user := req.URL.Query().Get("user")
		password := usersMap[user]

		rw.Write([]byte(fmt.Sprintf(`{"user": "%s", "password": "%s"}`, user, password)))

		return
	}

	if req.Method != "POST" {
		rw.WriteHeader(404)
		return
	}
	if err := req.ParseForm(); err != nil {
		rw.WriteHeader(400)
		return
	}
	user := req.PostForm.Get("user")
	password := req.PostForm.Get("password")

	if user == "" || password == "" {
		rw.WriteHeader(400)
		return
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	if _, ok := usersMap[user]; ok {
		rw.WriteHeader(401)
	}

	usersMap[user] = turn.GenerateAuthKey(user, config.Realm, password)
}

func main() {
	cfgFile := flag.String("c", "turns.json", "config file")
	flag.Parse()

	c, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Panic(err)
	}
	if err = json.Unmarshal(c, &config); err != nil {
		log.Panic(err)
	}

	udpListener, err := net.ListenPacket("udp4", fmt.Sprintf("0.0.0.0:%d", config.Port))
	if err != nil {
		log.Panicf("Failed to create TURN server listener: %s", err)
	}

	for user, password := range config.Auth {
		usersMap[user] = turn.GenerateAuthKey(user, config.Realm, password)
	}

	s, err := turn.NewServer(turn.ServerConfig{
		Realm: config.Realm,
		// Set AuthHandler callback
		// This is called everytime a user tries to authenticate with the TURN server
		// Return the key for that user, or false when no user is found
		AuthHandler: func(username string, realm string, srcAddr net.Addr) ([]byte, bool) {
			usersLock.RLock()
			defer usersLock.RUnlock()
			if key, ok := usersMap[username]; ok {
				return key, true
			}
			return nil, false
		},
		// PacketConnConfigs is a list of UDP Listeners and the configuration around them
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: udpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP(config.Public), // Claim that we are listening on IP passed by user (This should be your Public IP)
					Address:      "0.0.0.0",                  // But actually be listening on every interface
				},
			},
		},
	})
	if err != nil {
		log.Panic(err)
	}

	defer s.Close()

	if config.API != "" {
		http.HandleFunc("/user", handleUser)
		log.Panic(http.ListenAndServe(config.API, nil))
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

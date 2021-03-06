package nfsmain

import (
	"log"
	"net"

	"github.com/Unknwon/goconfig"
	"github.com/mangalaman93/nfs/voip"
)

var (
	apps   map[string]AppLine
	server *StoppableServer
)

func Start(config *goconfig.ConfigFile) error {
	port, err := config.GetValue("CONTROLLER", "port")
	if err != nil {
		return err
	}
	vl, err := voip.NewVoipLine(config)
	if err != nil {
		return err
	}

	apps = make(map[string]AppLine)
	apps[vl.GetDB()] = vl
	log.Println("[INFO] registered db", vl.GetDB(), "with VoipLine instance")
	server, err = NewStoppableServer(config, apps)
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	tl, _ := l.(*net.TCPListener)
	log.Println("[INFO] listening for data on", l.Addr())
	server.Start(tl)

	return vl.Start()
}

func Stop() {
	// stop the server first so that app methods
	// are not called after calling Stop()
	server.Stop()
	log.Println("[INFO] stopped http server")

	// and now stop apps
	for _, app := range apps {
		app.Stop()
	}
	log.Println("[INFO] stopped all applications")
	log.Println("[INFO] exiting control loop")
}

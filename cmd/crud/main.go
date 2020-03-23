package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/banch0/crud2/pkg/crud/services/agents"
	"gopkg.in/natefinch/lumberjack.v2"

	conf "github.com/banch0/crud2/pkg/config"
	"github.com/banch0/crud2/pkg/crud/services/houses"
	"github.com/banch0/crud2/pkg/crud/services/owners"
	"github.com/banch0/crud2/pkg/mux"
	"github.com/banch0/crud2/pkg/token"
	"github.com/banch0/crud2/pkg/user"

	"github.com/banch0/crud2/cmd/crud/app"
	agent "github.com/banch0/crud2/cmd/crud/app/api/agents"
	house "github.com/banch0/crud2/cmd/crud/app/api/houses"
	api "github.com/banch0/crud2/cmd/crud/app/api/owners"

	"github.com/banch0/crud2/pkg/jwt"

	"github.com/banch0/crud2/pkg/di"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	dsn = flag.String("dsn", "", "Postgres DSN")
)

// DSN ...
type DSN string

// ErrCreatePool ...
var ErrCreatePool = errors.New("can't create pool")

func main() {
	config := &conf.Config{}

	log.SetOutput(&lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     60,
	})
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("App config file: %v\n", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("App config file: %v\n", err)
	}

	flag.Parse()

	*dsn = `postgres://` + config.DBUser + ":"
	*dsn += config.DBPass + "@" + config.DBHost + ":"
	*dsn += config.DBPort + "/" + config.DBName

	addr := net.JoinHostPort(config.Host, config.Port)
	secret := jwt.Secret("secret")

	start(addr, *dsn, secret)
}

func start(addr string, dsn string, secret jwt.Secret) {
	container := di.NewContainer()
	container.Provide(
		mux.NewExactMux,
		app.NewServer,
		func() jwt.Secret { return secret },
		func() DSN { return DSN(dsn) },
		NewConnectionPool,
		token.NewService,
		user.NewService,
		agents.NewService,
		owners.NewService,
		houses.NewService,
		api.NewMainServer,
		agent.NewMainServer,
		house.NewHouseServer,
	)

	var appService *app.Server
	container.Component(&appService)
	appService.Start()

	log.Println("Server starting ...")
	panic(http.ListenAndServe(addr, appService))
}

// NewConnectionPool ...
func NewConnectionPool() *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), string(*dsn))
	if err != nil {
		panic(ErrCreatePool)
	}
	return pool
}

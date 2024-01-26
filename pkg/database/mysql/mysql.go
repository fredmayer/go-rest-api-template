package mysql

import (
	"context"
	"database/sql"
	"github.com/fredmayer/go-rest-api-template/pkg/logging"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Options struct {
	Host     string
	Port     string
	User     string
	Password string
	Db       string
}

var (
	config          Options
	Conn            *sql.DB
	cancelKeepAlive context.CancelFunc
)

const (
	MAX_OPEN_CONNECTIONS = 20
	MAX_IDLE_CONNEXTIONS = 10
	MAX_LIFE             = time.Minute * 3
)

func NewClient(ctx context.Context, opt Options) *sql.DB {
	config = opt
	cs := getCS()
	logging.Get().Debugf("Mysql config:%v \n", cs)

	//dial
	var err error
	Conn, err = dial(ctx, cs)
	if err != nil {
		logging.Get().Panicf("MySQL connection error: %v \n", err)
	}

	//TODO go keepAlive
	var ctxWC context.Context
	ctxWC, cancelKeepAlive = context.WithCancel(ctx)
	go keepAlive(ctxWC)

	return Conn
}

func Stop(ctx context.Context) {
	cancelKeepAlive() //Завершаем контекст keepAlive
	err := Conn.Close()

	if err != nil {
		logging.Get().Panicf("Error close Mysql connection: %v \n", err)
	}
}

func keepAlive(ctx context.Context) {
	for {
		time.Sleep(time.Second * 3)

		select {
		case <-ctx.Done():
			_ = Conn.Close()
			return
		default:
			ctxPing, _ := context.WithTimeout(ctx, time.Second*2)

			err := Conn.PingContext(ctxPing)
			if err == nil {
				continue
			}

			logging.Log().Errorf("Connection to MySQL lost: %v \n", err)

			c, err := dial(ctx, getCS())
			if err != nil {
				logging.Log().Errorf("Restoring... %v \n", err)
				continue
			}

			Conn = c
			logging.Log().Infof("Clickhouse reconnected\n")
		}
	}
}

func getCS() string {
	return config.User + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Db
}

func dial(ctx context.Context, cs string) (*sql.DB, error) {
	db, err := sql.Open("mysql", cs)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(MAX_LIFE)
	db.SetMaxOpenConns(MAX_OPEN_CONNECTIONS)
	db.SetMaxIdleConns(MAX_IDLE_CONNEXTIONS)

	return db, nil
}

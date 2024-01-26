package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/fredmayer/go-rest-api-template/pkg/logging"
	"log"
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
	config Options
	Conn   driver.Conn
)

func (o *Options) getUrl() []string {
	return []string{fmt.Sprintf("%s:%s", o.Host, o.Port)}
}

func NewClient(ctx context.Context, cfg Options) driver.Conn {
	config = cfg
	fmt.Println(config)
	c, err := dial(ctx)
	if err != nil {
		log.Panic(err)
	}
	Conn = c

	go keepAlive(ctx)

	return Conn
}

func dial(ctx context.Context) (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: config.getUrl(),
		Auth: clickhouse.Auth{
			Database: config.Db,
			Username: config.User,
			Password: config.Password,
		},
		// DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
		// 	dialCount++
		// 	var d net.Dialer
		// 	return d.DialContext(ctx, "tcp", addr)
		// },
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		// Protocol: clickhouse.HTTP,
	})

	if err != nil {
		return nil, err
	}
	logging.Log().Println("Connected...")

	v, err := conn.ServerVersion()
	if err != nil {
		return nil, err
	}
	logging.Log().Debug(v)

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func keepAlive(ctx context.Context) {
	for {
		time.Sleep(time.Second * 3)

		select {
		case <-ctx.Done():
			_ = Conn.Close()
			return
		default:
			ctxPing, _ := context.WithTimeout(ctx, time.Second*3)

			err := Conn.Ping(ctxPing)
			if err == nil {
				continue
			}

			logging.Log().Errorf("Connection to clickhouse lost: %v \n", err)

			c, err := dial(ctx)
			if err != nil {
				logging.Log().Errorf("Restoring... %v \n", err)
				continue
			}

			Conn = c
			logging.Log().Infof("Clickhouse reconnected\n")
		}
	}
}

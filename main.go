package main

import (
	"net"
	"net/http"
	"os"
	pb "wordfilter/proto"

	cli "gopkg.in/urfave/cli.v2"

	log "github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	go func() {
		log.Info(http.ListenAndServe("0.0.0.0:6060", nil))
	}()
	app := &cli.App{
		Name: "wordfilter",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "listen",
				Value: ":10000",
				Usage: "listening address:port",
			},
			&cli.StringFlag{
				Name:  "dictionary",
				Value: "./dictionary.txt",
				Usage: "path for dictionary",
			},
			&cli.StringFlag{
				Name:  "dirty",
				Value: "./dirty.txt",
				Usage: "path for dirty words",
			},
			&cli.StringFlag{
				Name:  "replace-word",
				Value: "*",
				Usage: "replace-word for sensitive data",
			},
		},

		Action: func(c *cli.Context) error {
			log.Println("listen:", c.String("listen"))
			log.Println("replace-word:", c.String("replace-word"))
			log.Println("dictionary:", c.String("dictionary"))
			log.Println("dirty:", c.String("dirty"))
			// 监听
			lis, err := net.Listen("tcp", c.String("listen"))
			log.Println("listen:", c.String("listen"))
			if err != nil {
				log.Panic(err)
				os.Exit(-1)
			}
			log.Info("listening on ", lis.Addr())

			// 注册服务
			s := grpc.NewServer()
			ins := &server{}
			ins.init(c)
			pb.RegisterWordFilterServiceServer(s, ins)

			// 开始服务
			return s.Serve(lis)
		},
	}
	app.Run(os.Args)
}

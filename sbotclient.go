package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"go.cryptoscope.co/muxrpc"
	"go.cryptoscope.co/netwrap"
	"go.cryptoscope.co/secretstream"
	"go.cryptoscope.co/ssb"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// @@@ move this to another package within the go.crytoscope codebase
type noopHandler struct{}

func (h noopHandler) HandleConnect(ctx context.Context, edp muxrpc.Endpoint) {
	srv := edp.(muxrpc.Server)
	fmt.Println("event", "onConnect", "addr", srv.Remote())
}

func (h noopHandler) HandleCall(ctx context.Context, req *muxrpc.Request, edp muxrpc.Endpoint) {
	fmt.Println("event", "onCall", "args", fmt.Sprintf("%v", req.Args), "method", req.Method)
}

func initClient(pathToKeyfile string) (client muxrpc.Endpoint, err error) {
	// NB: This is hardcoded b/c I deleted the "connect to a
	// remote" section of the code.
	const localhost = "localhost:8008"
	sbotAppKey, err := base64.StdEncoding.DecodeString("1KHLiKZvAvjbY1ziZEHMXawbCEIM6qwjCDm3VYRan/s=")
	if err != nil {
		return nil, err
	}

	localKey, err := ssb.LoadKeyPair(pathToKeyfile)
	if err != nil {
		return nil, err
	}

	c, err := secretstream.NewClient(localKey.Pair, sbotAppKey)
	if err != nil {
		return nil, errors.Wrap(err, "error creating secretstream.Client")
	}

	var remotPubKey = localKey.Pair.Public
	plainAddr, err := net.ResolveTCPAddr("tcp", localhost)
	if err != nil {
		// @@@ remove / reword
		return nil, errors.Wrapf(err, "init: base64 decode of --remoteKey failed")
	}

	conn, err := netwrap.Dial(plainAddr, c.ConnWrapper(remotPubKey))
	if err != nil {
		return nil, errors.Wrap(err, "error dialing")
	}
	/* coming soon:
	conn, err := net.Dial("unix", "/home/cryptix/.ssb/socket")
	if err != nil {
		return nil,  errors.Wrap(err, "error dialing unix sock")
	}
	*/
	var rwc io.ReadWriteCloser = conn
	// logs every muxrpc packet
	// if ctx.Bool("verbose") {
	// 	rwc = debug.Wrap(log, rwc)
	// }
	pkr := muxrpc.NewPacker(rwc)

	h := noopHandler{}
	client = muxrpc.HandleWithRemote(pkr, &h, conn.RemoteAddr())

	longctx := context.Background()
	longctx, shutdownFunc := context.WithCancel(longctx)
	signalc := make(chan os.Signal)
	signal.Notify(signalc, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalc
		fmt.Println("killed. shutting down")
		shutdownFunc()
		time.Sleep(1 * time.Second)
		check(pkr.Close())
		os.Exit(0)
	}()
	// logging.SetCloseChan(signalc)
	go func() {
		err := client.(muxrpc.Server).Serve(longctx)
		check(err)
		// if this returns, you can't return anything.
		// Maybe this should cancel the context??
	}()
	log.Println("init", "done")
	return client, nil
}
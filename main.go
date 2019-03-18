package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/go-kit/kit/log/term"
	"github.com/QOSGroup/qos/app"
	"github.com/tendermint/tendermint/libs/log"
	tmrpc "github.com/tendermint/tendermint/rpc/client"
	clictx "github.com/QOSGroup/qbase/client/context"
	"strings"
	"time"
	"os/signal"
	"syscall"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"os/user"
)

var logger = log.NewNopLogger()

// FIXME config slice
type Config struct {
	Name	string	`json:"name"`
	Addr  	string 	`json:"address"`
	Pass 	string	`json:"password"`
}

func main() {
	var durationInt, txsRate, connections int
	var verbose bool
	var qosPath, configFile, outputFormat, broadcastTxMethod string

	flagSet := flag.NewFlagSet("qos-bench", flag.ExitOnError)
	flagSet.IntVar(&connections, "c", 1, "Connections to keep open per endpoint")
	flagSet.IntVar(&durationInt, "T", 30, "Exit after the specified amount of time in seconds")
	flagSet.IntVar(&txsRate, "R", 100, "Txs per second to send in a connection")
	flagSet.StringVar(&qosPath, "home", "~/.qoscli", "Setup qos home path")
	flagSet.StringVar(&configFile, "config", "./config.json", "File deployed from which account to prepare Test Transactions")
	flagSet.StringVar(&outputFormat, "output-format", "plain", "Output format: plain or json")
	flagSet.StringVar(&broadcastTxMethod, "broadcast-tx-method", "async", "Broadcast method: async (no guarantees; fastest), sync (ensures tx is checked) or commit (ensures tx is checked and committed; slowest)")
	flagSet.BoolVar(&verbose, "v", false, "Verbose output")
	flagSet.Usage = func() {
		fmt.Printf(`
QOS blockchain benchmarking tool.

Usage:
	qos-bench [-c 1] [-T 10] [-r 1000] [endpoints] [-output-format <plain|json> [-broadcast-tx-method <async|sync|commit>]]

Examples:
	qos-bench -v -T 10 -r 10 -output-format plain -broadcast-tx-method async localhost:26657
`)
		fmt.Println("Flags:")

		flagSet.PrintDefaults()
	}

	// Parse endpoint.
	flagSet.Parse(os.Args[1:])
	if flagSet.NArg() == 0 {
		flagSet.Usage()
		os.Exit(1)
	}

	// Enable verbose module.
	if verbose {
		if outputFormat == "json" {
			printErrorAndExit("Verbose mode not supported with json output.")
		}
		// Color errors red
		colorFn := func(keyvals ...interface{}) term.FgBgColor {
			for i := 1; i < len(keyvals); i += 2 {
				if _, ok := keyvals[i].(error); ok {
					return term.FgBgColor{Fg: term.White, Bg: term.Red}
				}
			}
			return term.FgBgColor{}
		}
		logger = log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), colorFn)
		fmt.Printf("Running %ds test @ %s\n", durationInt, flagSet.Arg(0))
	}

	// Check broadcast method.
	if broadcastTxMethod != "async" &&
		broadcastTxMethod != "sync" &&
		broadcastTxMethod != "commit" {
		printErrorAndExit("broadcast-tx-method should be either 'sync', 'async' or 'commit'.")
	}

	// Load config file.
	config, err := Load(configFile)
	if err != nil {
		printErrorAndExit(err.Error())
	}

	// Parse and find qos home directory.
	path, err := ParsePath(qosPath)
	if err != nil {
		printErrorAndExit(err.Error())
	}

	// Init values.
	endpoints     := strings.Split(flagSet.Arg(0), ",")
	client        := tmrpc.NewHTTP(endpoints[0], "/websocket")

	// Log out test parameter.
	fmt.Println("time duration: ", durationInt)
	fmt.Println("transacter rate: ", txsRate)
	fmt.Println("transacter broadcast method: ", broadcastTxMethod)

	// Prepare qos transactions, this step takes some times.
	transacters := PrepareTransacters(
		config,
		path,

		client,
		endpoints,
		connections,
		durationInt,
		txsRate,
		"broadcast_tx_"+broadcastTxMethod,
	)

	// Time duration.
	timeStart := time.Now()
	logger.Info("Time last transacter started", "t", timeStart)
	duration := time.Duration(durationInt) * time.Second
	timeEnd := timeStart.Add(duration)
	logger.Info("End time for calculation", "t", timeEnd)

	// anchor the last block
	initialHeight := LatestBlockHeight(client)
	logger.Info("Latest block height", "h", initialHeight)

	// Start broadcasting tx.
	for _, t := range transacters {
		t.Start()
	}

	// Quit when interrupted or received SIGTERM.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			fmt.Printf("captured %v, exiting...\n", sig)
			for _, t := range transacters {
				t.Stop()
			}
			os.Exit(1)
		}
	}()

	// Wait until transacters have begun until we get the start time.
	<-time.After(duration)
	for i, t := range transacters {
		t.Stop()
		numCrashes :=
			countCrashes(t.connsBroken)
		if numCrashes != 0 {
			fmt.Printf("%d connections crashed on transacter #%d\n", numCrashes, i)
		}
	}

	logger.Debug("Time all transacters stopped", "t", time.Now())

	// State txs from initial height to current height.
	stats, err := calculateStatistics(
		client,
		initialHeight,
		timeStart,
		timeEnd,
		durationInt,
	)
	if err != nil {
		printErrorAndExit(err.Error())
	}

	// Print it in format.
	printStatistics(stats, outputFormat)
}

func Load(filename string) (Config, error) {
	var config Config

	// ReadFile func read all data from filename, and return as []byte type.
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	// Read in json format，decode into data struct
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, err
}

func ParsePath(qosPath string) (string, error) {
	var path string
	if !filepath.IsAbs(qosPath) {
		switch qosPath {
		case "~/.qoscli":
			user, err := user.Current()
			if nil == err {
				path = filepath.Join(user.HomeDir, ".qoscli")
			}
		default:
			wd, _ := os.Getwd()
			path, _ = filepath.Abs(filepath.Join(wd, qosPath))
		}
	}

	// Check qosPath
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		printErrorAndExit(err.Error())
	}
	return path, err
}

func LatestBlockHeight(client tmrpc.Client) int64 {

	status, err := client.Status()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return status.SyncInfo.LatestBlockHeight
}

func countCrashes(crashes []bool) int {
	count := 0
	for i := 0; i < len(crashes); i++ {
		if crashes[i] {
			count++
		}
	}
	return count
}

func PrepareTransacters(
	config Config,
	qosPath string,
	client tmrpc.Client,
	endpoints []string,
	connections,
	durationInt int,
	txsRate int,
	broadcastTxMethod string,
) []*transacter {
	fmt.Println("Start Preparing Test Transactions ...")
	transacters := make([]*transacter, len(endpoints))
	ctx := clictx.NewCLIContext().WithCodec(app.MakeCodec()).WithClient(client)

	for i, e := range endpoints {
		t := newTransacter(config, qosPath, ctx, e, connections, durationInt, txsRate, broadcastTxMethod)
		t.SetLogger(logger)
		t.PrepareTx()
		transacters[i] = t
	}
	fmt.Println("Test Transactions All Ready !!!")
	return transacters
}

func printErrorAndExit(err string) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

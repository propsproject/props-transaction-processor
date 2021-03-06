package main

import (
	"encoding/json"
	"fmt"
	"github.com/propsproject/props-transaction-processor/core"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
)

var atom = zap.NewAtomicLevel()
var encoderCfg = zap.NewProductionEncoderConfig()
var logger *zap.SugaredLogger

func main() {
	pflag.StringP("verbose", "v", "debug", "Log verbosity info|warning|debug")
	pflag.String("token",  "", "PROPS token contract address")
	pflag.IntP("worker-queue", "q", 100, "Set the maximum queue size before rejecting process requests")
	pflag.IntP("worker-threads", "t", 0, "Set the number of worker threads to use for processing requests in parallel")
	pflag.StringP("config-file-path", "f", "", "Path to configuration file. Other arguments ignored if this flag is set")
	pflag.BoolP("config-file", "c", false, "If flag is set configurations will be loaded from ConfigFilePath")
	pflag.Parse()

	viper.BindPFlag( "use-config", pflag.Lookup("config-file"))
	if viper.GetBool("use-config") {
		viper.BindPFlag( "config-file-path", pflag.Lookup("config-file-path"))
		err := parseConfigFile()
		if err != nil {
			logger.Errorf("error parsing configuration file:  ", err)
			os.Exit(1)
		}
	} else {
		viper.BindPFlags(pflag.CommandLine)
	}

	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	atom.SetLevel(zap.DebugLevel)

	var tmpLogger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	tmpLogger = tmpLogger.With(
		zap.String("app", viper.GetString("app")),
		zap.String("name", viper.GetString("name")),
		zap.String("env", viper.GetString("environment")),
	)

	defer tmpLogger.Sync()

	logger = tmpLogger.Sugar()

	logger.Info("Starting the transaction processor")

	// Set some default values in the logger
	//logger.SetDefaultKeyValues(zap.String("app", viper.GetString("app")), zap.String("name", viper.GetString("name")), zap.String("env", viper.GetString("environment")))

	tp := core.NewTransactionProcessor(viper.GetString("validator_url"))
	err := tp.Start()
	if err != nil {
		logger.Errorf("Processor stopped: ", err)
	}
}

func parseConfigFile() error {
	config := viper.GetString("config-file-path")
	if config == "" {
		return fmt.Errorf("illegal argument for config file path, path must be specified")
	}

	abs, err := filepath.Abs(config)
	if err != nil {
		return fmt.Errorf("error reading filepath: (%s)", err)
	}

	// get the config name
	base := filepath.Base(abs)

	// get the path
	path := filepath.Dir(abs)
	viper.SetConfigType("json")
	viper.SetConfigName(strings.Split(base, ".")[0])
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Find and read the config file; Handle errors reading the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading configuration file: (%s)", err)
	}
	// parse json values back into viper
	var settlementAddressJson map[string]interface{}
	var validSignersJson map[string]interface{}
	err1 := json.Unmarshal([]byte(viper.GetString("settlement_from_addresses")),&settlementAddressJson)
	if err1 != nil {
		return fmt.Errorf("error reading configuration file settlement_from_addresses malformed: (%s)", err1)
	}
	viper.Set("settlement_from_addresses_map", settlementAddressJson)
	err2 := json.Unmarshal([]byte(viper.GetString("valid_signers_addresses")),&validSignersJson)
	if err2 != nil {
		return fmt.Errorf("error reading configuration file settlement_from_addresses malformed: (%s)", err2)
	}
	viper.Set("valid_signers_addresses_map", validSignersJson)


	return nil
}

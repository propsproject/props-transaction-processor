package state

import (
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

type State struct {
	context *processor.Context
}

var atom = zap.NewAtomicLevel()
var encoderCfg = zap.NewProductionEncoderConfig()
var logger *zap.SugaredLogger

var doOnce sync.Once

func NewState(context *processor.Context) *State {
	initLogger()

	return &State{context: context}
}

func initLogger() {
	doOnce.Do(func() {
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
	})
}
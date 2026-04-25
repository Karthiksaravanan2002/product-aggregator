package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
	Log, _ = zap.NewProduction()
}

func Sync() {
	if Log != nil {
		Log.Sync()
	}
}

var (
	String = zap.String
	Int    = zap.Int
	Error  = zap.Error
)

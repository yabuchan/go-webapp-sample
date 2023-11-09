package logger

import (
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func build(cfg *Config) (*zap.Logger, error) {
	var zapCfg = cfg.ZapConfig
	enc, _ := newEncoder(zapCfg)
	writer, errWriter := openWriters(cfg)

	atom := zap.NewAtomicLevel()
	//zapCfg.Level = atom
	log := zap.New(zapcore.NewCore(enc, writer, &atom), buildOptions(zapCfg, errWriter)...)
	go 	logLevelWatcher(&atom, log)
	return log, nil
}

func logLevelWatcher(atom *zap.AtomicLevel, log *zap.Logger) {
	var current, past string
	for {
		time.Sleep(100*time.Millisecond)
		current = readLog(log)
		if current != past {
			past = current
			na, err := zapcore.ParseLevel(current)
			if err != nil {
				continue
			}
			atom.SetLevel(na)
			log.Info("log level is changed", zap.String("LEVEL", log.Level().CapitalString()))
		}
	}
}

var mx sync.Mutex
func readLog(log *zap.Logger) string {
	mx.Lock()
	defer mx.Unlock()
	dat, err := os.ReadFile("/logLevel")
	if err != nil {
		return ""
	}
	lines := strings.TrimSuffix(string(dat), "\n")
	return lines
}

func newEncoder(cfg zap.Config) (zapcore.Encoder, error) {
	switch cfg.Encoding {
	case "console":
		return zapcore.NewConsoleEncoder(cfg.EncoderConfig), nil
	case "json":
		return zapcore.NewJSONEncoder(cfg.EncoderConfig), nil
	}
	return nil, errors.New("failed to set encoder")
}

func openWriters(cfg *Config) (zapcore.WriteSyncer, zapcore.WriteSyncer) {
	writer := open(cfg.ZapConfig.OutputPaths, &cfg.LogRotate)
	errWriter := open(cfg.ZapConfig.ErrorOutputPaths, &cfg.LogRotate)
	return writer, errWriter
}

func open(paths []string, rotateCfg *lumberjack.Logger) zapcore.WriteSyncer {
	writers := make([]zapcore.WriteSyncer, 0, len(paths))
	for _, path := range paths {
		writer := newWriter(path, rotateCfg)
		writers = append(writers, writer)
	}
	writer := zap.CombineWriteSyncers(writers...)
	return writer
}

func newWriter(path string, rotateCfg *lumberjack.Logger) zapcore.WriteSyncer {
	switch path {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	}
	sink := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   path,
			MaxSize:    rotateCfg.MaxSize,
			MaxBackups: rotateCfg.MaxBackups,
			MaxAge:     rotateCfg.MaxAge,
			Compress:   rotateCfg.Compress,
		},
	)
	return sink
}

func buildOptions(cfg zap.Config, errWriter zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errWriter)}
	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if cfg.Development {
		stackLevel = zap.WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}
	return opts
}

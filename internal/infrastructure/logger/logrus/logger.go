package logrus

import (
	"chat/internal/infrastructure/logger"
	"github.com/sirupsen/logrus"
)

func NewLogger(lvl logger.Level) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return logger, nil
}

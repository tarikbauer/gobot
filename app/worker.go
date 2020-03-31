package app

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/tarikbauer/gobot/application"
)

func RunWorker(logger logrus.Logger, service application.Bot, c chan <- error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "logger", logger)
	c <- service.Run(ctx)
}

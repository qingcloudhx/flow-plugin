//Created by zhbinary on 2019-04-10.
//Email: zhbinary@gmail.com
package utils

import "log"
import "go.uber.org/zap"

var Log = NewZap()

func NewZap() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	return logger.Sugar()
}

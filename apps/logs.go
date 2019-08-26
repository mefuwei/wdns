/*
 Auther : F.W
 Create time  2018/12/29
*/
package apps

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {

	logger.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	// logger.SetReportCaller(true)
	logger.Out = os.Stdout
	logger.SetLevel(logrus.DebugLevel)

}

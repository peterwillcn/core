package log

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/bshuster-repo/logrus-logstash-hook"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

var (
	root = &logger{}
)

// New returns a new logger with the given key/value .
// New is a convenient alias for Root().New
func New(key string, value interface{}) Logger {
	if root == nil {
		fmt.Println("no root")
	}
	return root.New(key, value)
}

// Init setups default field and add hook
func Init(id []byte) {
	appSession := os.Getenv("APPSESSION")
	nodId := byteTohex(id)
	logIp := os.Getenv("LOGIP")
	clientIP := os.Getenv("NODEIP")
	logrus.SetOutput(ioutil.Discard)

	if logIp != "" {
		logrus.SetLevel(logrus.DebugLevel)

		hook, err := logrustash.NewHook("tcp", logIp, appSession)
		if err != nil {
			fmt.Println(err)
			logrus.Error(err)
		}
		logrus.AddHook(hook)
	}

	path := "./vault/doslog.txt"
	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)
	if err != nil {
		fmt.Println("rotatelogs.New err", err)
	}

	logrus.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
		},
		&logrus.JSONFormatter{},
	))
	if appSession != "" {
		root.entry = logrus.WithFields(logrus.Fields{
			"appSession": appSession,
			"clientip":   clientIP,
			"nodeID":     nodId,
		})
	} else {
		root.entry = logrus.WithFields(logrus.Fields{
			"clientip": clientIP,
			"nodeID":   nodId,
		})
	}
}

func byteTohex(a []byte) string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "" + string(result)
}

// AddField is a convenient alias for Root().AddField
func AddField(key string, value interface{}) {
	root.AddField(key, value)
}

// Debug is a convenient alias for Root().Debug
func Debug(msg string) {
	root.Debug(msg)
}

// Info is a convenient alias for Root().Info
func Info(msg string) {
	root.Info(msg)
}

// Warn is a convenient alias for Root().Warn
func Warn(msg string) {
	root.Warn(msg)
}

// Error is a convenient alias for Root().Error
func Error(err error) {
	if root != nil {
		root.Error(err)
	}
}

// Fatal is a convenient alias for Root().Fatal
func Fatal(err error) {
	root.Fatal(err)
}

// Event is a convenient alias for Root().Event
func Event(e string, f map[string]interface{}) {
	root.Event(e, f)
}

// TimeTrack is a convenient alias for Root().TimeTrack
func TimeTrack(start time.Time, e string, info map[string]interface{}) {
	root.TimeTrack(start, e, info)
}

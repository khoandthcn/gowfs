package gowfs

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "gowfs")

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{})
}

const WebHdfsVer string = "/webhdfs/v1"

type Configuration struct {
	Addr                  string // host:port
	BasePath              string // initial base path to be appended
	User                  string // user.name to use to connect
	Passwd                string // user.passwd for simple authen
	ConnectionTimeout     time.Duration
	DisableKeepAlives     bool
	DisableCompression    bool
	ResponseHeaderTimeout time.Duration
	UseHTTPS              bool
	BasicAuth             bool
}

func NewConfiguration() *Configuration {
	return &Configuration{
		ConnectionTimeout:     time.Second * 17,
		DisableKeepAlives:     false,
		DisableCompression:    true,
		ResponseHeaderTimeout: time.Second * 17,
		UseHTTPS:              true,
		BasicAuth:             true,
	}
}

func (conf *Configuration) GetNameNodeUrl() (*url.URL, error) {
	if &conf.Addr == nil {
		return nil, errors.New("Configuration namenode address not set.")
	}

	scheme := "http"
	if conf.UseHTTPS {
		scheme = "https"
	}

	urlStr := fmt.Sprintf("%s://%s%s%s", scheme, conf.Addr, WebHdfsVer, conf.BasePath)

	if &conf.User == nil || len(conf.User) == 0 {
		u, _ := user.Current()
		conf.User = u.Username
	}
	urlStr = urlStr + "?user.name=" + conf.User

	u, err := url.Parse(urlStr)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (conf *Configuration) GetBasicAuthString() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", conf.User, conf.Passwd))))
}

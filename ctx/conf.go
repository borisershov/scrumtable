package ctx

import (
	conf "github.com/nixys/nxs-go-conf"
)

type confOpts struct {
	LogFile  string       `conf:"logfile" conf_extraopts:"default=stdout"`
	LogLevel string       `conf:"loglevel" conf_extraopts:"default=info"`
	MySQL    mysqlConf    `conf:"mysql" conf_extraopts:"required"`
	Telegram telegramConf `conf:"telegram" conf_extraopts:"required"`
}

type mysqlConf struct {
	Host     string `conf:"host" conf_extraopts:"required"`
	DB       string `conf:"db" conf_extraopts:"required"`
	User     string `conf:"user" conf_extraopts:"required"`
	Password string `conf:"password" conf_extraopts:"required"`
}

type telegramConf struct {
	BotAPI    string `conf:"bot_api" conf_extraopts:"required"`
	RedisHost string `conf:"redis" conf_extraopts:"required"`
}

type telegramProxyConf struct {
	Type     string `conf:"type"`
	Host     string `conf:"host"`
	Login    string `conf:"login"`
	Password string `conf:"password"`
}

func confRead(confPath string) (confOpts, error) {

	var c confOpts

	err := conf.Load(&c, conf.Settings{
		ConfPath:    confPath,
		ConfType:    conf.ConfigTypeYAML,
		UnknownDeny: true,
	})
	if err != nil {
		return c, err
	}

	return c, err
}

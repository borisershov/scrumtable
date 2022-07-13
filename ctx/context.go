package ctx

import (
	"github.com/nixys/scrumtable/ds/mysql"
	tgbot "github.com/nixys/scrumtable/modules/bot"

	appctx "github.com/nixys/nxs-go-appctx/v2"
)

// Ctx defines application custom context
type Ctx struct {
	Conf  confOpts
	MySQL mysql.MySQL
	Bot   tgbot.Bot
}

// Init initiates application custom context
func (c *Ctx) Init(opts appctx.CustomContextFuncOpts) (appctx.CfgData, error) {

	// Read config file
	conf, err := confRead(opts.Config)
	if err != nil {
		return appctx.CfgData{}, err
	}

	// Set application context
	c.Conf = conf

	// Connect to MySQL
	c.MySQL, err = mysql.Connect(mysql.Param{
		Host:     c.Conf.MySQL.Host,
		Database: c.Conf.MySQL.DB,
		User:     c.Conf.MySQL.User,
		Password: c.Conf.MySQL.Password,
	})
	if err != nil {
		return appctx.CfgData{}, err
	}

	c.Bot, err = tgbot.Init(tgbot.Settings{
		MySQL:     c.MySQL,
		APIToken:  c.Conf.Telegram.BotAPI,
		RedisHost: c.Conf.Telegram.RedisHost,
	})
	if err != nil {
		return appctx.CfgData{}, err
	}

	return appctx.CfgData{
		LogFile:  c.Conf.LogFile,
		LogLevel: c.Conf.LogLevel,
	}, nil
}

// Reload reloads application custom context
func (c *Ctx) Reload(opts appctx.CustomContextFuncOpts) (appctx.CfgData, error) {

	opts.Log.Debug("reloading context")

	c.MySQL.Close()

	return c.Init(opts)
}

// Free frees application custom context
func (c *Ctx) Free(opts appctx.CustomContextFuncOpts) int {

	opts.Log.Debug("freeing context")

	c.MySQL.Close()

	return 0
}

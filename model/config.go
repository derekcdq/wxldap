package model

import (
	"github.com/akkuman/parseConfig"
	"github.com/wonderivan/logger"
	"os"
)

type ConfigFile struct {
	FileName string
	FileDir string
	FilePath string
}

type WxConfig struct {
	ApiUrl string
	CorpID string
	CorpSecret string
	Interval int
}

type LdapConfig struct {
	LdapHost string
	BaseDn string
	RootDn string
	Password string
	DefaultPassword string
}

type SmtpConfig struct {
	MailHost string
	MailUser string
	Password string
}

type CallbackConfig struct {
	Token           string
	ReceiverId      string
	EncodingAeskey  string

}

var (
	wxConfig  *WxConfig
	ldapConfig *LdapConfig
	smtpConfig *SmtpConfig
	callbackConfig *CallbackConfig
)

func (t *ConfigFile) Init (fileName string , fileDir string) *ConfigFile {
	filePath := "./" + fileDir + "/" + fileName
	_, err := os.Stat(filePath)
	if err != nil {
		filePath = "../" + fileDir + "/" + fileName
		_, err = os.Stat(filePath)
		if err != nil {
			logger.Info("配置文件错误:", err)
			os.Exit(0)
		}
	}
	t.FileName = fileName
	t.FileDir  = fileDir
	t.FilePath = filePath
	return t
}

func (t *WxConfig) Init () interface{} {
	configFile := new (ConfigFile)
	configFile.Init ("wechat.json","conf")
	config := parseConfig.New(configFile.FilePath)
	t.ApiUrl = config.Get("wechatConfig > apiUrl").(string)
	t.CorpID = config.Get("wechatConfig > corpId").(string)
	t.CorpSecret = config.Get("wechatConfig > corpSecret").(string)
	return t
}

func ( t *LdapConfig ) Init () interface{} {
	configFile := new (ConfigFile)
	configFile.Init ("ldap.json","conf")
	config := parseConfig.New(configFile.FilePath)
	t.LdapHost = config.Get("ldapConfig > ldapHost").(string)
	t.BaseDn = config.Get("ldapConfig > baseDn").(string)
	t.RootDn = config.Get("ldapConfig > rootDn").(string)
	t.Password = config.Get("ldapConfig > password").(string)
	t.DefaultPassword = config.Get("ldapConfig > defaultPassword").(string)
	return t
}

func ( t *SmtpConfig ) Init () interface{} {
	configFile := new (ConfigFile)
	configFile.Init ("smtp.json","conf")
	config := parseConfig.New(configFile.FilePath)
	t.MailHost = config.Get("smtpConfig > mailHost").(string)
	t.MailUser = config.Get("smtpConfig > mailUser").(string)
	t.Password = config.Get("smtpConfig > mailPasswd").(string)
	return t
}

func ( t *CallbackConfig ) Init () interface{} {
	configFile := new (ConfigFile)
	configFile.Init ("callback.json","conf")
	config := parseConfig.New(configFile.FilePath)
	t.Token = config.Get("callbackConfig > token").(string)
	t.ReceiverId = config.Get("callbackConfig > receiverId").(string)
	t.EncodingAeskey = config.Get("callbackConfig > encodingAeskey").(string)
	return t
}


func InitConfig () {
	wxConfig = new(WxConfig)
	wxConfig.Init()
	ldapConfig = new(LdapConfig)
	ldapConfig.Init()
	smtpConfig = new(SmtpConfig)
	smtpConfig.Init()
	callbackConfig = new(CallbackConfig)
	callbackConfig.Init()

}






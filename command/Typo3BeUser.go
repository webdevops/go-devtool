package command

import (
	"fmt"
	"github.com/tredoe/osutil/user/crypt/md5_crypt"
)

const (
	t3BeUserConfig = `
options.clearCache.system = 1
options.clearCache.all = 1
options.enableShowPalettes = 1
options.pageTree.showPageIdWithTitle = 1
options.pageTree.showPathAboveMounts = 1
options.pageTree.showDomainNameWithTitle = 1
admPanel.enable.edit = 1
admPanel.module.edit.forceDisplayFieldIcons = 1
admPanel.hide = 0
setup.default.thumbnailsByDefault = 1
setup.default.enableFlashUploader = 0
setup.default.recursiveDelete = 1
setup.default.showHiddenFilesAndFolders = 1
setup.default.resizeTextareas_Flexible = 1
setup.default.copyLevels = 99
setup.default.rteResize = 99
setup.default.moduleData.web_list.bigControlPanel = 1
setup.default.moduleData.web_list.clipBoard = 1
setup.default.moduleData.web_list.localization = 1
setup.default.moduleData.web_list.showPalettes = 1
setup.default.moduleData.file_list.bigControlPanel = 1
setup.default.moduleData.file_list.clipBoard = 1
setup.default.moduleData.file_list.localization = 1
setup.default.moduleData.file_list.showPalettes = 1`

	t3BeUserUC = "a:9:{s:19:\"thumbnailsByDefault\";i:1;s:15:\"recursiveDelete\";i:1;s:25:\"showHiddenFilesAndFolders\";i:1;s:8:\"edit_RTE\";i:1;s:15:\"resizeTextareas\";i:1;s:24:\"resizeTextareas_Flexible\";i:1;s:10:\"copyLevels\";i:99;s:9:\"rteResize\";i:99;s:10:\"moduleData\";a:4:{s:10:\"web_layout\";a:1:{s:8:\"function\";s:1:\"1\";}s:8:\"web_list\";a:4:{s:15:\"bigControlPanel\";s:1:\"1\";s:9:\"clipBoard\";s:1:\"1\";s:12:\"localization\";s:1:\"1\";s:12:\"showPalettes\";s:1:\"1\";}s:6:\"web_ts\";a:5:{s:8:\"function\";s:87:\"TYPO3\\CMS\\Tstemplate\\Controller\\TypoScriptTemplateObjectBrowserModuleFunctionController\";s:15:\"ts_browser_type\";s:5:\"setup\";s:16:\"ts_browser_const\";s:5:\"subst\";s:19:\"ts_browser_fixedLgd\";s:1:\"0\";s:23:\"ts_browser_showComments\";s:1:\"1\";}s:9:\"file_list\";a:4:{s:15:\"bigControlPanel\";s:1:\"1\";s:9:\"clipBoard\";s:1:\"1\";s:12:\"localization\";s:1:\"1\";s:12:\"showPalettes\";s:1:\"1\";}}}"
)

type Typo3BeUser struct {
	Options MysqlCommonOptions `group:"common"`
	Positional struct {
		Schema string `description:"Schema" required:"1"`
	} `positional-args:"true"`
	Username string           `long:"typo3-user"       description:"TYPO3 username"    default:"dev"`
	Password string           `long:"typo3-password"   description:"TYPO3 password"    default:"dev"`
}

func (conf *Typo3BeUser) Execute(args []string) error {
	fmt.Println("Starting TYPO3 BE user generator")
	conf.Options.Init()

	userId := "NULL"

	fmt.Println(" - Creating salted MD5 password")
	password := typo3PasswordGenerator(conf.Password)

	sql := `SELECT uid
              FROM be_users
             WHERE username = %s
               AND deleted = 0`
	sql = fmt.Sprintf(sql, mysqlQuote(conf.Username))
	result := conf.Options.ExecQuery(conf.Positional.Schema, sql)

	for _, row := range result.Row {
		rowList := row.GetList()
		userId = rowList["uid"]
	}

	sql = `INSERT INTO be_users
	                 (uid, tstamp, crdate, realName, username, password, TSconfig, uc, admin, disable, starttime, endtime)
               VALUES(%s, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 'DEVELOPMENT', %s, %s, %s, %s, 1, 0, 0, 0)
	       ON DUPLICATE KEY UPDATE
               realName  = VALUES(realName),
               password  = VALUES(password),
               TSconfig  = VALUES(TSconfig),
               disable   = VALUES(disable),
               starttime = VALUES(starttime),
               endtime   = VALUES(endtime)`
	sql = fmt.Sprintf(sql, userId, mysqlQuote(conf.Username), mysqlQuote(password), mysqlQuote(t3BeUserConfig), mysqlQuote(t3BeUserUC))
	conf.Options.ExecStatement(conf.Positional.Schema, sql)

	if userId != "NULL" {
		fmt.Println(fmt.Sprintf(" - Updated user \"%s\" (UID: %s)", conf.Username, userId))
	} else {
		fmt.Println(fmt.Sprintf(" - Created user \"%s\"", conf.Username))
	}

	return nil
}

func typo3PasswordGenerator(password string) string {
	salt := fmt.Sprintf("$1$%s$", randomString(6))

	crypter := md5_crypt.New()
	ret, _ := crypter.Generate([]byte(password), []byte(salt))

	return ret
}

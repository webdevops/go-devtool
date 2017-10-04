package command

import (
	"github.com/webdevops/go-stubfilegenerator"
	"fmt"
	"strconv"
	"path/filepath"
)

type Typo3Stubs struct {
	Options MysqlCommonOptions `group:"common"`
	Positional struct {
		Schema string `description:"Schema" required:"1"`
		Typo3Root string `description:"TYPO3 root path" required:"1"`
	} `positional-args:"true"`
	Force  bool   `short:"f"  long:"force"      description:"Overwrite existing files"`
}

type storage struct {
	Uid string
	Name string
	Path string
}

type storageFile struct {
	Uid string
	Path string
	RelPath string
	AbsPath string
	ImageWidth string
	ImageHeight string
}

func (conf *Typo3Stubs) Execute(args []string) error {
	fmt.Println("Starting TYPO3 fileadmin stub generator")
	conf.Options.Init()

	sql := `SELECT uid,
                   name,
                   ExtractValue(configuration, '//field[@index=\'basePath\']/value/text()') as storagepath
              FROM sys_file_storage
             WHERE deleted = 0
              AND driver = 'local'`
	result := conf.Options.ExecQuery(conf.Positional.Schema, sql)

	for _, row := range result.Row {
		rowList := row.GetList()
		storage := storage{rowList["uid"], rowList["name"], rowList["storagepath"]}
		conf.processStorage(storage)
	}

	return nil
}

func (conf *Typo3Stubs) processStorage(storage storage) {
	stubgen := stubfilegenerator.StubGenerator()
	stubgen.Image.Text = append(stubgen.Image.Text, "Size: %IMAGE_WIDTH% * %IMAGE_HEIGHT%")

	if conf.Force {
		stubgen.Overwrite = true
	}

	sql := `SELECT f.uid,
                   f.identifier,
                   fm.width as meta_width,
                   fm.height as meta_height
              FROM sys_file f
                   LEFT JOIN sys_file_metadata fm
                     ON fm.file = f.uid
                    AND fm.t3ver_oid = 0
              WHERE f.storage = ` + storage.Uid;
	result := conf.Options.ExecQuery(conf.Positional.Schema, sql)

	for _, row := range result.Row {
		rowList := row.GetList()

		file := storageFile{}
		file.ImageWidth = "800"
		file.ImageHeight = "400"

		switch len(rowList) {
		case 4:
			file.ImageWidth = rowList["meta_width"]
			file.ImageHeight = rowList["meta_height"]
			fallthrough
		case 2:
			file.Uid = rowList["uid"]
			file.Path = filepath.Join(storage.Path, rowList["identifier"])
			file.RelPath = filepath.Join(conf.Positional.Typo3Root, file.Path)
			file.AbsPath, _ = filepath.Abs(file.RelPath)
		}

		stubgen.TemplateVariables["PATH"] = file.Path
		stubgen.TemplateVariables["IMAGE_WIDTH"] = file.ImageWidth
		stubgen.TemplateVariables["IMAGE_HEIGHT"] = file.ImageHeight
		stubgen.Image.Width, _ = strconv.Atoi(file.ImageWidth)
		stubgen.Image.Height, _ = strconv.Atoi(file.ImageHeight)
		stubgen.GenerateStub(file.AbsPath)
	}

}

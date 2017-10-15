package command

import (
	"os"
	"strconv"
	"runtime"
	"path/filepath"
	"github.com/webdevops/go-stubfilegenerator"
	"github.com/remeh/sizedwaitgroup"
)

type Typo3Stubs struct {
	Options MysqlCommonOptions `group:"common"`
	Positional struct {
		Database string `description:"Database" required:"1"`
	} `positional-args:"true"`
	RootPath    string  `long:"path"            description:"TYPO3 root path"`
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
	Logger.Main("Starting TYPO3 fileadmin stub generator")
	if err := conf.Options.Init(); err != nil {
		return err
	}

	rootPath, err := conf.GetTypo3Root()
	if err != nil {
		return err
	}

	sql := `SELECT uid,
                   name,
                   ExtractValue(configuration, '//field[@index=\'basePath\']/value/text()') as storagepath
              FROM sys_file_storage
             WHERE deleted = 0
              AND driver = 'local'`
	result := conf.Options.ExecQuery(conf.Positional.Database, sql)

	for _, row := range result.Row {
		rowList := row.GetList()
		storage := storage{rowList["uid"], rowList["name"], rowList["storagepath"]}
		err := conf.processStorage(storage, rootPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (conf *Typo3Stubs) processStorage(storage storage, rootPath string) error {
	Logger.Step("processing storage \"%v\" (path: %v)", storage.Name, storage.Path)
	stubgen := stubfilegenerator.NewStubGenerator()
	err := stubgen.SetRootPath(filepath.Join(rootPath, storage.Path))
	if err != nil {
		return err
	}

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
	result := conf.Options.ExecQuery(conf.Positional.Database, sql)

	swg := sizedwaitgroup.New(runtime.GOMAXPROCS(0) * 10)
	for _, row := range result.Row {
		row := row.GetList()

		swg.Add()
		go func(row map[string]string, stubgen stubfilegenerator.StubGenerator) {
			defer swg.Done()

			file := storageFile{}
			file.ImageWidth = "800"
			file.ImageHeight = "400"

			switch len(row) {
			case 4:
				file.ImageWidth = row["meta_width"]
				file.ImageHeight = row["meta_height"]
				fallthrough
			case 2:
				file.Uid = row["uid"]
				file.Path = row["identifier"]
			}

			Logger.Item(file.Path)
			stubgen.TemplateVariables["PATH"] = filepath.Join(storage.Path, file.Path)
			stubgen.TemplateVariables["IMAGE_WIDTH"] = file.ImageWidth
			stubgen.TemplateVariables["IMAGE_HEIGHT"] = file.ImageHeight
			stubgen.Image.Width, _ = strconv.Atoi(file.ImageWidth)
			stubgen.Image.Height, _ = strconv.Atoi(file.ImageHeight)
			stubgen.Generate(file.Path)
		} (row, stubgen.Clone())
		swg.Wait()
	}

	return nil
}


func (conf *Typo3Stubs) GetTypo3Root() (string, error) {
	if conf.RootPath == "" {
		// use current working dir as root path
		rootPath, err := os.Getwd()
		if err != nil {
			return "", err
		}
		conf.RootPath = rootPath
	}

	return conf.RootPath, nil
}


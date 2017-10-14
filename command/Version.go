package command

type Version struct {
	ShowOnlyVersion  bool `long:"dump"  description:"show only version number and exit"`
	Name string
	Version string
	Author string
}

func (conf *Version) Execute(args []string) error {
	if conf.ShowOnlyVersion {
		Logger.Println(conf.Version)
	} else {
		Logger.Printlnf("%s version %s", conf.Name, conf.Version)
		Logger.Printlnf("Copyright (C) 2017 %s", conf.Author)
	}

	return nil
}

package config

import (
	"github.com/qdm12/golibs/params"
)

type Filepaths struct {
	Srv  string
	Work string
}

func (f *Filepaths) get(env params.Env) (err error) {
	f.Srv, err = f.getSrv(env)
	if err != nil {
		return err
	}
	f.Work, err = f.getWork(env)
	if err != nil {
		return err
	}
	return nil
}

func (f *Filepaths) getSrv(env params.Env) (filepath string, err error) {
	return env.Path("FILEPATH_SRV", params.Default("./srv"))
}

func (f *Filepaths) getWork(env params.Env) (filepath string, err error) {
	return env.Path("FILEPATH_WORK", params.Default("./work"))
}

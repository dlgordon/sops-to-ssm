package main

import (
	"flag"
)

func AddStandardFlags(f *flag.FlagSet, inputFile *string) *sopsData {
	var data sopsData
	f.Func("sops-file-path", "Specific a path to read your local env files", func(s string) error {
		tmpdata, err := LoadSopsFile(s)
		if err != nil {
			return err
		}
		data = tmpdata
		return nil
	})
	return &data
}

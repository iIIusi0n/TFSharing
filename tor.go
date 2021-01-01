package main

import (
	"github.com/cretz/bine/process/embedded"
	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil/geoipembed"
	"os"
)

func StartClientTor() (*tor.Tor, error) {
	tor, err := tor.Start(nil, &tor.StartConf{
		ProcessCreator: embedded.NewCreator(),
		TempDataDirBase: os.Getenv("TEMP"),
		DebugWriter: nil,
		GeoIPFileReader: geoipembed.GeoIPReader,
	})
	if err != nil {
		return nil, err
	}
	return tor, nil
}

func StartServerTor() (*tor.Tor, error) {
	tor, err := tor.Start(nil, &tor.StartConf{
		ProcessCreator: embedded.NewCreator(),
		DataDir: os.Getenv("TEMP") + "\\data-dir-server",
		DebugWriter: nil,
		GeoIPFileReader: geoipembed.GeoIPReader,
	})
	if err != nil {
		return nil, err
	}
	return tor, nil
}
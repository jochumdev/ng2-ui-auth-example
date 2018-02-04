package shared

import (
	"log"
	"path/filepath"
	"strings"

	ini "gopkg.in/ini.v1"
)

func IniConfigParse(cfile string, debug bool) (string, *ini.File) {
	cfile, _ = filepath.Abs(cfile)

	cfg, err := ini.Load(cfile)
	if err != nil {
		log.Fatalf("ERROR: Failed to read: " + cfile)
	}
	cfg.BlockMode = false // See: http://go-ini.github.io/ini/

	baseDir := filepath.Dir(cfile)
	extends := cfg.Section("").Key("extends").MustString("")
	for _, extend := range strings.Split(extends, ";") {
		if extend == "" {
			continue
		}

		extend = filepath.Join(baseDir, strings.Trim(extend, " "))
		if debug {
			log.Printf("Loading config file: %s\n", extend)
		}
		if err = cfg.Append(extend); err != nil {
			log.Fatalf("ERROR: Loading '%s': %v\n", extend, err)
		}
	}

	return cfile, cfg
}

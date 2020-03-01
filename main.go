// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path"
	"runtime"

	"github.com/api-routerd/api-routerd/cmd/conf"
	"github.com/api-routerd/api-routerd/cmd/router"
	"github.com/api-routerd/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

func main() {
	share.InitLog()

	err := conf.InitConf()
	if err != nil {
		log.Errorf("Failed to init conf file %s: %s", conf.ConfFile, err)
	}

	log.Infof("api-routerd: v%s (built %s)", conf.Version, runtime.Version())
	log.Infof("Start Server at %s:%s", conf.IPFlag, conf.PortFlag)

	err = router.StartRouter(conf.IPFlag, conf.PortFlag, path.Join(conf.ConfPath, conf.TLSCert), path.Join(conf.ConfPath, conf.TLSKey))
	if err != nil {
		log.Fatalf("Failed to init api-routerd: %v", err)
		os.Exit(1)
	}
}

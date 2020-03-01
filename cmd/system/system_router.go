// SPDX-License-Identifier: Apache-2.0

package system

import (
	"github.com/api-routerd/api-routerd/cmd/system/conf"
	"github.com/api-routerd/api-routerd/cmd/system/coredump"
	"github.com/api-routerd/api-routerd/cmd/system/group"
	"github.com/api-routerd/api-routerd/cmd/system/hostname"
	"github.com/api-routerd/api-routerd/cmd/system/journal"
	"github.com/api-routerd/api-routerd/cmd/system/kmod"
	"github.com/api-routerd/api-routerd/cmd/system/login"
	"github.com/api-routerd/api-routerd/cmd/system/resolv"
	"github.com/api-routerd/api-routerd/cmd/system/resolve"
	"github.com/api-routerd/api-routerd/cmd/system/sysctl"
	"github.com/api-routerd/api-routerd/cmd/system/timedate"
	"github.com/api-routerd/api-routerd/cmd/system/timesync"
	"github.com/api-routerd/api-routerd/cmd/system/user"

	"github.com/gorilla/mux"
)

// RegisterRouterSystem register with mux
func RegisterRouterSystem(router *mux.Router) {
	n := router.PathPrefix("/system").Subrouter()

	// system conf
	conf.RegisterRouterSystemConf(n)

	// coredump
	coredump.RegisterRouterCoreDump(n)

	// group
	group.RegisterRouterGroup(n)

	// hostname
	hostname.InitHostname()
	hostname.RegisterRouterHostname(n)

	// journald
	journal.InitJournalConf()
	journal.RegisterRouterJournal(n)

	// kmod
	kmod.RegisterRouterKMod(n)

	// login
	login.RegisterRouterLogin(n)

	// sysctl
	sysctl.RegisterRouterSysctl(n)

	// /etc/resolv
	resolv.RegisterRouterResolv(n)

	// resolved
	resolve.RegisterRouterResolve(n)

	// timedate
	timedate.InitTimeDate()
	timedate.RegisterRouterTimeDate(n)

	// timesync
	timesync.RegisterRouterTimeSync(n)

	// user
	user.RegisterRouterUser(n)
}

//go:build darwin
// +build darwin

package skyenv

const (
	//OS detection at runtime
	OS = "mac"
	// SkywirePath is the path to the installation folder.
	SkywirePath = "/Library/Application Support/Skywire"
	// Configjson is the config name generated by the script included with the installation on mac
	Configjson = ConfigName
)

// PackageConfig contains installation paths (for mac)
func PackageConfig() PkgConfig {
	var pkgconfig PkgConfig
	pkgconfig.Launcher.BinPath = "/Applications/Skywire.app/Contents/MacOS/apps"
	pkgconfig.LocalPath = "/Library/Application Support/Skywire/local"
	pkgconfig.Hypervisor.DbPath = "/Library/Application Support/Skywire/users.db"
	pkgconfig.Hypervisor.EnableAuth = true
	return pkgconfig
}

// UserConfig contains installation paths (for mac)
func UserConfig() PkgConfig {
	var usrconfig PkgConfig
	usrconfig.Launcher.BinPath = "/Applications/Skywire.app/Contents/MacOS/apps"
	usrconfig.LocalPath = HomePath() + "/.skywire/local"
	usrconfig.Hypervisor.DbPath = HomePath() + "/.skywire/users.db"
	usrconfig.Hypervisor.EnableAuth = true
	return usrconfig
}

// UpdateCommand returns the commands which are run when the update button is clicked in the ui
func UpdateCommand() []string {
	return []string{`echo "update not implemented for macOS. Download a new version from the release section here: https://github.com/skycoin/skywire/releases"`}
}

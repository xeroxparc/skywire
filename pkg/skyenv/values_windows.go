//go:build windows
// +build windows

package skyenv

const (
	//OS detection at runtime
	OS = "win"
	// SkywirePath is the path to the installation folder for the .msi
	SkywirePath = "C:/Program Files/Skywire"
)

// PackageConfig contains installation paths (for windows)
func PackageConfig() PkgConfig {
	var pkgconfig PkgConfig
	pkgconfig.Launcher.BinPath = "C:/Program Files/Skywire/apps"
	pkgconfig.LocalPath = "C:/Program Files/Skywire/local"
	pkgconfig.Hypervisor.DbPath = "C:/Program Files/Skywire/users.db"
	pkgconfig.Hypervisor.EnableAuth = true
	return pkgconfig
}

// UserConfig contains installation paths (for windows)
func UserConfig() PkgConfig {
	var usrconfig PkgConfig
	usrconfig.Launcher.BinPath = "C:/Program Files/Skywire/apps"
	usrconfig.LocalPath = HomePath() + "/.skywire/local"
	usrconfig.Hypervisor.DbPath = HomePath() + "/.skywire/users.db"
	usrconfig.Hypervisor.EnableAuth = true
	return usrconfig
}

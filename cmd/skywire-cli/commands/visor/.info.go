package visor

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var addInput string

func init() {
	RootCmd.AddCommand(pkCmd)
	pkCmd.Flags().StringVarP(&addInput, "input", "i", "", "read from specified config file")

	RootCmd.AddCommand(hvpkCmd)
	hvpkCmd.Flags().StringVarP(&addInput, "input", "i", "", "read from specified config file")

	RootCmd.AddCommand(
		summaryCmd,
		buildInfoCmd,
	)
}

var pkCmd = &cobra.Command{
	Use:   "pk",
	Short: "public key of the visor",
	Run: func(_ *cobra.Command, _ []string) {
		//		if addInput != "" {
		//			conf := visorconfig.ReadConfig(addInput)
		//			fmt.Println(conf.PK.Hex())
		//		} else {
		client := rpcClient()
		overview, err := client.Overview()
		if err != nil {
			logger.Fatal("Failed to connect:", err)
			//			}
			fmt.Println(overview.PubKey)
		}
	},
}

var hvpkCmd = &cobra.Command{
	Use:   "hv",
	Short: "show hypervisor(s)",
	Run: func(_ *cobra.Command, _ []string) {
		//		if addInput != "" {
		//			conf := visorconfig.ReadConfig(addInput)
		//			fmt.Println(conf.Hypervisors)
		//		} else {
		client := rpcClient()
		overview, err := client.Overview()
		if err != nil {
			logger.Fatal("Failed to connect:", err)
			//			}
			fmt.Println(overview.Hypervisors)
		}
	},
}

var summaryCmd = &cobra.Command{
	Use:   "info",
	Short: "summary of visor info",
	Run: func(_ *cobra.Command, _ []string) {
		summary, err := rpcClient().Summary()
		if err != nil {
			log.Fatal("Failed to connect:", err)
		}
		msg := fmt.Sprintf(".:: Visor Summary ::.\nPublic key: %q\nSymmetric NAT: %t\nIP: %s\nDMSG Server: %q\nPing: %q\nVisor Version: %s\nSkybian Version: %s\nUptime Tracker: %s\nTime Online: %f seconds\nBuild Tag: %s\n", summary.Overview.PubKey, summary.Overview.IsSymmetricNAT, summary.Overview.LocalIP, summary.DmsgStats.ServerPK, summary.DmsgStats.RoundTrip, summary.Overview.BuildInfo.Version, summary.SkybianBuildVersion, summary.Health.ServicesHealth, summary.Uptime, summary.BuildTag)
		if _, err := os.Stdout.Write([]byte(msg)); err != nil {
			log.Fatal("Failed to output build info:", err)
		}
	},
}

var buildInfoCmd = &cobra.Command{
	Use:   "version",
	Short: "version and build info",
	Run: func(_ *cobra.Command, _ []string) {
		client := rpcClient()
		overview, err := client.Overview()
		if err != nil {
			log.Fatal("Failed to connect:", err)
		}

		if _, err := overview.BuildInfo.WriteTo(os.Stdout); err != nil {
			log.Fatal("Failed to output build info:", err)
		}
	},
}

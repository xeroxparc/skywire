package commands

import (
	"fmt"
	"net/http"
	"os"

	"github.com/SkycoinProject/dmsg"
	"github.com/SkycoinProject/dmsg/disc"
	"github.com/SkycoinProject/skycoin/src/util/logging"
	"github.com/spf13/cobra"

	"github.com/SkycoinProject/skywire-mainnet/pkg/hypervisor"
	"github.com/SkycoinProject/skywire-mainnet/pkg/util/buildinfo"
	"github.com/SkycoinProject/skywire-mainnet/pkg/util/pathutil"
)

const configEnv = "SW_HYPERVISOR_CONFIG"

// nolint:gochecknoglobals
var (
	log = logging.MustGetLogger("hypervisor")

	configPath     string
	mock           bool
	mockEnableAuth bool
	mockVisors     int
	mockMaxTps     int
	mockMaxRoutes  int
)

// nolint:gochecknoinits
func init() {
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "./hypervisor-config.json", "hypervisor config path")
	rootCmd.Flags().BoolVarP(&mock, "mock", "m", false, "whether to run hypervisor with mock data")
	rootCmd.Flags().BoolVar(&mockEnableAuth, "mock-enable-auth", false, "whether to enable user management in mock mode")
	rootCmd.Flags().IntVar(&mockVisors, "mock-visors", 5, "number of visors to have in mock mode")
	rootCmd.Flags().IntVar(&mockMaxTps, "mock-max-tps", 10, "max number of transports per mock visor")
	rootCmd.Flags().IntVar(&mockMaxRoutes, "mock-max-routes", 30, "max number of routes per visor")
}

// nolint:gochecknoglobals
var rootCmd = &cobra.Command{
	Use:   "hypervisor",
	Short: "Manages Skywire Visors",
	Run: func(_ *cobra.Command, args []string) {
		if _, err := buildinfo.Get().WriteTo(os.Stdout); err != nil {
			log.Printf("Failed to output build info: %v", err)
		}

		// Prepare config.
		if configPath == "" {
			configPath = pathutil.FindConfigPath(args, -1, configEnv, pathutil.HypervisorDefaults())
		}
		var conf hypervisor.Config
		conf.FillDefaults(mock)
		if err := conf.Parse(configPath); err != nil {
			log.WithError(err).Fatalln("failed to parse config file")
		}
		log.WithField("config", conf).
			Info()

		// Prepare hypervisor.
		m, err := hypervisor.New(conf)
		if err != nil {
			log.Fatalln("Failed to start hypervisor:", err)
		}

		if mock {
			// Mock mode.
			err := m.AddMockData(hypervisor.MockConfig{
				Visors:            mockVisors,
				MaxTpsPerVisor:    mockMaxTps,
				MaxRoutesPerVisor: mockMaxRoutes,
				EnableAuth:        mockEnableAuth,
			})
			if err != nil {
				log.Fatalln("Failed to add mock data:", err)
			}
		} else {
			// Prepare dmsg client.
			dmsgC := dmsg.NewClient(conf.PK, conf.SK, disc.NewHTTP(conf.DmsgDiscovery), dmsg.DefaultConfig())
			go dmsgC.Serve()

			dmsgL, err := dmsgC.Listen(conf.DmsgPort)
			if err != nil {
				log.WithField("addr", fmt.Sprintf("dmsg://%s:%d", conf.PK, conf.DmsgPort)).
					Fatal("Failed to listen over dmsg.")
			}
			go func() {
				if err := m.ServeRPC(dmsgC, dmsgL); err != nil {
					log.WithError(err).
						Fatal("Failed to serve RPC client over dmsg.")
				}
			}()
			log.WithField("addr", fmt.Sprintf("dmsg://%s:%d", conf.PK, conf.DmsgPort)).
				Info("Serving RPC client over dmsg.")
		}

		// Serve HTTP.
		log.WithField("http_addr", conf.HttpAddr).
			Info("Serving HTTP.")

		if err := http.ListenAndServe(conf.HttpAddr, m); err != nil {
			log.WithError(err).
				Fatal("Hypervisor exited with error.")
		}
		log.Info("Good bye!")
	},
}

// Execute executes root CLI command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

package dmsgpty

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/skycoin/dmsg/pkg/dmsgpty"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/spf13/cobra"

	"github.com/skycoin/skywire-utilities/pkg/cmdutil"
	"github.com/skycoin/skywire/cmd/skywire-cli/internal"
	"github.com/skycoin/skywire/pkg/visor"
)

var rpcAddr string
var ptyPort string
var masterLogger = logging.NewMasterLogger()
var packageLogger = masterLogger.PackageLogger("dmsgpty")

func init() {
	RootCmd.PersistentFlags().StringVarP(&rpcAddr, "rpc", "", "localhost:3435", "RPC server address")
	RootCmd.PersistentFlags().StringVarP(&ptyPort, "port", "p", "22", "port of remote visor dmsgpty")
}

// RootCmd is the command that contains sub-commands which interacts with dmsgpty.
var RootCmd = &cobra.Command{
	Use:   "dmsgpty",
	Short: "Some simple commands of dmsgpty to remote visor",
}

func init() {
	RootCmd.AddCommand(
		listOfVisors,
		executeCommand,
	)
}

var listOfVisors = &cobra.Command{
	Use:   "list",
	Short: "fetch list of connected visor's PKs",
	Run: func(_ *cobra.Command, _ []string) {
		remoteVisors := rpcClient().RemoteVisors()
		var msg string
		for idx, pk := range remoteVisors {
			msg += fmt.Sprintf("%d. %s\n", idx+1, pk)
		}
		if _, err := os.Stdout.Write([]byte(msg)); err != nil {
			packageLogger.Fatal("Failed to output build info:", err)
		}
	},
}

var executeCommand = &cobra.Command{
	Use:   "start <pk>",
	Short: "start dmsgpty for specific visor by its dmsg address pk:port",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		cli := dmsgpty.DefaultCLI()
		addr := internal.ParsePK("pk", args[0])
		port, _ := strconv.ParseUint(ptyPort, 10, 16) //nolint
		ctx, cancel := cmdutil.SignalContext(context.Background(), nil)
		defer cancel()
		return cli.StartRemotePty(ctx, addr, uint16(port), dmsgpty.DefaultCmd)
	},
}

func rpcClient() visor.API {
	const rpcDialTimeout = time.Second * 5
	conn, err := net.DialTimeout("tcp", rpcAddr, rpcDialTimeout)
	if err != nil {
		packageLogger.Fatal("RPC connection failed:", err)
	}
	return visor.NewRPCClient(packageLogger, conn, visor.RPCPrefix, 0)
}
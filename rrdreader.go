package main

import (
	"github.com/spf13/cobra"
	"github.com/tonypoch/rrdreader.git/clicmd"
)

func main() {
	var rootCmd = &cobra.Command{Use: "rrd-reader"}
	var address string
	var port int
	var rootdir string
	var outputdir string

	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "ip address")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "TCP port")
	rootCmd.PersistentFlags().StringVarP(&rootdir, "rootdir", "d", "", "collectd root directory")
	rootCmd.PersistentFlags().StringVarP(&outputdir, "output", "o", "/tmp/images/", "output image directory")

	var cmdRrdRun = &cobra.Command{
		Use:   "run",
		Short: "Run rrd-reader in server mode",
		Run:   clicmd.CmdRun,
	}
	rootCmd.AddCommand(cmdRrdRun)
	rootCmd.Execute()
}

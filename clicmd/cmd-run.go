package clicmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/tonypoch/rrdreader/htmlviewer"
	"github.com/tonypoch/rrdreader/rrdtool"
)

func CmdRun(cmd *cobra.Command, args []string) {
	//Initialize our config structure
	var presetPtr *rrdtool.GraphPreset

	rpcAddress, _ := cmd.Flags().GetString("address")
	rpcPort, _ := cmd.Flags().GetInt("port")
	rootDir, _ := cmd.Flags().GetString("rootdir")
	outputDir, _ := cmd.Flags().GetString("output")

	hostame, _ := os.Hostname()

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		log.Fatalf("No such directory: %s", outputDir)
	}

	if rpcPort == 0 {
		rpcPort = 8080
	}
	if rootDir == "" {
		rootDir = "/var/lib/collectd/" + hostame
	}

	connectString := fmt.Sprintf("%s:%d", rpcAddress, rpcPort)

	go htmlviewer.RunHTMLServer(connectString, outputDir)

	for {
		for group, config := range rrdtool.Configs {
			pattern := rootDir + "/" + config.Directory
			matches, err := filepath.Glob(pattern)
			if err != nil {
				fmt.Println(err)
			}
			for _, match := range matches {

				rrdfiles, err := rrdtool.FileSearcher(match)
				if err != nil {
					log.Printf("Error in test: %s", err)
				}

				switch group {
				case "Network":
					rrdtool.NetworkPreset.VarsNames = nil
					presetPtr = &rrdtool.NetworkPreset
				case "Memory":
					rrdtool.SystemPreset.VarsNames = nil
					presetPtr = &rrdtool.SystemPreset
				case "Disk":
					rrdtool.DiskPreset.VarsNames = nil
					presetPtr = &rrdtool.DiskPreset
				case "CPU":
					rrdtool.SystemPreset.VarsNames = nil
					presetPtr = &rrdtool.SystemPreset
				case "Load":
					rrdtool.SystemPreset.VarsNames = nil
					presetPtr = &rrdtool.SystemPreset
				case "Temp":
					rrdtool.SystemPreset.VarsNames = nil
					presetPtr = &rrdtool.SystemPreset
				case "SMART":
					rrdtool.SystemPreset.VarsNames = nil
					presetPtr = &rrdtool.SystemPreset
				default:
					log.Fatalln("There is no such config!")
				}

				rrdtool.RRDProcess(rrdfiles, presetPtr, config, outputDir)
			}
		}
		time.Sleep(300 * time.Second)
	}
}

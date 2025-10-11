package player

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	simuCmd.Flags().StringP("prefix", "P", "", "user (required)")
	simuCmd.MarkFlagRequired("prefix")
	simuCmd.Flags().IntP("count", "C", 10, "simulator count")
	simuCmd.Flags().StringP("host", "H", "http://192.168.1.11", "use default")
}

var simuCmd = &cobra.Command{
	Use:   "simu",
	Short: "simulation",
	Long:  "simulation load test",
	Run: func(cmd *cobra.Command, args []string) {
		prefix, _ := cmd.Flags().GetString("prefix")
		count, _ := cmd.Flags().GetInt("count")
		host, _ := cmd.Flags().GetString("host")
		var wt sync.WaitGroup
		wt.Add(count)
		start := time.Now()
		for i := range count {
			go func() {
				sim := Simulator{Player: fmt.Sprintf("%s_%d", prefix, i), Host: host}
				err := sim.Play()
				if err != nil {
					fmt.Printf("Failed from %d %s\n", i,err.Error())
				}
				wt.Done()
			}()
		}
		wt.Wait()
		dur := time.Since(start)
		fmt.Printf("Total duration %d\n", dur.Milliseconds())
	},
}

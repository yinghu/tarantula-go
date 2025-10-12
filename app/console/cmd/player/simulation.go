package player

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	simuCmd.Flags().StringP("prefix", "P", "", "user (required)")
	simuCmd.MarkFlagRequired("prefix")
	simuCmd.Flags().IntP("count", "C", 10, "simulator count")
	simuCmd.Flags().IntP("batch", "B", 10, "simulator batch")
	simuCmd.Flags().StringP("host", "H", "http://192.168.1.11", "use default")
}

var simuCmd = &cobra.Command{
	Use:   "simu",
	Short: "simulation",
	Long:  "simulation load test",
	Run: func(cmd *cobra.Command, args []string) {
		prefix, _ := cmd.Flags().GetString("prefix")
		count, _ := cmd.Flags().GetInt("count")
		batch, _ := cmd.Flags().GetInt("batch")
		host, _ := cmd.Flags().GetString("host")
		wt := make(chan bool, 1)
		bt := make(chan int, 1)
		done := false
		suc := 0
		fail := 0
		sub := 0
		tc := 0
		start := time.Now()
		bt <- 1
		for {
			if done {
				break
			}
			select {
			case t := <-wt:
				sub++
				if t {
					suc++
				} else {
					fail++
				}
				if sub == count {
					tc++
					if tc == batch {
						done = true
					} else {
						sub = 0
						bt <- 1
					}
				}
			case <-bt:
				for i := range count {
					go func() {
						sim := Simulator{Player: fmt.Sprintf("%s_%d_%d", prefix, tc, i), Host: host}
						err := sim.Play()
						wt <- err == nil
					}()
				}
			}
		}
		dur := time.Since(start)
		avg := dur.Milliseconds() / int64(count*batch)
		fmt.Printf("avg dur per call (ms) :%d with success [%d] failure [%d]\n", avg, suc, fail)
	},
}

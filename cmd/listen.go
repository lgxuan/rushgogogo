/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"rushgogogo/pkgs/config"
	"rushgogogo/pkgs/proxy"

	"github.com/spf13/cobra"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "The address you want to listen to. (default ':8081')",
	Long: `Start an HTTP proxy server on the given address (default :8081).
Any traffic sent to this port will be forwarded through the proxy,
allowing you to inspect, modify or log requests before they reach
their original destination.`,
	Run: func(cmd *cobra.Command, args []string) {
		port := ":8081"
		if len(args) > 0 {
			port = args[0]
		}

		// 获取线程数设置
		threads, _ := cmd.Flags().GetInt("threads")
		if threads > 0 {
			// 更新配置文件中的线程数
			if err := config.UpdateThreadCount("config.yaml", threads); err != nil {
				fmt.Printf("Warning: failed to update thread count: %v\n", err)
			} else {
				fmt.Printf("Thread count set to %d\n", threads)
			}
		}

		fmt.Println("Listening on port", port)
		proxy.ListenAddress(port)
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)

	// 添加线程数参数
	listenCmd.Flags().IntP("threads", "t", 10, "Number of message processing threads")
}

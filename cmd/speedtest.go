/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/showwin/speedtest-go/speedtest"
	"github.com/spf13/cobra"
)

var speedtestCmd = &cobra.Command{
	Use:   "speedtest",
	Short: "Uji kecepatan internet (Ping, Download, Upload)",
	Long:  `Perintah ini akan mencari server terdekat dan melakukan pengujian kecepatan jaringan WiFi atau LAN kamu.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sedang mencari server terdekat...")

		var speedtestClient = speedtest.New()

		serverList, err := speedtestClient.FetchServers()
		if err != nil {
			fmt.Printf("Gagal mengambil daftar server: %v\n", err)
			return
		}

		targets, err := serverList.FindServer([]int{})
		if err != nil || len(targets) == 0 {
			fmt.Println("Gagal menemukan server terdekat.")
			return
		}

		s := targets[0]
		fmt.Printf("Server ditemukan: %s (%s)\n", s.Sponsor, s.Name)
		fmt.Println("Memulai pengujian... (Mohon tunggu beberapa detik)")

		err = s.PingTest(nil)
		if err != nil {
			fmt.Printf("Gagal tes Ping: %v\n", err)
		}

		err = s.DownloadTest()
		if err != nil {
			fmt.Printf("Gagal tes Download: %v\n", err)
		}

		err = s.UploadTest()
		if err != nil {
			fmt.Printf("Gagal tes Upload: %v\n", err)
		}

		fmt.Printf("\n=== HASIL SPEED TEST ===\n")
		fmt.Printf("Ping     : %d ms\n", s.Latency.Milliseconds())
		fmt.Printf("Download : %.2f Mbps\n", s.DLSpeed/1000000)
		fmt.Printf("Upload   : %.2f Mbps\n", s.ULSpeed/1000000)
		fmt.Println("==============================")
	},
}

func init() {
	rootCmd.AddCommand(speedtestCmd)
}

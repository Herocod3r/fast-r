package commands

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	tm "github.com/buger/goterm"

	"github.com/herocod3r/fast-r/pkg/cache/file"
	"github.com/herocod3r/fast-r/pkg/client"

	"github.com/apoorvam/goterminal"
	ct "github.com/daviddengcn/go-colortext"

	"github.com/herocod3r/fast-r/pkg/network"

	fast_r "github.com/herocod3r/fast-r"

	"github.com/spf13/cobra"
)

var (
	config = fast_r.Config{EnableCaching: true}
)

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Executes a speed test",
		Run:   executeRun,
	}
	cmd.Flags().DurationVarP(&config.MaxTime, "max-time", "m", time.Minute*3, "Maximum Time To Run A Test, eg 10m for 10 minutes")
	cmd.Flags().IntVarP(&config.MaximumConnections, "connections-limit", "l", 1, "Parallel Connections")
	cmd.Flags().BoolVar(&config.EnableCaching, "enable-cache", true, "Enable caching")
	return cmd
}

func executeRun(cmd *cobra.Command, args []string) {
	if config.MaximumConnections < 1 || config.MaximumConnections > 10 {
		config.MaximumConnections = 1
	}

	if config.MaxTime.Minutes() < 1 || config.MaxTime.Minutes() > 10 {
		config.MaxTime = time.Minute * 3
	}

	server, err := executeGetServers()
	if err != nil {
		return
	}

	fmt.Println("")

	downloadClient, err := executeDownload(server)
	if err != nil {
		return
	}

	fmt.Println("")

	uploadClient, err := executeUpload(server)
	if err != nil {
		return
	}

	ct.ResetColor()
	tm.Clear()
	ct.Foreground(ct.Yellow, true)
	fmt.Println(fmt.Sprintf("Your Download Speed Is %.2f Mbs ", (downloadClient.Speed)/1000))
	fmt.Println(fmt.Sprintf("Your Upload Speed Is %.2f Mbs ", (uploadClient.Speed)/1000))
	//fmt.Println(fmt.Sprintf("Your Upload Speed Is %.1f Mbs", (float64(uploadStream.TotalBytes)/float64(125000))/uploadStream.TotalTime.Seconds()))
	ct.ResetColor()
}

func executeGetServers() (*network.Server, error) {
	writer := goterminal.New(os.Stdout)
	ct.Foreground(ct.White, true)
	fmt.Fprintln(writer, "Selecting server ...")
	writer.Print()
	//ct.ResetColor()
	servers, er := client.GetServersList(config.EnableCaching, file.NewFileSystemCache())
	if er != nil || len(servers) < 1 {
		ct.Foreground(ct.Red, true)
		if errors.Is(er, network.NetworkAccessErr) {
			fmt.Fprintln(writer, "A network error has occurred, please check and ensure your network is still active")
			return nil, er
		}
		fmt.Fprintln(writer, "An Error has occurred, please try again later")
		writer.Print()
		ct.ResetColor()
		return nil, er
	}
	fmt.Println()
	server := &servers[0]
	fmt.Fprintln(writer, "")

	ct.Foreground(ct.Blue, true)
	fmt.Fprintln(writer, "----------SERVER SELECTED-----------")
	fmt.Fprintln(writer, "====================================")
	fmt.Fprintln(writer, "::Server::     ", server.Name)
	fmt.Fprintln(writer, "====================================")
	writer.Print()
	ct.ResetColor()
	return server, nil
}

func executeDownload(server *network.Server) (*client.Download, error) {
	syncLock := new(sync.Mutex)
	ct.Foreground(ct.Green, true)
	tm.Clear()
	downloadClient := client.NewDownload(server, func(curSpeed float64) {
		syncLock.Lock()
		speed := curSpeed / 1000 //kib => mib
		tm.Flush()               // Call it every time at the end of rendering
		tm.MoveCursor(1, 1)

		tm.Println(fmt.Sprintf("Download Speed is %.2f Mbs   ", speed))
		tm.Flush()
		syncLock.Unlock()
	})

	er := downloadClient.Start()
	if er != nil && downloadClient.Speed < 1 {
		fmt.Println("An error occurred unable to complete the request", er.Error())
		return nil, er
	}
	ct.ResetColor()
	tm.Clear()
	return downloadClient, nil
}

func executeUpload(server *network.Server) (*client.Upload, error) {
	syncLock := new(sync.Mutex)
	ct.Foreground(ct.Green, true)
	tm.Clear()
	uploadClient := client.NewUpload(server, func(curSpeed float64) {
		syncLock.Lock()
		speed := curSpeed / 1000 //kib => mib
		tm.Flush()               // Call it every time at the end of rendering
		tm.MoveCursor(1, 1)

		tm.Println(fmt.Sprintf("Upload Speed is %.2f Mbs   ", speed))
		tm.Flush()
		syncLock.Unlock()
	})

	er := uploadClient.Start()
	if er != nil && uploadClient.Speed < 1 {
		fmt.Println("An error occurred unable to complete the request", er.Error())
		return nil, er
	}
	ct.ResetColor()
	tm.Clear()
	return uploadClient, nil
}

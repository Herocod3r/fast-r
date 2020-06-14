package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/apoorvam/goterminal"
	ct "github.com/daviddengcn/go-colortext"

	cache2 "github.com/herocod3r/fast-r/pkg/cache"

	"github.com/herocod3r/fast-r/pkg/cache/file"
	"github.com/herocod3r/fast-r/pkg/network/http"

	"github.com/herocod3r/fast-r/pkg/network"

	fast_r "github.com/herocod3r/fast-r"

	"github.com/spf13/cobra"
)

var (
	config = fast_r.Config{}
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
	writer := goterminal.New(os.Stdout)
	ct.Foreground(ct.White, true)
	fmt.Fprintln(writer, "Selecting server ...")
	writer.Print()
	//ct.ResetColor()
	server, er := getServer()
	if er != nil {
		ct.Foreground(ct.Red, true)
		if errors.Is(er, network.NetworkAccessErr) {
			fmt.Fprintln(writer, "A network error has occured, please check and ensure your network is still active")
			return
		}
		fmt.Fprintln(writer, "An Error has occured, please try again later")
		writer.Print()
		ct.ResetColor()
	}
	fmt.Fprintln(writer, "")

	ct.Foreground(ct.Blue, false)
	fmt.Fprintln(writer, "----------SERVER SELECTED-----------")
	fmt.Fprintln(writer, "====================================")
	fmt.Fprintln(writer, "::Server::     ", server.Name)
	fmt.Fprintln(writer, "====================================")
	writer.Print()
	ct.ResetColor()

	//if _, ok := er.(*network.Error); ok {
	//	//log network error
	//}

}

func getServer() (*network.Server, error) {
	service := http.NewSpeedTestService()
	if config.EnableCaching {
		server, _ := getServerFromCache()
		if server != nil {
			return server, nil
		}
	}

	client, er := service.GetClientInfo()
	if er != nil {
		return nil, er
	}

	servers, er := service.GetServers(1, client)
	if er != nil {
		return nil, er
	}
	cache := file.NewFileSystemCache()
	data, _ := json.Marshal(servers[0])
	cache.Set("server", string(data))
	return &servers[0], nil
}

func getServerFromCache() (*network.Server, error) {
	cache := file.NewFileSystemCache()
	clientData, er := cache.Get("server")
	if er != nil {
		if errors.Is(er, cache2.StoreNotActiveErr) {
			//log caching not supported
		}
		return nil, er
	}
	client := network.Server{}
	er = json.Unmarshal([]byte(clientData), &client)
	if er != nil {
		return nil, er
	}
	return &client, nil
}

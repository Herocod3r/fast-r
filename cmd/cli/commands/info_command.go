package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/herocod3r/fast-r/pkg/network"
	"github.com/herocod3r/fast-r/pkg/network/http"

	"github.com/apoorvam/goterminal"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/spf13/cobra"
)

func NewInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Fetches your network info",
		Run:   getUserInfo,
	}
}

func getUserInfo(cmd *cobra.Command, args []string) {

	writer := goterminal.New(os.Stdout)
	ct.Foreground(ct.White, true)
	fmt.Fprintln(writer, "Processing your client info ...")
	writer.Print()
	ct.ResetColor()
	service := http.NewSpeedTestService()
	client, er := service.GetClientInfo()
	writer.Clear()
	if errors.Is(er, network.NetworkAccessErr) {
		ct.Foreground(ct.Red, true)
		fmt.Fprintln(writer, "Unable to get your network information, make sure you are connected and try again")
		writer.Print()
		ct.ResetColor()
		return
	}

	ct.Foreground(ct.Green, true)
	fmt.Fprintln(writer, "====================================")
	fmt.Fprintln(writer, "::IP::     ", client.Ip)
	fmt.Fprintln(writer, "::ISP::    ", client.Isp)
	fmt.Fprintln(writer, "::COUNTRY::", client.Location)
	fmt.Fprintln(writer, "====================================")
	writer.Print()
	ct.ResetColor()
	return

}

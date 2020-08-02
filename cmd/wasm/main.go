package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"syscall/js"

	"github.com/herocod3r/fast-r/pkg/cache/file"

	"github.com/herocod3r/fast-r/pkg/client"
	"github.com/herocod3r/fast-r/pkg/network"

	"github.com/herocod3r/fast-r/pkg/network/http"
)

var (
	done = make(chan struct{})
)

func main() {
	http.ConfigUrl = "https://cors-anywhere.herokuapp.com/" + http.ConfigUrl
	http.ServerUrl = "https://cors-anywhere.herokuapp.com/" + http.ServerUrl

	printCallback := js.FuncOf(printResult)
	defer printCallback.Release()

	getServerCallback := js.FuncOf(getServer)
	defer getServerCallback.Release()

	setPrintResult := js.Global().Get("setPrintResult")
	if !setPrintResult.IsNull() {
		setPrintResult.Invoke(printCallback)
	}

	setGetServerResult := js.Global().Get("setGetServerResult")
	if !setGetServerResult.IsNull() {
		setGetServerResult.Invoke(getServerCallback)
	}

	downloadResult := js.Global().Get("setDownloadTestResult")
	downloadCallback := js.FuncOf(startDownload)
	defer downloadCallback.Release()
	if !downloadResult.IsNull() {
		downloadResult.Invoke(downloadCallback)
	}

	uploadResult := js.Global().Get("setUploadTestResult")
	uploadCallback := js.FuncOf(startUpload)
	defer uploadCallback.Release()
	if !uploadResult.IsNull() {
		uploadResult.Invoke(uploadCallback)
	}

	<-done
}

func printResult(value js.Value, args []js.Value) interface{} {
	value1 := args[0].String()
	v1, err := strconv.Atoi(value1)
	if err != nil {
		fmt.Errorf("error %s", err.Error())
		return err
	}
	value2 := args[1].String()
	v2, err := strconv.Atoi(value2)
	if err != nil {
		fmt.Errorf("error %s", err.Error())
		return err
	}

	fmt.Printf("%d\n", v1+v2)
	done <- struct{}{}
	return nil
}

func getServer(value js.Value, args []js.Value) interface{} {
	callback := args[0]

	go func() {
		list, err := client.GetServersList(true, file.NewFileSystemCache())
		resp := struct {
			Error   string
			Servers []network.Server
		}{}
		if err != nil {
			resp = struct {
				Error   string
				Servers []network.Server
			}{Error: err.Error()}
		} else {
			resp = struct {
				Error   string
				Servers []network.Server
			}{Servers: list}
		}
		fmt.Print(callback)
		objBytes, _ := json.Marshal(resp)
		callback.Invoke(js.ValueOf(string(objBytes)))
		done <- struct{}{}

	}()
	return nil
}

func startDownload(value js.Value, args []js.Value) interface{} {
	serverJson := args[0].String()
	server := network.Server{}
	_ = json.Unmarshal([]byte(serverJson), &server)
	curCallback := args[1]
	finalCallback := args[2]
	var downloadClient *client.Download

	go func() {
		syncLock := new(sync.Mutex)
		downloadClient = client.NewDownload(&server, func(curSpeed float64) {
			syncLock.Lock()
			curCallback.Invoke(curSpeed)
			syncLock.Unlock()
		})

		_ = downloadClient.Start()
		finalCallback.Invoke(downloadClient.Speed)
		done <- struct{}{}
	}()

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			done <- struct{}{}
		}()
		if downloadClient == nil {
			return js.ValueOf("Process Has Not Started Because client is nil")
		}

		er := downloadClient.Stop()
		if er != nil {
			return js.ValueOf(er.Error())
		}
		return nil
	})
}

func startUpload(value js.Value, args []js.Value) interface{} {
	serverJson := args[0].String()
	server := network.Server{}
	_ = json.Unmarshal([]byte(serverJson), &server)
	curCallback := args[1]
	finalCallback := args[2]
	var uploadClient *client.Upload

	go func() {
		syncLock := new(sync.Mutex)
		uploadClient = client.NewUpload(&server, func(curSpeed float64) {
			syncLock.Lock()
			curCallback.Invoke(curSpeed)
			syncLock.Unlock()
		})

		_ = uploadClient.Start()
		finalCallback.Invoke(uploadClient.Speed)
		done <- struct{}{}
	}()

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			done <- struct{}{}
		}()
		if uploadClient == nil {
			return js.ValueOf("Process Has Not Started")
		}

		er := uploadClient.Stop()
		if er != nil {
			return js.ValueOf(er.Error())
		}
		return nil
	})
}

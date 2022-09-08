package main

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

const (
	Proxy_Path                       = "Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings"
	Proxy_Enable                     = "ProxyEnable"
	Proxy_Server                     = "ProxyServer"
	Wininet_Dll                      = "wininet.dll"
	Internet_Set_Option              = "InternetSetOptionW"
	INTERNET_OPTION_REFRESH          = 37
	INTERNET_OPTION_SETTINGS_CHANGED = 39
)

func ProxyFlush() {
	dll := syscall.NewLazyDLL(Wininet_Dll)
	fn := dll.NewProc(Internet_Set_Option)
	fn.Call(0, INTERNET_OPTION_REFRESH, 0, 0)
	fn.Call(0, INTERNET_OPTION_SETTINGS_CHANGED, 0, 0)
}
func main() {
	key, err := registry.OpenKey(registry.CURRENT_USER, Proxy_Path, registry.ALL_ACCESS)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	bin := make([]byte, 8)
	bin[7] = 1
	err = key.SetDWordValue(Proxy_Enable, 0)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	err = key.SetStringValue(Proxy_Server, "127.0.0.1:80")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

}

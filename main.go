package main

import(
	"fmt"
	"runtime"
	"os"
	"strings"
	"time"
	"crypto/tls"
	libmessage "github.com/ProhtMeyhet/libgomessage"
)

func main() {
	runtime.GOMAXPROCS(4)

	flags := newFlagConfig()
	flags.parse()
	validate(flags)

	if flags.Echo {
		fmt.Println(flags.Title)
		fmt.Print(flags.Message)
	}

	message := &libmessage.Message{ Title: flags.Title,
					Message: flags.Message,
					Icon: flags.Icon }

	switch(strings.ToLower(flags.Type)) {
	case EMPTY:
		flags.usage()
		os.Exit(EMPTY_TYPE)
	case "notify":
		to := &libmessage.To{ /*From: flags.From,*/
				To: []string{ flags.To } }
		notify := libmessage.NewNotify()
		send(flags.Timeout, notify, message, to)
	case "plain":
		fallthrough
	case "tcp":
		if flags.To == "" {
			fmt.Println("To cannot be empty with tcp")
			os.Exit(EMPTY_TO_TCP)
		}

		config := libmessage.NewTcpPlainConfig()

		if !flags.NoSSL {
			tlsConfig := tls.Config{}
			if flags.SSLNoVerify {
				tlsConfig.InsecureSkipVerify = true
			}
			config.SetSSLConfig(tlsConfig)
		}

		tcp := libmessage.NewTcpPlain(config)

		to := tcp.GetTo()
		to.AddAddress(flags.To)

		send(flags.Timeout, tcp, message, to)
	case "androidpn":
		fallthrough
	case "xmppandroidpn":
		if flags.To == EMPTY {
			fmt.Println("To cannot be empty with xmppandroidpn!")
			os.Exit(EMPTY_TO_ANDROIDPN)
		}
		to := &libmessage.To{ /*From: flags.From,*/
					To: []string{ flags.To } }
		android := libmessage.NewXmppAndroidpn(libmessage.NewTcpPlainConfig())
		send(flags.Timeout, android, message, to)
	case "stdout":
		to := &libmessage.To{ /*From: flags.From,*/
				To: []string{ flags.To } }
		stdout := libmessage.NewStdout()
		send(flags.Timeout, stdout, message, to)
	}
}

func send(timeout int,
		messenger libmessage.SendMessageInterface,
		message *libmessage.Message,
		to libmessage.ToInterface) {

	go messenger.Send(message, to)

	select {
	case result := <-messenger.GetResult():
		if result.Result == libmessage.FAILURE {
			fmt.Println("Message could not be send!")
			fmt.Println(result.ErrorString)
			os.Exit(11)
		}
	case <-time.After(time.Duration(timeout)*time.Second):
		fmt.Println("Timeout!")
		os.Exit(12)
	}
}

func validate(flags *flagConfig) {
	if flags.Title == "" || flags.Message == "" {
		flags.usage()
		os.Exit(1)
	}
}

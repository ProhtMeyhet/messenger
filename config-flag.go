package main

import(
	"fmt"
	"bufio"
	"bytes"
	"flag"
	"io/ioutil"
	"os"
)

const EMPTY = ""

type flagConfig struct {
	To, From, Type, Title, Message, Icon string
	Stdin, StdinReadLine, StdinTitle, ShowConfig, Echo, NoSSL, SSLNoVerify bool
	Timeout int
}

func newFlagConfig() *flagConfig {
	return &flagConfig{}
}

func (flags *flagConfig) parse() {
	flag.StringVar(&flags.Title, "title", "", "[mandatory] message title")
	flag.StringVar(&flags.Message, "message", "", "[mandatory] the message")
	//TODO add multiple to support
	flag.StringVar(&flags.To, "to", "", "send to")
	flag.StringVar(&flags.From, "from", "terminal", "who this message is from")
	flag.StringVar(&flags.Type, "type", "notify", "send via this type")
	flag.BoolVar(&flags.NoSSL, "nossl", false, "disable ssl encryption")
	flag.BoolVar(&flags.SSLNoVerify, "no-ssl-verify", false, "disable ssl certificate verification")
	flag.BoolVar(&flags.Stdin, "stdin", false, "read message from stdin till EOF (automatically choosen [eg. for pipe] if --message is not given)")
	flag.BoolVar(&flags.StdinReadLine, "stdin-line", false, "read message from stdin till line break")
	flag.StringVar(&flags.Icon, "icon", "info", "an optional icon to be displayed. see: http://standards.freedesktop.org/icon-naming-spec/icon-naming-spec-latest.html#status")
	flag.IntVar(&flags.Timeout, "timeout", 30, "timeout in seconds")
	// flag.IntVar(&flags.MinHead, "short-head", -1, "for short messengers (libnotify) - use only N first lines of message")
	// flag.IntVar(&flags.ShortTail, "short-tail", -1, "for short messengers (libnotify) - use only N last lines of message")
	flag.BoolVar(&flags.StdinTitle, "stdin-title", false, "first line of stdin is title")
	flag.BoolVar(&flags.Echo, "echo", false, "echo title & message to stdout. true if reading from stdin")
	//flag.StringVar(&flags.Group, "group", "", "send to this group (defined in config)")
	//flag.BoolVar(&flags.ShowConfig, "show-config", false, "display config and exit")

	flag.Parse()

	if flags.Title == "" {
		flags.StdinTitle = true
	}

	if flags.StdinReadLine || flags.StdinTitle {
		flags.Stdin = true
		if !flags.Echo {
			flags.Echo = true
		}
	}

	// stdin
	if (!flags.StdinReadLine && flags.Title != "" && flags.Message == "") ||
		(flags.StdinReadLine && flags.Title == "") ||
		flags.Stdin {
		flags.readFromStdin()

		if !flags.Echo {
			flags.Echo = true
		}
	}
}

func (flags *flagConfig) readFromStdin() {
	reader := bufio.NewReader(os.Stdin)

	var e error
	var stdinTitleByte []byte
	var stdinMessageByte []byte
	if !flags.StdinReadLine {
		stdinMessageByte, e = ioutil.ReadAll(reader)

		if flags.StdinTitle {
			buffer := bytes.NewBuffer(stdinMessageByte)
			reader = bufio.NewReader(buffer)
			stdinTitleByte, e = flags.reallyReadLine(reader)
			stdinTitle := string(stdinTitleByte)
			if e == nil && stdinTitle != "" {
				flags.Title = stdinTitle
			}
			stdinMessageByte, e = ioutil.ReadAll(reader)

		}
	} else {
		fmt.Printf("Title: ")
		stdinTitleByte, e = flags.reallyReadLine(reader)
		stdinTitle := string(stdinTitleByte)
		if e == nil && stdinTitle != "" {
			flags.Title = stdinTitle
		}
		fmt.Printf("Message: ")
		stdinMessageByte, e = flags.reallyReadLine(reader)
	}

	stdinMessage := string(stdinMessageByte)
	if e == nil && stdinMessage != "" {
		flags.Message = stdinMessage
	}
}

func (flags *flagConfig) reallyReadLine(reader *bufio.Reader) ([]byte, error) {
	var isPrefix bool
	var readStdinTmp []byte
	readStdin, isPrefix, e := reader.ReadLine()
	if isPrefix && e == nil {
		for ; isPrefix && e == nil ; {
			readStdinTmp, isPrefix, e = reader.ReadLine()
			readStdin = append(readStdin, readStdinTmp...)
		}
	}
	return readStdin, e
}

func (flags *flagConfig) usage() {
	flag.Usage()
}

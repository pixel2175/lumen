package parser

import (
	"errors"
	"fmt"
	"os"
	"strings"

	. "lumen/src/utils/errs"
	"lumen/src/utils/log"
)

type Action int

const (
	None Action = iota
	Get
	Set
	List
)

type Options struct {
	Action  Action
	Monitor string
	Value   string
}

func PrintHelp() {
	fmt.Println(`usage:
  lumen [-m monitor] get
  lumen [-m monitor] set <value>
  lumen get monitors

examples:
  lumen -m DP-1 set [+10%|20%]
  lumen -m=DP-1 get
  lumen get monitors   # list supported monitors
  lumen get`)
}

func Parse() (opt Options) {
	var err error
	defer Catch(&err, func(err *error) {
		log.Warn((*err).Error())
		PrintHelp()
		os.Exit(1)
	})

	args := os.Args[1:]

	if len(args) == 0 {
		TryE(errors.New("no arguments given"))
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "-h" || arg == "--help":
			PrintHelp()
			os.Exit(0)

		case arg == "get":
			opt.Action = Get
			if i+1 < len(args) && args[i+1] == "monitors" {
				i++
				opt.Action = List
			}

		case arg == "set":
			opt.Action = Set
			if i+1 >= len(args) {
				TryE(errors.New("set needs a value"))
			}
			i++
			opt.Value = args[i]

		case arg == "-m":
			if i+1 >= len(args) {
				TryE(errors.New("-m needs a monitor name"))
			}
			i++
			opt.Monitor = args[i]

		case strings.HasPrefix(arg, "-m="):
			opt.Monitor = strings.TrimPrefix(arg, "-m=")
			if opt.Monitor == "" {
				TryE(errors.New("-m= needs a monitor name"))
			}

		default:
			TryE(errors.New("unknown argument: " + arg))
		}
	}

	if opt.Action == None {
		TryE(errors.New("no action given (get or set)"))
	}

	if opt.Monitor == "" && opt.Action != List {
		log.Info("no monitor given, using default")
	}

	return opt
}

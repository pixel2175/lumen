package actions

import (
	"fmt"

	"lumen/src/backend"
	"lumen/src/parser"
	. "lumen/src/utils/errs"
	"lumen/src/utils/log"
)

func Run(opt parser.Options) {
	var err error
	defer Catch(&err, func(err *error) {
		log.Die(log.Title("Actions"), "Run: %s", (*err).Error())
	})

	if opt.Action == parser.List {
		for i, m := range backend.List() {
			fmt.Printf("%d: %s (%s)\n", i, m.Name, m.Kind)
		}
		return
	}

	b := backend.Load(opt.Monitor)

	switch opt.Action {
	case parser.Get:
		fmt.Println(TryV(b.Get()))

	case parser.Set:
		p := backend.Resolve(b, opt.Value)
		TryE(b.Set(p))
		log.Info(log.Title(b.Name()),"%d%%", p)
	}
}

package backend

import (
	"errors"
	"os"
	"strconv"
	"strings"

	. "lumen/src/utils/errs"
	"lumen/src/utils/log"
)

type Backend interface {
	Name() string
	Get() (int, error)
	Set(percent int) error
}

func Load(monitor string) Backend {
	var err error
	defer Catch(&err, func(err *error) {
		log.Die(log.Title("Backend"), "Load: %s", (*err).Error())
	})

	if monitor == "" || isInternal(monitor) {
		if s, e := loadSysfs(); e == nil {
			return s
		}
	}
	return TryV(loadDDC(monitor))
}

func Resolve(b Backend, in string) int {
	var err error
	defer Catch(&err, func(err *error) {
		log.Die(log.Title("Backend"), "Resolve: %s", (*err).Error())
	})

	rel := strings.HasPrefix(in, "+") || strings.HasPrefix(in, "-")
	v, e := strconv.Atoi(strings.TrimSuffix(in, "%"))
	n := TryV(v, e, errors.New("invalid value: "+in))
	if rel {
		n += TryV(b.Get())
	}
	return min(max(n, 0), 100)
}

func isInternal(m string) bool {
	return strings.HasPrefix(m, "eDP") || strings.HasPrefix(m, "LVDS") || strings.HasPrefix(m, "DSI")
}

func readInt(path string) int {
	data := TryV(os.ReadFile(path))
	return TryV(strconv.Atoi(strings.TrimSpace(string(data))))
}

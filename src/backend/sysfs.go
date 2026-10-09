package backend

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	. "lumen/src/utils/errs"
)

type sysfs struct {
	dir string
	max int
}

func loadSysfs() (b Backend, err error) {
	defer Catch(&err)

	dirs, _ := filepath.Glob("/sys/class/backlight/*")
	if len(dirs) == 0 {
		TryE(errors.New("no backlight device"))
	}
	m := readInt(dirs[0] + "/max_brightness")
	if m <= 0 {
		TryE(errors.New("bad max_brightness"))
	}
	return &sysfs{dir: dirs[0], max: m}, nil
}

func (s *sysfs) Name() string { return "sysfs:" + filepath.Base(s.dir) }

func (s *sysfs) Get() (p int, err error) {
	defer Catch(&err)

	v := readInt(s.dir + "/brightness")
	return (v*100 + s.max/2) / s.max, nil
}

func (s *sysfs) Set(p int) (err error) {
	defer Catch(&err)

	raw := strconv.Itoa(p * s.max / 100)
	TryE(os.WriteFile(s.dir+"/brightness", []byte(raw), 0644))
	return nil
}

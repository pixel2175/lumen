package backend

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	. "lumen/src/utils/errs"
)

type ddc struct {
	f    *os.File
	name string
	max  int
}

func loadDDC(monitor string) (b Backend, err error) {
	defer Catch(&err)

	links, _ := filepath.Glob("/sys/class/drm/card*-*/ddc")
	for _, l := range links {
		conn := filepath.Base(filepath.Dir(l))
		if monitor != "" && !strings.HasSuffix(conn, "-"+monitor) {
			continue
		}
		if d, e := openDDC(l, conn); e == nil {
			return d, nil
		}
	}
	TryE(errors.New("no DDC monitor found (is i2c-dev loaded?)"))
	return nil, nil
}

func openDDC(link, name string) (d *ddc, err error) {
	defer Catch(&err)

	real := TryV(filepath.EvalSymlinks(link))
	f := TryV(os.OpenFile("/dev/"+filepath.Base(real), os.O_RDWR, 0))

	if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), 0x0703, 0x37); e != 0 {
		f.Close()
		TryE(e)
	}

	x := &ddc{f: f, name: name}
	_, m, e := x.vcp()
	if e != nil || m <= 0 {
		f.Close()
		TryE(errors.New("no DDC response"))
	}
	x.max = m
	return x, nil
}

func (d *ddc) vcp() (cur, max int, err error) {
	defer Catch(&err)

	b := []byte{0x51, 0x82, 0x01, 0x10, 0}
	b[4] = 0x6E ^ b[0] ^ b[1] ^ b[2] ^ b[3]
	TryV(d.f.Write(b))
	time.Sleep(40 * time.Millisecond)

	r := make([]byte, 11)
	TryV(d.f.Read(r))
	if r[2] != 0x02 || r[3] != 0 {
		TryE(errors.New("bad DDC reply"))
	}
	max = int(r[6])<<8 | int(r[7])
	cur = int(r[8])<<8 | int(r[9])
	return
}

func (d *ddc) Name() string { return "ddc:" + d.name }

func (d *ddc) Get() (p int, err error) {
	defer Catch(&err)

	cur, max, e := d.vcp()
	TryE(e)
	return (cur*100 + max/2) / max, nil
}

func (d *ddc) Set(p int) (err error) {
	defer Catch(&err)

	v := p * d.max / 100
	b := []byte{0x51, 0x84, 0x03, 0x10, byte(v >> 8), byte(v), 0}
	b[6] = 0x6E ^ b[0] ^ b[1] ^ b[2] ^ b[3] ^ b[4] ^ b[5]
	TryV(d.f.Write(b))
	return nil
}

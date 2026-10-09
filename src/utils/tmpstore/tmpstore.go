package tmpstore

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func tmpFile(name, ext string) string {
	return filepath.Join("/tmp/lumen", strconv.Itoa(os.Getuid())+"-"+name+"."+ext)
}

func ReadTmp(name, ext string) int {
	data, err := os.ReadFile(tmpFile(name, ext))
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return n
}

func WriteTmp(name, ext string, n int) {
	os.MkdirAll("/tmp/lumen", 0777)
	os.Chmod("/tmp/lumen", 0777)
	os.WriteFile(tmpFile(name, ext), []byte(strconv.Itoa(n)), 0644)
}

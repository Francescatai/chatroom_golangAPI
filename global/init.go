package global

import (
	"os"
	"path/filepath"
	"sync"
)

func init() {
	Init()
}

var RootDir string

var once = new(sync.Once)

func Init() {
	once.Do(func() {
		inferRootDir()
		initConfig()
	})
}

var rootDir string

// inferRootDir 推導出項目 root document
func inferRootDir() {
	cwd, err := os.Getwd() //獲得當前工作目錄
	if err != nil {
		panic(err)
	}
	var infer func(d string) string
	infer = func(d string) string {
    // 確認項目根目錄下存在 template
		if exists(d + "/template") {
			return d
		}

		return infer(filepath.Dir(d)) //使用遞迴+Dir()往上層目錄查找
	}

	rootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
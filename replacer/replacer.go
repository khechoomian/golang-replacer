package replacer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//Replace is replacer struct
type Replace struct {
	// Directory that find and replace
	DIR string
	// Old word that want to replace
	Old string
	// New word that want to replace
	New string
}

var DirectoryFined []string

const SplitItem string = "/"
const MergeItem string = "/"
const Base string = "./"

func (rp *Replace) CheckItem(path string, info os.FileInfo) error {
	if info.IsDir() {
		DirectoryFined = append(DirectoryFined, path)
		return nil
	}
	GrepFile([]byte(rp.Old), []byte(rp.New), path)
	rp.FindAndRenameFile(path, info)
	return nil
}
func (rp *Replace) Find() {
	err := filepath.Walk(rp.DIR,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			rp.CheckItem(path, info)
			return nil
		})
	rp.FindAndRenameDIR(DirectoryFined)
	if err != nil {
		log.Println("Err(1)", err)
	}
}

// FindAndRenameFile rename file that in path with info.name
func (rp *Replace) FindAndRenameFile(path string, info os.FileInfo) {
	if strings.Contains(info.Name(), rp.Old) {
		str := strings.Split(path, SplitItem)
		lenstr := len(str)
		filenamereplace := strings.Replace(info.Name(), rp.Old, rp.New, -2)
		str[lenstr-1] = filenamereplace
		NewFileName := MergeString(str)
		fmt.Println("NewFileName: ", NewFileName)
		err := os.Rename(path, NewFileName)
		if err != nil {
			log.Println("Err(2) ", err)
		}
	}
}

// FindAndRenameDIR rename folders recrusivly
func (rp *Replace) FindAndRenameDIR(dir []string) {
	for i := len(dir) - 1; i > 0; i-- {
		if strings.Contains(dir[i], SplitItem) {
			str := strings.Split(dir[i], SplitItem)
			lenstr := len(str)
			if strings.Contains(str[lenstr-1], rp.Old) {
				folderreplace := strings.Replace(str[lenstr-1], rp.Old, rp.New, -2)
				str[lenstr-1] = folderreplace
				NewFileName := MergeString(str)
				err := os.Rename(dir[i], NewFileName)
				if err != nil {
					log.Println("Err(3) ", err)
				}
			}
		} else {
			if strings.Contains(dir[i], rp.Old) {
				myText := strings.Replace(dir[i], rp.Old, rp.New, -2)
				err := os.Rename(dir[i], myText)
				if err != nil {
					log.Println("Err(4) ", err)
				}
			}
		}
	}
}

// GrepFile find string in the file and replace it
func GrepFile(oldb, newb []byte, fn string) (err error) {
	var f *os.File
	if f, err = os.OpenFile(fn, os.O_RDWR, 0); err != nil {
		return
	}
	defer func() {
		if cErr := f.Close(); err == nil {
			err = cErr
		}
	}()
	var b []byte
	if b, err = ioutil.ReadAll(f); err != nil {
		return
	}
	if bytes.Index(b, oldb) < 0 {
		return
	}
	r := bytes.Replace(b, oldb, newb, -2)
	if err = f.Truncate(0); err != nil {
		return
	}
	_, err = f.WriteAt(r, 0)
	return
}

//MergeString merge array of string like Unix/Linux like directory
//If you give []string["A","B","C"] return "A/B/C"
func MergeString(str []string) string {
	var rst string
	for _, s := range str {
		// if rst == "" {
		// 	rst = s
		// 	continue
		// }
		rst = rst + MergeItem + s
	}
	return rst
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	fileList = "FileList.txt"
	prjsPath = "c:/prjs/"
)

var (
	prjNum         string
	numFilesCopied int
)

func main() {
	fmt.Print("prj_bkp")

	if len(os.Args) != 2 {
		log.Fatal("number args: please provide prj num")
	}

	prjNum = os.Args[1]
	PrjDirName, err := getPrjDirName(prjNum)
	if err != nil {
		log.Fatal(err)
	}
	BkpPath := prjsPath + PrjDirName + "/bkps/" + makeBkpDirName()
	err = os.Mkdir(BkpPath, os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}

	fl, err := os.Open(prjsPath + "/" + PrjDirName + "/bkps/" + fileList)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(fl)
	for scanner.Scan() {
		fIn, err := os.Open(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		//fn := filepath.Base(scanner.Text()) todo: remv
		fn := filepath.Base(fIn.Name())
		fOut, err := os.OpenFile(BkpPath+"/"+fn, os.O_RDWR|os.O_CREATE, 666)
		if err != nil {
			log.Fatal(err)
		}
		n, err := io.Copy(fOut, fIn)
		if err != nil {
			log.Fatal(err)
		}
		numFilesCopied++
		fIn.Close()
		fOut.Close()
		fmt.Printf("%s, num bytes copied: %d\n", fn, n)
	}

	fmt.Printf("success: Number files copied: %d", numFilesCopied)
}

func makeBkpDirName() string {
	t := time.Now()
	return t.Format("20060102150405")
}

func getPrjDirName(prjNum string) (string, error) {
	fis, err := ioutil.ReadDir(prjsPath) //[]FileInfo  input is string, directory name
	if err != nil {
		log.Fatal(err)
	}
	count := 0
	prjName := ""
	for _, fi := range fis {
		if fi.IsDir() {
			count++
			prjName = fi.Name()
		}
	}
	if count != 1 {
		return "", errors.New("Prd Dir not found or found multiple: " + prjNum)
	}
	return prjName, nil
}

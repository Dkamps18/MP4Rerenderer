package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//go:embed version
var version string

var total int64
var rendered map[string]int64
var ex string

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("Not enough arguments")
		os.Exit(0)
	}

	e, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flag := flag.NewFlagSet("ConvertToMP4", flag.ExitOnError)
	path := flag.String("p", filepath.Dir(e), "Path (defaults to current working directory)")
	flag.StringVar(&ex, "exec", "ffmpeg", "Overwrite FFmpeg executable")
	ver := flag.Bool("version", false, "Get version info")
	flag.Parse(os.Args[1:])

	if *ver {
		fmt.Println("MP4Rerenderer version "+version, runtime.GOOS+"/"+runtime.GOARCH)
		os.Exit(0)
	}

	ri, err := os.ReadFile("rendered.json")
	if err == nil {
		json.Unmarshal(ri, &rendered)
	} else {
		rendered = make(map[string]int64)
	}
	go func() {
		process(*path)
		fmt.Println("Finished, saved", humanfilesize(total))
		os.Exit(0)
	}()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	fmt.Println("Canceled, saved", humanfilesize(total))
}

func process(dir string) {
	e, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range e {
		n := v.Name()
		if v.IsDir() {
			if n == "." || n == ".." {
				continue
			}
			process(filepath.Join(dir, n))
			continue
		}
		if !strings.HasSuffix(n, ".mp4") {
			continue
		}
		fp := filepath.Join(dir, n)
		if _, ok := rendered[fp]; ok {
			fmt.Println("Ignoring previously rendered file", fp)
			continue
		}
		nfp := filepath.Join(dir, n[:len(n)-4]+"_converted.mp4")
		fmt.Println("Converting", fp)
		_, err := exec.Command("ffmpeg", "-i", fp, nfp).Output()
		if err != nil {
			fmt.Println("Conversion failed")
			continue
		}
		old, _ := os.Stat(fp)
		new, _ := os.Stat(nfp)
		if old.Size() > new.Size() {
			dif := old.Size() - new.Size()
			fmt.Println("Saved", humanfilesize(dif))
			total += dif
			os.Remove(fp)
			os.Rename(nfp, fp)
		} else {
			fmt.Println("File is larger after conversion")
			os.Remove(nfp)
		}
		markredered(fp)
	}
}

func markredered(fp string) {
	rendered[fp] = time.Now().Unix()
	r, _ := json.Marshal(rendered)
	f, err := os.Create("rendered.json")
	if err != nil {
		return
	}
	f.Write(r)
	f.Close()
}

var suffixes = []string{"B", "KB", "MB", "GB", "TB"}

func humanfilesize(size int64) string {
	if size == 0 {
		return "0B"
	}

	base := math.Log(float64(size)) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix)
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/ValleyZw/comgo"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var data [][]string
var (
	flagFile    string
	flagHelp    bool
	flagChannel uint
	flagDetail  bool
	flagVersion bool
)

func init() {
	flag.StringVar(&flagFile, "f", "", "")
	flag.StringVar(&flagFile, "file", "", "")
	flag.BoolVar(&flagHelp, "h", false, "")
	flag.BoolVar(&flagHelp, "help", false, "")
	flag.UintVar(&flagChannel, "c", 1, "")
	flag.UintVar(&flagChannel, "channel", 1, "")
	flag.BoolVar(&flagDetail, "d", false, "")
	flag.BoolVar(&flagDetail, "detail", false, "")
	flag.BoolVar(&flagVersion, "v", false, "")
	flag.BoolVar(&flagVersion, "version", false, "")

	setFlag(flag.CommandLine)
}

func setFlag(flag *flag.FlagSet) {
	flag.Usage = func() {
		Help()
	}
}

func main() {
	args := os.Args[1:]
	arg, err := CommandLine(args)
	CheckError(err)
	flag.CommandLine.Parse(arg)

	if flagVersion {
		Version()
	}

	if flagHelp {
		Help()
	}

	file, err := os.Open(flagFile)
	defer file.Close()
	CheckError(err)

	cfg := comgo.New()
	err = cfg.ReadCFG(file)
	CheckError(err)

	if flagDetail {
		fmt.Println(cfg.GetAnalogChannelNames())
		os.Exit(1)
	}

	var filename string

	name := strings.TrimSuffix(flagFile, filepath.Ext(flagFile))
	if exist := PathExists(name + ".dat"); exist {
		filename = name + ".dat"
	} else if exist := PathExists(name + ".DAT"); exist {
		filename = name + ".DAT"
	} else {
		log.Fatal("Data file File not found.")
	}

	file, err = os.Open(filename)
	defer file.Close()
	CheckError(err)
	err = cfg.ReadDAT(file)
	CheckError(err)

	res, err := cfg.GetAnalogChannelData(uint16(flagChannel))
	CheckError(err)

	ra := cfg.GetSamplingRate()
	ti := cfg.GetStartTime()

	for i := 0; i < cfg.GetSamplingNumber(); i++ {
		s, _ := time.ParseDuration(strconv.FormatFloat(float64(i)/float64(ra), 'g', 1, 64) + "s")
		x := ti.Add(s).Format(AxisFormat)
		y := strconv.FormatFloat(res[i], 'f', -1, 32)
		data = append(data, []string{x, y})
	}

	file, err = os.OpenFile(name+".csv", os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
	CheckError(err)

	w := csv.NewWriter(file)
	w.WriteAll(data) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	log.Println("success!")
}

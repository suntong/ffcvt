package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

// The Encoding struct defines the structure to hold encoding values
type Encoding struct {
	VES string // video encoding method set
	AES string // audio encoding method set
	SES string // subtitle encoding method set
	VEP string // video encoding method prepend
	AEP string // audio encoding method prepend
	SEP string // subtitle encoding method prepend
	VEA string // video encoding method append
	AEA string // audio encoding method append
	ABR string // audio bitrate
	CRF string // the CRF value: 0-51. Higher CRF gives lower quality
	Ext string // extension for the output file
}

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var Defaults map[string]Encoding

func init() {
	initVars() // define flag vars
	debug(Opts.CRF, 3)

	flag.Usage = Usage
	flag.Parse()

	initDefaults() // read & populate defaults from ffcvt.json
	debug(Opts.CRF, 3)

	getDefault() // re-read defaults based on cli -cfg & -t
	//fmt.Fprintf(os.Stderr, "Defaults: '%+v'\n", Defaults)
	debug(Opts.CRF, 3)
	flag.Parse() // update Opts again from cli to overwrite defaults if necessary
	initVals()   // update as per cli settings
}

//go:embed ffcvt.json
// see https://pkg.go.dev/embed
var f embed.FS

func initDefaults() {
	data, err := f.ReadFile("ffcvt.json")
	checkError(err)
	json.Unmarshal(data, &Defaults)
}

func getDefault() {
	if Opts.Cfg != "" {
		data, err := ioutil.ReadFile(Opts.Cfg)
		checkError(err)
		err = json.Unmarshal(data, &Defaults)
		checkError(err)
	}
	if encDefault, ok := Defaults[Opts.Target]; ok {
		// debug(encDefault.Ext, 2)
		Opts.Encoding = encDefault
		// debug(Opts.Encoding.Ext, 2)
		// debug(Opts.Ext, 2)
	} else {
		log.Fatalf("[%s] Error: Wrong target option passed to -t.", progname)
	}

	if Opts.Suffix != "" {
		Opts.Suffix = "_" + Opts.Suffix
	}

}

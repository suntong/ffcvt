package main

import "log"

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

// The Encoding struct defines the structure to hold encoding values
type Encoding struct {
	AES string // audio encoding method set
	VES string // video encoding method set
	AEA string // audio encoding method append
	VEA string // video encoding method append
	ABR string // audio bitrate
	CRF string // the CRF value: 0-51. Higher CRF gives lower quality
}

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var Defaults map[string]Encoding

func init() {

	Defaults = map[string]Encoding{
		"x265-opus": {
			AES: "libopus",
			VES: "libx265",
			ABR: "64k",
			CRF: "28",
		},
		"x264-mp3": {
			AES: "libmp3lame",
			AEA: "-q:a 3",
			VES: "libx264",
			ABR: "256k",
			CRF: "23",
		},
		"youtube": {
			// https://trac.ffmpeg.org/wiki/Encode/YouTube
			// https://trac.ffmpeg.org/wiki/Encode/HighQualityAudio
			AES: "libvorbis",
			AEA: "-q:a 5",
			VES: "libx264",
			VEA: "-pix_fmt yuv420p",
			ABR: "",
			CRF: "20",
		},
	}

}

func getDefault() {
	// preserve command line values
	abr, crf := Opts.ABR, Opts.CRF

	if encDefault, ok := Defaults[Opts.Target]; ok {
		Opts.Encoding = encDefault
	} else {
		log.Fatal(progname + " Error: Wrong target option passed to -t.")
	}

	// restore command line values
	if abr != "" {
		Opts.ABR = abr
	}
	if crf != "" {
		Opts.CRF = crf
	}

	if Opts.Suffix != "" {
		Opts.Suffix = "_" + Opts.Suffix
	}

}

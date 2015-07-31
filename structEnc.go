package main

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
			AEA: "-q:a 2",
			VES: "libx264",
			ABR: "256k",
			CRF: "23",
		},
	}

}

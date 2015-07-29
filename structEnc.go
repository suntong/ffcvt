package main

// The Encoding struct defines the structure to hold encoding values
type Encoding struct {
	AES string // audio encoding method set
	VES string // video encoding method set
	AEA string // audio encoding method append
	VEA string // video encoding method append
	ABR string // audio bitrate
	Crf string // the CRF value: 0-51. Higher CRF gives lower quality
}

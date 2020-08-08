// Command derekify is a sample program to process a image using Derek Bradley's "Adaptive Thresholding using the Integral Image".
//
//    Usage of derekify:
//      -cluster int
//            cluster size (default 100)
//      -format string
//            output format jpg or png (default "jpg")
//      -in string
//            input filename (use stdin if empty)
//      -out string
//            output filename (use stdout if empty)
//      -threshold int
//            threshold [0, 100] (default 85)
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	orientation "github.com/takumakei/exif-orientation"
	"github.com/takumakei/go-derek"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	var input, output, outfmt string
	var clusterSize, threshold int
	flag.StringVar(&input, "in", "", "input filename (use stdin if empty)")
	flag.StringVar(&output, "out", "", "output filename (use stdout if empty)")
	flag.StringVar(&outfmt, "format", "jpg", "output format jpg or png")
	flag.IntVar(&clusterSize, "cluster", 100, "cluster size")
	flag.IntVar(&threshold, "threshold", 85, "threshold [0, 100]")
	flag.Parse()

	// Read file

	var in io.Reader = os.Stdin
	if input != "" {
		f, err := os.Open(input)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	b, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	// Decode image

	r := bytes.NewReader(b)

	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	// Normalize orientation

	r.Reset(b)
	o, _ := orientation.Read(r) // ignore any error. no problem.

	img = orientation.Normalize(img, o)

	// Validation

	if clusterSize < 0 {
		clusterSize = 0
	}

	switch {
	case threshold < 0:
		threshold = 0
	case threshold > 100:
		threshold = 100
	}

	// Processing

	img = derek.Process(img, clusterSize, threshold)

	// Output

	var out io.Writer = os.Stdout
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	switch strings.ToLower(filepath.Ext(output)) {
	case ".png":
		outfmt = "png"
	case ".jpg", ".jpeg":
		outfmt = "jpg"
	}

	if outfmt == "png" {
		return png.Encode(out, img)
	} else {
		return jpeg.Encode(out, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	}
}

package exif

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const TIMEOUT = 60 * time.Second

var exiftool string
var environ []string = os.Environ()

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	exiftool, err = exec.LookPath(fmt.Sprintf("%s/exiftool/exiftool", dir))
	if err != nil {
		log.Fatal(err)
	}
}

func Exec(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	log.Printf("Runing exiftool %v", args)
	if err := exec.CommandContext(ctx, exiftool, args...).Run(); err != nil {
		return err
	}
	return nil
}

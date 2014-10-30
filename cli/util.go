package cli

import (
	"bufio"
	"fmt"
	"github.com/LapisBlue/Tar/skin"
	"github.com/LapisBlue/Tar/util"
	"image"
	"image/png"
	"io"
	"os"
	"strings"
)

func printError(err error, description ...interface{}) int {
	fmt.Fprintln(os.Stderr, description...)
	fmt.Fprintln(os.Stderr, err)
	return 1
}

func readFrom(source string, args []string) []string {
	switch source {
	case "ARGS", "args":
		return args
	case "STDIN", "stdin":
		// TODO: Let this start the generation once the first line is read
		lines, err := readLines(os.Stdin)
		if err != nil {
			printError(err, "Failed to read from STDIN")
			return nil
		}

		return lines
	default:
		lines, err := readFile(source)
		if err != nil {
			printError(err, "Unable to read from %s: %s\n", source, err)
			return nil
		}

		return lines
	}
}

func isStdout(out string) bool {
	return out == "STDOUT" || out == "stdout"
}

func readLines(reader io.Reader) (lines []string, err error) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) > 0 {
			lines = append(lines, scanner.Text())
		}
	}

	err = scanner.Err()
	return
}

func readFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	return readLines(file)
}

func downloadSkins(players []string) (result []*skin.Skin) {
	fmt.Printf("Downloading %d skin(s), please wait...\n", len(players))

	watch := util.GlobalWatch().Mark()
	result = make([]*skin.Skin, len(players))

	var err error
	for i, player := range players {
		watch.Mark()
		result[i], err = skin.Download(player)
		if err != nil {
			printError(err, "Failed to download skin:", player, watch)
			continue
		}

		fmt.Println("Downloaded skin:", player, watch)
	}

	fmt.Println("Finished downloading skins", watch)
	return
}

func saveResults(players []string, results []image.Image, dest string) {
	format := strings.Contains(dest, "%s")
	fmt.Printf("Saving %d image(s), please wait...\n", len(results))

	watch := util.GlobalWatch().Mark()
	for i, player := range players {
		watch.Mark()

		result := results[i]
		if result == nil {
			continue
		}

		name := player
		if format {
			name = fmt.Sprintf(dest, name)
		} else {
			name = dest
		}

		file, err := os.Create(name)
		if err != nil {
			printError(err, "Failed to open file: ", name, watch)
			continue
		}

		err = png.Encode(file, result)
		if err != nil {
			printError(err, "Failed to write image to file:", name, watch)
			continue
		}

		fmt.Println("Saved image:", player, watch)
	}

	fmt.Println("Finished saving images", watch)
	return
}

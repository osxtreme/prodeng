package main

// TODO:
// -Support more complex interleave of multiple hits besides line color?
// -Fix Color priority order hard coded for lines?
// -support more colors and hi-colors and bold

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type arrayFlags []string

// Implements needed method for flag.Value (String) - used as default value
// or at least in help text it does
func (s *arrayFlags) String() string {
	// Empty string for default
	return ""
}

func (s *arrayFlags) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {

	var ansi_red = "\033[31m"
	var ansi_yellow = "\033[33m"

	var ansi_normal = "\033[0m"

	var flag_yellow_lines arrayFlags
	flag.Var(&flag_yellow_lines, "yl", "String(s) to trigger whole line in yellow")

	var flag_yellow_strings arrayFlags
	flag.Var(&flag_yellow_strings, "ys", "String(s) to make yellow")

	var flag_red_lines arrayFlags
	flag.Var(&flag_red_lines, "rl", "String(s) to trigger whole line in red")

	var flag_red_strings arrayFlags
	flag.Var(&flag_red_strings, "rs", "String(s) to make red")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tail -f <log> | grep stuff | %s [flags]\n",
			filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		// For multiple colors: If whole line is modified, don't print ansi_normal, revert to line
		line_color := ansi_normal

		for _, yl := range flag_yellow_lines {
			if strings.Contains(line, yl) {
				line = ansi_yellow + line + ansi_normal
				line_color = ansi_yellow
			}
		}

		for _, rl := range flag_red_lines {
			if strings.Contains(line, rl) {
				line = ansi_red + line + ansi_normal
				line_color = ansi_red
			}
		}

		for _, ys := range flag_yellow_strings {
			//fmt.Println("Looking for:", ys)
			// Faster or slower to do strings.Contains check first?
			line = strings.Replace(line, ys, ansi_yellow+ys+line_color, -1)
		}

		for _, rs := range flag_red_strings {
			// Faster or slower to do strings.Contains check first?
			line = strings.Replace(line, rs, ansi_red+rs+line_color, -1)
		}

		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

}

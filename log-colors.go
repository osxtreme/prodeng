package main

// TODO:
// -Support more complex interleave of multiple hits besides line color?
// -Fix Color priority order hard coded for lines?
// -support more colors and hi-colors and bold
// -support -s or -v to do grep -v like behavior (supress/skip line)?

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

func reg_aFlag(name string, desc string) *arrayFlags {
	var af arrayFlags
	flag.Var(&af, name, desc)
	return &af
}

func main() {

	// http://pueblo.sourceforge.net/doc/manual/ansi_color_codes.html
	// non-bold color codes for the lines include turning off bold, so that mid-line strings
	// can reset from bold to standard (or vice versa) as needed
	var ansi_red = "\033[0;31m"
	//	var ansi_bold_red = "\033[1;31m"

	var ansi_yellow = "\033[0;33m"
	var ansi_bold_yellow = "\033[1;33m"

	var ansi_normal = "\033[0m"

	flag_yellow_lines := reg_aFlag("yl", "String(s) to trigger whole line in yellow")
	flag_yellow_strings := reg_aFlag("ys", "String(s) to make yellow")
	flag_byellow_strings := reg_aFlag("bys", "String(s) to make bold yellow")
	flag_red_lines := reg_aFlag("rl", "String(s) to trigger whole line in red")
	flag_red_strings := reg_aFlag("rs", "String(s) to make red")

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

		for _, yl := range *flag_yellow_lines {
			if strings.Contains(line, yl) {
				line = ansi_yellow + line + ansi_normal
				line_color = ansi_yellow
			}
		}

		for _, rl := range *flag_red_lines {
			if strings.Contains(line, rl) {
				line = ansi_red + line + ansi_normal
				line_color = ansi_red
			}
		}

		for _, s := range *flag_yellow_strings {
			//fmt.Println("Looking for:", ys)
			// Faster or slower to do strings.Contains check first?
			line = strings.Replace(line, s, ansi_yellow+s+line_color, -1)
		}

		for _, s := range *flag_byellow_strings {
			// Faster or slower to do strings.Contains check first?
			line = strings.Replace(line, s, ansi_bold_yellow+s+line_color, -1)
		}

		for _, s := range *flag_red_strings {
			// Faster or slower to do strings.Contains check first?
			line = strings.Replace(line, s, ansi_red+s+line_color, -1)
		}

		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

}

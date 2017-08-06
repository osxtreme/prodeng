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
	return "[string]"
}

func (s *arrayFlags) Set(value string) error {
	*s = append(*s, value)
	return nil
}

const (
	// http://pueblo.sourceforge.net/doc/manual/ansi_color_codes.html
	// non-bold color codes for the lines include turning off bold, so that mid-line strings
	// can reset from bold to standard (or vice versa) as needed

	ansi_bold        = "\033[1m"
	ansi_red         = "\033[0;31m"
	ansi_bold_red    = "\033[1;31m"
	ansi_yellow      = "\033[0;33m"
	ansi_bold_yellow = "\033[1;33m"

	ansi_reset = "\033[0m"
)

type Triggers struct {
	triggers      map[*arrayFlags]string
	trigger_order []*arrayFlags
}

func NewTriggers() *Triggers {
	t := new(Triggers)
	t.triggers = make(map[*arrayFlags]string)
	return t
}

func register_trigger(t *Triggers, style string, flag_name string, code string) {
	var af arrayFlags
	flag.Var(&af, flag_name, "String to trigger "+style)
	t.triggers[&af] = code
	// Store the order that we received the trigger in for priority
	// and because range over maps is always randomized
	t.trigger_order = append(t.trigger_order, &af)
}

func main() {

	lines := NewTriggers()
	items := NewTriggers()

	register_trigger(lines, "yellow line", "yl", ansi_yellow)
	register_trigger(items, "yellow item", "ys", ansi_yellow)
	register_trigger(lines, "bold yellow line", "byl", ansi_bold_yellow)
	register_trigger(items, "bold yellow item", "bys", ansi_bold_yellow)
	register_trigger(lines, "red line", "rl", ansi_red)
	register_trigger(items, "red item", "rs", ansi_red)
	register_trigger(lines, "bold red line", "brl", ansi_bold_red)
	register_trigger(items, "bold red item", "brs", ansi_bold_red)
	register_trigger(lines, "bold line", "bl", ansi_bold)
	register_trigger(items, "bold item", "bs", ansi_bold)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tail -f <log> | grep stuff | %s [flags]\n",
			filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		// Save for strings in multiple colors: If whole line is modified, don't print ansi_reset, revert to line
		line_color := ansi_reset

		for _, line_trigger_strings := range lines.trigger_order {
			line_style_code := lines.triggers[line_trigger_strings]
			for _, line_trigger_string := range *line_trigger_strings {
				if strings.Contains(line, line_trigger_string) {
					line = line_style_code + line + ansi_reset
					line_color = line_style_code
				}
			}
		}

		for _, item_trigger_strings := range items.trigger_order {
			item_style_code := items.triggers[item_trigger_strings]
			for _, item_trigger_string := range *item_trigger_strings {
				// Faster or slower to do strings.Contains check first?
				line = strings.Replace(line, item_trigger_string, item_style_code+item_trigger_string+line_color, -1)
			}
		}

		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

}

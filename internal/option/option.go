package option

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Options struct {
	Exclude string // -x or --exclude
	Version bool // -v or --version
	Help    bool // -h or --help
}

var flagUsages = map[string]string {
	"x": "Exclude files matching the pattern",
	"v": "Show version information",
	"h": "Show help information",
}

func Parse(args []string) (*Options, error) {
	opts := &Options{}

	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)

	flagSet.StringVar(&opts.Exclude, "x", "", "x")
	flagSet.StringVar(&opts.Exclude, "exclude", "", "x")
	flagSet.BoolVar(&opts.Version, "v", false, "v")
	flagSet.BoolVar(&opts.Version, "version", false, "v")
	flagSet.BoolVar(&opts.Help, "h", false, "h")
	flagSet.BoolVar(&opts.Help, "help", false, "h")

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", flagSet.Name())
		printFlags(flagSet)
	}

	err := flagSet.Parse(args[1:])
	if err != nil {
		flag.Usage()
		return nil, err
	}

	return opts, nil
}

// -x, --excludeのように複数の名前を持つオプションをまとめて表示する
func printFlags(flagSet *flag.FlagSet) {
	type flagInfo struct {
		names []string
		usage string
	}
	flagMap := make(map[string]*flagInfo)

	flagSet.VisitAll(func(f *flag.Flag) {
        key := f.Usage
        if fi, exists := flagMap[key]; exists {
            fi.names = append(fi.names, "-" + f.Name)
        } else {
            flagMap[key] = &flagInfo{
                names: []string{"-" + f.Name},
                usage: flagUsages[f.Usage],
            }
        }
    })

    var keys []string
    for k := range flagMap {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    for _, key := range keys {
        fi := flagMap[key]
		sort.Slice(fi.names, func(i, j int) bool {
            return len(fi.names[i]) < len(fi.names[j])
        })
        fmt.Fprintf(os.Stderr, "  %s\n        %s\n", strings.Join(fi.names, ", "), fi.usage)
    }
}


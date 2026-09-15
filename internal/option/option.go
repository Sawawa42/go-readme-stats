package option

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
)

type Options struct {
	Include         string   // -i or --include
	IncludePatterns []string // -i or --includeで指定されたパターンをカンマで分割して格納
	Role            string   // -r or --role
	Roles           []string // -r or --roleをカンマで分割し大文字化したもの(未指定ならOWNER)
	ExcludeRepo     string   // -e or --exclude-repo
	ExcludeRepos    []string // -e or --exclude-repoをカンマで分割して格納
	SizeWeight      float64  // --size-weight
	CountWeight     float64  // --count-weight
	Verbose         bool     // -v or --verbose
	Help            bool     // -h or --help
	FlagSet         *flag.FlagSet
}

var flagUsages = map[string]string{
	"i": "Include only languages matching the pattern",
	"r": "Repository roles to include: OWNER, ORGANIZATION_MEMBER, COLLABORATOR (default: OWNER)",
	"e": "Exclude repositories (\"name\" or \"owner/name\")",
	"sw": "Weight of byte size in the usage index (default: 0.5)",
	"cw": "Weight of repository count in the usage index (default: 0.5)",
	"v": "Show per-repository breakdown (may expose private repository names)",
	"h": "Show help information",
}

func Parse(args []string) (*Options, error) {
	opts := &Options{}

	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)

	flagSet.StringVar(&opts.Include, "i", "", "i")
	flagSet.StringVar(&opts.Include, "include", "", "i")
	flagSet.StringVar(&opts.Role, "r", "", "r")
	flagSet.StringVar(&opts.Role, "role", "", "r")
	flagSet.StringVar(&opts.ExcludeRepo, "e", "", "e")
	flagSet.StringVar(&opts.ExcludeRepo, "exclude-repo", "", "e")
	flagSet.Float64Var(&opts.SizeWeight, "size-weight", 0.5, "sw")
	flagSet.Float64Var(&opts.CountWeight, "count-weight", 0.5, "cw")
	flagSet.BoolVar(&opts.Verbose, "v", false, "v")
	flagSet.BoolVar(&opts.Verbose, "verbose", false, "v")
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

	if opts.Include != "" {
		opts.IncludePatterns = strings.Split(opts.Include, ",")
	}
	if opts.ExcludeRepo != "" {
		opts.ExcludeRepos = strings.Split(opts.ExcludeRepo, ",")
	}
	opts.Roles, err = parseRoles(opts.Role)
	if err != nil {
		return nil, err
	}
	if opts.SizeWeight < 0 || opts.CountWeight < 0 {
		return nil, fmt.Errorf("--size-weight and --count-weight must be non-negative")
	}
	opts.FlagSet = flagSet

	return opts, nil
}

var validRoles = []string{"OWNER", "ORGANIZATION_MEMBER", "COLLABORATOR"}

func parseRoles(role string) ([]string, error) {
	if role == "" {
		return []string{"OWNER"}, nil
	}
	var roles []string
	for r := range strings.SplitSeq(role, ",") {
		r = strings.ToUpper(strings.TrimSpace(r))
		if !slices.Contains(validRoles, r) {
			return nil, fmt.Errorf("invalid role %q (valid: %s)", r, strings.Join(validRoles, ", "))
		}
		roles = append(roles, r)
	}
	return roles, nil
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
			fi.names = append(fi.names, "-"+f.Name)
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

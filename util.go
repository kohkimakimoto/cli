package cli

import (
	"strings"
)

func GetStringFlagValue(args []string, flags ...string) string {
	if len(args) != 0 {
		for i, s := range args {
			if s == "--" {
				break
			}
			for _, flag := range flags {
				if len(s) >= len(flag) && s[0:len(flag)] == flag {
					if len(s) >= len(flag + "=") && s[0:len(flag + "=")] == flag + "=" {
						return strings.Split(s, "=")[1]
					} else if s == flag {
						return args[i + 1]
					}
				}
			}

		}
	}

	return ""
}

func GetBoolFlagValue(args []string, flags ...string) bool {
	if len(args) != 0 {
		for _, s := range args {
			if s == "--" {
				break
			}
			for _, flag := range flags {
				if s == flag {
					return true
				}
			}
		}
	}

	return false
}

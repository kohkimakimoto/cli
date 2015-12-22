package cli

import(
	"testing"
)

func TestGetStringFlagValue(t *testing.T) {

	args := []string{"greet", "--name", "Jeremy"}
	name := GetStringFlagValue(args, "--name")
	if name != "Jeremy" {
		t.Error(name)
	}

	args = []string{"greet", "--name=Jeremy"}
	name = GetStringFlagValue(args, "--name")
	if name != "Jeremy" {
		t.Error(name)
	}

	args = []string{"greet", "-n=Jeremy"}
	name = GetStringFlagValue(args, "-n")
	if name != "Jeremy" {
		t.Error(name)
	}

	args = []string{"greet", "-n=Jeremy"}
	name = GetStringFlagValue(args, "--name", "-n")
	if name != "Jeremy" {
		t.Error(name)
	}

	args = []string{"greet", "--name=Jeremy"}
	name = GetStringFlagValue(args, "--name", "-n")
	if name != "Jeremy" {
		t.Error(name)
	}

	args = []string{"greet", "--", "--name=Jeremy"}
	name = GetStringFlagValue(args, "--name", "-n")
	if name != "" {
		t.Error(name)
	}
}

func TestGetBoolFlagValue(t *testing.T) {
	args := []string{"greet", "--abc", "Jeremy"}
	b := GetBoolFlagValue(args, "--abc")
	if !b {
		t.Error(b)
	}

	args = []string{"greet", "-a", "Jeremy"}
	b = GetBoolFlagValue(args, "--abc", "-a")
	if !b {
		t.Error(b)
	}
}

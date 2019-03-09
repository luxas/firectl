package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DaemonBackend interface {
	RunDetached(args []string) (id string, err error)
}

type DockerBackend struct{}

func (db *DockerBackend) RunDetached(args []string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	execpath := filepath.Join(wd, "firectl")
	argMap, err := parseArgs(args)
	if err != nil {
		return "", err
	}
	kernel := argMap["kernel"]
	if len(kernel) == 0 {
		kernel = filepath.Join(wd, "vmlinux")
	}
	rootdrive := argMap["root-drive"]
	if !filepath.IsAbs(rootdrive) {
		rootdrive = filepath.Join(wd, rootdrive)
	}
	allArgs := []string{
		"run",
		"-itd",
		"--privileged",
		"--device=/dev/kvm",
		"-v=/var/lib/firecracker:/var/lib/firecracker",
		fmt.Sprintf("-v=%s:/host-firectl", execpath),
		fmt.Sprintf("-v=%s:%s", kernel, kernel),
		fmt.Sprintf("-v=%s:%s", rootdrive, rootdrive),
		"--name fc$(date +%s)",
		"luxas/firectl",
		"/host-firectl",
		"--containerized",
	}
	for _, arg := range args {
		if arg == "-d" || arg == "--daemon" {
			continue
		}
		allArgs = append(allArgs, arg)
	}
	out, err := exec.Command("docker", allArgs).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing docker with arguments %v. output: %q, error: %v", allArgs, string(out), err)
	}
	return "", nil
}

func parseArgs(arguments []string) (map[string]string, error) {
	resultingMap := map[string]string{}
	for i, arg := range arguments {
		key, val, err := parseArg(arg)
		if err != nil && i != 0 {
			return fmt.Errorf("argument %q could not be parsed correctly", arg)
		}
		resultingMap[key] = val
	}
	return resultingMap
}

func parseArg(arg string) (string, string, error) {
	if strings.HasPrefix(arg, "--") {
		arg = strings.TrimPrefix(arg, "-")
	} else {
		if strings.HasPrefix(arg, "-") {
			arg = strings.TrimPrefix(arg, "-")
		} else {
			return "", "", errors.New("the argument should start with '--' or '-'")
		}
	}
	if !strings.Contains(arg, "=") {
		return "", "", errors.New("the argument should have a '=' between the flag and the value")
	}
	// Split the string on =. Return only two substrings, since we want only key/value, but the value can include '=' as well
	kv := strings.SplitN(arg, "=", 2)
	// Make sure both a key and value is present
	if len(kv) != 2 {
		return "", "", errors.New("the argument must have both a key and a value")
	}
	if len(kv[0]) == 0 {
		return "", "", errors.New("the argument must have a key")
	}
	return kv[0], kv[1], nil
}

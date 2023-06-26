package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os/exec"
	"strings"
)

// Check containerd version
func containerdVersion() (string, map[string]string) {
	nv, err := exec.Command("containerd", "--version").Output()
	if err != nil {
		log.Fatal(err)
	}
	v := strings.TrimSuffix(string(nv), "\n")
	// containerd github.com/containerd/containerd Version GitCommit
	c := strings.SplitN(v, " ", 4)
	if len(c) == 4 && c[0] == "containerd" {
		v = strings.Replace(c[2], "v", "", 1)
		if c[3] != "" {
			return v, map[string]string{"GitCommit": c[3]}
		}
	}
	return v, nil
}


// List all the created containers
func nerdctlContainers(all bool) []map[string]interface{} {
	args := []string{"ps"}
	if all {
		args = append(args, "-a")
	}
	args = append(args, "--format", "{{json .}}")
	nc, err := exec.Command("nerdctl", args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	var containers []map[string]interface{}
	scanner := bufio.NewScanner(bytes.NewReader(nc))
	for scanner.Scan() {
		var container map[string]interface{}
		err = json.Unmarshal(scanner.Bytes(), &container)
		if err != nil {
			log.Fatal(err)
		}
		containers = append(containers, container)
	}
	return containers
}

// Inspect container
func nerdctlContainer(name string) (map[string]interface{}, error) {
	args := []string{"container", "inspect", "--mode", "dockercompat"}
	args = append(args, name, "--format", "{{json .}}")
	nc, err := exec.Command("nerdctl", args...).Output()
	if err != nil {
		return nil, err
	}
	var image map[string]interface{}
	err = json.Unmarshal(nc, &image)
	if err != nil {
		log.Fatal(err)
	}
	return image, nil
}


func nerdctlContainerInspect(id string) []map[string]interface{} {
	nc, err := exec.Command("nerdctl", "container", "inspect", id).Output()
	if err != nil {
		log.Fatal(err)
	}
	var inspect []map[string]interface{}
	err = json.Unmarshal(nc, &inspect)
	if err != nil {
		log.Fatal(err)
	}
	return inspect
}


func nerdctlLogs(id string, w io.Writer) error {
	cmd := exec.Command("nerdctl", "logs", "-f", id)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	containers[id] = NerdctlContainer{
		ID:  id,
		Cmd: cmd,
		Out: stdout,
		Err: stderr,
	}
	return nil
}


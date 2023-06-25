package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

func nerdctlVolumes(filter string) []map[string]interface{} {
	args := []string{"volume", "ls"}
	if filter != "" {
		args = append(args, "--filter", filter)
	}
	args = append(args, "--format", "{{json .}}")
	nc, err := exec.Command("nerdctl", args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	var volumes []map[string]interface{}
	scanner := bufio.NewScanner(bytes.NewReader(nc))
	for scanner.Scan() {
		var volume map[string]interface{}
		err = json.Unmarshal(scanner.Bytes(), &volume)
		if err != nil {
			log.Fatal(err)
		}
		volumes = append(volumes, volume)
	}
	return volumes
}

func nerdctlVolume(name string) (map[string]interface{}, error) {
	args := []string{"volume", "inspect"}
	args = append(args, name, "--format", "{{json .}}")
	nc, err := exec.Command("nerdctl", args...).Output()
	if err != nil {
		return nil, err
	}
	var volume map[string]interface{}
	err = json.Unmarshal(nc, &volume)
	if err != nil {
		log.Fatal(err)
	}
	return volume, nil
}

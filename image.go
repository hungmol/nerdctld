package nerdctld

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

func nerdctlImages(filter string) []map[string]interface{} {
	args := []string{"images"}
	if filter != "" {
		args = append(args, filter)
	}
	args = append(args, "--format", "{{json .}}")
	nc, err := exec.Command("nerdctl", args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	var images []map[string]interface{}
	scanner := bufio.NewScanner(bytes.NewReader(nc))
	for scanner.Scan() {
		var image map[string]interface{}
		err = json.Unmarshal(scanner.Bytes(), &image)
		if err != nil {
			log.Fatal(err)
		}
		images = append(images, image)
	}
	return images
}

func nerdctlImage(name string) (map[string]interface{}, error) {
	args := []string{"image", "inspect", "--mode", "dockercompat"}
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

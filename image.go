package nerdctld

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os/exec"
	"strings"
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

// Remove image nerdctl rmi
func nerdctlRmi(name string, w io.Writer) error {
	args := []string{"rmi"}
	args = append(args, name)
	nc, err := exec.Command("nerdctl", args...).Output()
	if err != nil {
		return err
	}
	removed := []map[string]string{}
	lines := strings.Split(string(nc), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Untagged: ") {
			image := strings.Replace(line, "Untagged: ", "", 1)
			removed = append(removed, map[string]string{"Untagged": image})
		} else if strings.HasPrefix(line, "Deleted:") {
			image := strings.Replace(line, "Deleted: ", "", 1)
			removed = append(removed, map[string]string{"Deleted": image})
		}
	}
	d, _ := json.Marshal(removed)
	_, err = w.Write(d)
	if err != nil {
		return err
	}
	return nil
}

// Save docker image
func nerdctlSave(names []string, w io.Writer) error {
	args := []string{"save"}
	args = append(args, names...)
	cmd := exec.Command("nerdctl", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	errors := make(chan error)
	go func() {
		defer stdout.Close()
		if _, err := io.Copy(w, stdout); err != nil {
			errors <- err
		}
		errors <- nil
	}()
	err = cmd.Run()
	if err != nil {
		return err
	}
	if err := <-errors; err != nil {
		return err
	}
	return nil
}

// Push image
func nerdctlPush(name string, w io.Writer) error {
	args := []string{"push"}
	nc, err := exec.Command("nerdctl", args...).Output()
	if err != nil {
		return err
	}
	lines := strings.Split(string(nc), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		data := map[string]string{"stream": line + "\n"}
		l, _ := json.Marshal(data)
		_, err = w.Write(l)
		if err != nil {
			return err
		}
		_, err = w.Write([]byte{'\n'})
		if err != nil {
			return err
		}
	}
	return nil
}

// Load image
func nerdctlLoad(quiet bool, r io.Reader, w io.Writer) error {
	args := []string{"load"}
	cmd := exec.Command("nerdctl", args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	errors := make(chan error)
	go func() {
		defer stdin.Close()
		if _, err := io.Copy(stdin, r); err != nil {
			errors <- err
		}
		errors <- nil
	}()
	nc, err := cmd.Output()
	if err != nil {
		return err
	}
	if err := <-errors; err != nil {
		return err
	}
	lines := strings.Split(string(nc), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		data := map[string]string{"stream": line + "\n"}
		l, _ := json.Marshal(data)
		_, err = w.Write(l)
		if err != nil {
			return err
		}
		_, err = w.Write([]byte{'\n'})
		if err != nil {
			return err
		}
	}
	return nil
}

// Pull image
func nerdctlPull(name string, w io.Writer) error {
	args := []string{"pull"}
	args = append(args, name)
	nc, err := exec.Command("nerdctl", args...).Output()
	if err != nil {
		return err
	}
	lines := strings.Split(string(nc), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		data := map[string]string{"stream": line + "\n"}
		l, _ := json.Marshal(data)
		_, err = w.Write(l)
		if err != nil {
			return err
		}
		_, err = w.Write([]byte{'\n'})
		if err != nil {
			return err
		}
	}
	return nil
}

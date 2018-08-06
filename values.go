package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"io"
)

const (
	defaultDirectoryPermission = 0755
	defaultValuesFilename      = "values.yaml"
)

type valuesCmd struct {
	chartPath    string
	values       valueFiles
	outputDir    string
	backupSuffix string
}

func (cmd *valuesCmd) run() error {
	cmd.values.prepend(path.Join(cmd.chartPath, defaultValuesFilename))
	vv, err := vals(cmd.values)
	if err != nil {
		return err
	}
	if cmd.outputDir != "" {
		return writeToFile(cmd.outputDir, defaultValuesFilename, cmd.backupSuffix, vv)
	} else {
		fmt.Println(string(vv))
	}
	return nil
}

// write the <data> to <output-dir>/<name>, backup <name> to <name><backup-suffix> first if exist
func writeToFile(outputDir string, name string, backupSuffux string, data []byte) error {
	outfileName := strings.Join([]string{outputDir, name}, string(filepath.Separator))

	err := ensureDirectoryForFile(outfileName)
	if err != nil {
		return err
	}

	err = ensureFileBackedup(outfileName, backupSuffux)
	if err != nil {
		return err
	}

	f, err := os.Create(outfileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	fmt.Printf("wrote %s\n", outfileName)
	return nil
}

// check if the file backed up
func ensureFileBackedup(file string, backupSuffix string) error {
	if _, err := os.Stat(file); err == nil {
		backupPath := path.Join(path.Dir(file), filepath.Base(file)+backupSuffix)
		ensureFileBackedup(backupPath, backupSuffix)
		copy(file, backupPath)
		fmt.Printf("backed up %s to %s\n", file, backupPath)
	}
	return nil
}

func copy(src, dest string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}

// check if the directory exists to create file. creates if don't exists
func ensureDirectoryForFile(file string) error {
	baseDir := path.Dir(file)
	_, err := os.Stat(baseDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return os.MkdirAll(baseDir, defaultDirectoryPermission)
}

// liberally borrows from Helm
func vals(values valueFiles) ([]byte, error) {
	base := map[string]interface{}{}

	// User specified a values files via -f/--values
	for _, filePath := range values {
		currentMap := map[string]interface{}{}
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return []byte{}, err
		}

		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return []byte{}, fmt.Errorf("failed to parse %s: %s", filePath, err)
		}
		// Merge with the previous map
		base = mergeValues(base, currentMap)
	}

	return yaml.Marshal(base)
}

// Copied from Helm.

func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = nextMap
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}

type valueFiles []string

func (v *valueFiles) String() string {
	return fmt.Sprint(*v)
}

func (v *valueFiles) Type() string {
	return "valueFiles"
}

func (v *valueFiles) Set(value string) error {
	for _, filePath := range strings.Split(value, ",") {
		*v = append(*v, filePath)
	}
	return nil
}

func (v *valueFiles) prepend(value string) error {
	for _, filePath := range strings.Split(value, ",") {
		*v = append(valueFiles{filePath}, *v...)
	}
	return nil
}

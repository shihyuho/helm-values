package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"path"
)

var (
	values  valueFiles
	output  string
	verbose bool
)

func main() {
	cmd := &cobra.Command{
		Use:   "helm values [flags] CHART",
		Short: fmt.Sprintf("merge one or more YAML files of values"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("chart is required")
			}
			return run(args[0])
		},
	}

	f := cmd.Flags()
	f.VarP(&values, "values", "f", "specify one or more YAML files of values")
	f.StringVarP(&output, "output", "o", "", "write to a file, instead of stdout")
	f.BoolVarP(&verbose, "verbose", "v", false, "enable verbose output.")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(chart string) error {
	values.prepend(path.Join(chart, "values.yaml"))
	if verbose {
		fmt.Println("merging values:", values.String())
	}
	vv, err := vals()
	if err != nil {
		return err
	}
	if output != "" {
		return ioutil.WriteFile(output, vv, os.ModePerm)
	}
	fmt.Println(string(vv))
	return nil
}

// liberally borrows from Helm
func vals() ([]byte, error) {
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

package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
)

type spec struct {
	Replicas             int                    `yaml:"replicas"`
	RevisionHistoryLimit int                    `yaml:"revisionHistoryLimit"`
	MinReadySeconds      int                    `yaml:"minReadySeconds"`
	Strategy             map[string]interface{} `yaml:"strategy"`
	Template             struct {
		Metadata struct {
			Labels      map[string]string      `yaml:"labels"`
			Annotations map[string]interface{} `yaml:"annotations"`
		} `yaml:"metadata"`
		Spec map[string]interface{}
	} `yaml:"template"`
}

type deployment struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   map[string]interface{} `yaml:"metadata"`
	Spec       spec                   `yaml:"spec"`
}

func transformDeployment(text string, graceful bool) {
	requiredType := "Deployment"
	annotation := "pod.beta.kubernetes.io/init-containers"
	var d deployment

	mErr := yaml.Unmarshal([]byte(text), &d)
	if mErr != nil {
		panic(mErr)
	}

	if d.Kind != requiredType {
		panic(fmt.Sprintf("This is not a %s. Got: %s", requiredType, d.Kind))
	}

	containerAnnotation := d.Spec.Template.Metadata.Annotations[annotation]
	if containerAnnotation == nil && graceful {
		fmt.Println(text)
		return
	} else if containerAnnotation == nil {
		panic(fmt.Sprintf("Missing annotation '%s' on template metadata", annotation))
	}

	// Parse annotation JSON
	var i []interface{}
	json.Unmarshal([]byte(containerAnnotation.(string)), &i)

	initContainers := make([]interface{}, 0)
	initContainers = append(initContainers, containerAnnotation)
	d.Spec.Template.Spec["initContainers"] = i

	result, err := yaml.Marshal(d)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(result))
}

var graceful bool

var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "deployment migrates v1.5-1.7 initContainer annotations to native spec values.",
	Long:  "deployment migrates v1.5-1.7 initContainer annotations to native spec values.",
	Run: func(cmd *cobra.Command, args []string) {
		fi, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}
		if fi.Size() < 1 {
			panic(fmt.Errorf("stdin is empty"))
		}

		scanner := bufio.NewScanner(os.Stdin)
		var inputBuffer bytes.Buffer
		for scanner.Scan() {
			inputBuffer.WriteString(scanner.Text() + "\n")
		}

		text := inputBuffer.String()
		transformDeployment(text, graceful)
	},
}

func init() {
	deploymentCmd.PersistentFlags().BoolVarP(&graceful, "graceful", "g", false, "fail gracefully")
	rootCmd.AddCommand(deploymentCmd)
}

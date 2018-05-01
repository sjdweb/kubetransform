package cmd

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
)

type secret struct {
	APIVersion string `yaml:"apiVersion"`
	Data       map[string]string
	Kind       string
}

func transformSecret(text string) {
	requiredType := "Secret"
	var myStuff secret

	mErr := yaml.Unmarshal([]byte(text), &myStuff)
	if mErr != nil {
		panic(mErr)
	}

	if myStuff.Kind != requiredType {
		panic(fmt.Sprint("This is not a secret. Got: ", myStuff.Kind))
	}

	for k, v := range myStuff.Data {
		decoded, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			fmt.Println("Error:", err)
		}
		myStuff.Data[k] = string(decoded)
	}

	for k, v := range myStuff.Data {
		fmt.Printf("%s: \"%s\"\n", k, v)
	}
}

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "secret reverses a Kubernetes YAML secret and output plaintext data.",
	Long:  "secret reverses a Kubernetes YAML secret and output plaintext data.",
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
		transformSecret(text)
	},
}

func init() {
	rootCmd.AddCommand(secretCmd)
}

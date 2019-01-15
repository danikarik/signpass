package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/danikarik/signpass/pkg/exec"
	"github.com/spf13/cobra"
)

func init() {
	keyCmd.Flags().StringVarP(&certFile, "in", "i", "", "Private p12 certificate path")
	keyCmd.Flags().StringVarP(&outputFile, "out", "o", "", "PEM output file path")
	keyCmd.Flags().StringVarP(&password, "pass", "p", "", "Private key password")
	keyCmd.Flags().BoolVarP(&native, "native", "n", false, "Use native std instead of OpenSSL")
	rootCmd.AddCommand(keyCmd)
}

var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Export the key pass in PEM format",
	RunE: func(cmd *cobra.Command, args []string) error {
		if certFile == "" || outputFile == "" || password == "" {
			return errors.New("missing required parameters")
		}
		if native {
			return runGoKeyCommand()
		}
		return runKeyCommand()
	},
}

func runKeyCommand() error {
	c, err := exec.New("openssl", "pkcs12", "-in", certFile, "-nocerts", "-out", outputFile, "-passin", "pass:"+password, "-passout", "pass:"+password)
	if err != nil {
		return err
	}
	defer c.Stderr.Close()
	defer c.Stdout.Close()
	err = c.Start()
	if err != nil {
		return err
	}
	cout := c.StderrStr()
	err = c.Wait()
	if err != nil {
		return err
	}
	if cout != macVerified {
		return errors.New("could verify MAC")
	}
	fs, err := os.Stat(outputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("could not generate file")
		}
		return err
	}
	fmt.Printf("%s is generated.\n", fs.Name())
	return nil
}

func runGoKeyCommand() error {
	return errors.New("not implemented yet")
}

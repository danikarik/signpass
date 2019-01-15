package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/danikarik/signpass/pkg/pass"

	"github.com/spf13/cobra"
)

var (
	certFile     string
	passFile     string
	outputFile   string
	password     string
	native       bool
	wwdrFile     string
	passCertFile string
	passKeyFile  string
	rawPackage   string
	passesDir    string
)

const (
	macVerified = "MAC verified OK"
)

func init() {
	rootCmd.Flags().StringVarP(&wwdrFile, "wwdr", "w", "", "Apple Worldwide Developer Relations Certification Authority")
	rootCmd.Flags().StringVarP(&passCertFile, "signer", "s", "", "Pass certificate path")
	rootCmd.Flags().StringVarP(&passKeyFile, "key", "k", "", "Pass key path")
	rootCmd.Flags().StringVarP(&rawPackage, "raw", "r", "", "Raw package path")
	rootCmd.Flags().StringVarP(&password, "pass", "p", "", "Private key password")
	rootCmd.Flags().StringVarP(&passesDir, "dir", "d", "", ".pkpass output directory")
}

var rootCmd = &cobra.Command{
	Use:   "signpass",
	Short: "Signpass is CLI for Apple PassKit Generation",
	RunE: func(cmd *cobra.Command, args []string) error {
		return createPackage()
	},
}

// Execute runs signature generation script.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func createPackage() error {
	// check for arguments
	if !checkargs() {
		return errors.New("missing required parameters")
	}
	// check if given raw package is exists
	err := pass.CheckDir(rawPackage)
	if err != nil {
		return err
	}
	// create package
	pack, err := pass.Pack()
	if err != nil {
		return nil
	}
	// copy template dir to temp dir
	err = pack.CopyFrom(rawPackage)
	if err != nil {
		return err
	}
	err = pack.SetSerialNumber()
	if err != nil {
		return err
	}
	err = pack.Manifest()
	if err != nil {
		return err
	}
	err = pack.Signature(wwdrFile, passCertFile, passKeyFile, password)
	if err != nil {
		return err
	}
	err = pack.Zip()
	if err != nil {
		return err
	}
	err = pack.PKPass(passesDir)
	if err != nil {
		return err
	}
	err = pack.Clean()
	if err != nil {
		return err
	}
	return nil
}

func checkargs() bool {
	if wwdrFile == "" ||
		passCertFile == "" ||
		passKeyFile == "" ||
		rawPackage == "" ||
		password == "" ||
		passesDir == "" {
		return false
	}
	return true
}

func getPassJSON(p string) (path string, found bool) {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.Name() == "pass.json" {
			found = true
			path = filepath.Join(p, file.Name())
			return
		}
	}
	return
}

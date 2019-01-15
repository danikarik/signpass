package pass

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Package used for pkpass generation.
type Package struct {
	tempDir      string
	serialNumber string
}

// Pack returns new instance of package.
func Pack() (*Package, error) {
	tmp, err := ioutil.TempDir("", "passkit")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Creating temp dir at %s\n", tmp)
	return &Package{tempDir: tmp}, nil
}

// Location returns package path.
func (p *Package) Location() string {
	return p.tempDir
}

// CopyFrom copies the given directory to temp dir.
func (p *Package) CopyFrom(dir string) error {
	err := copyDir(dir, p.tempDir)
	if err != nil {
		return err
	}
	fmt.Println("Copying pass to temp directory")
	return nil
}

// Clean deletes all unused file and dirs.
func (p *Package) Clean() error {
	return os.RemoveAll(p.tempDir)
}

// SetSerialNumber creates new UUID for pass.json.
func (p *Package) SetSerialNumber() error {
	uuid, err := newUUID()
	if err != nil {
		return err
	}
	err = updatePassJSON(p.tempDir, uuid)
	if err != nil {
		return err
	}
	p.serialNumber = uuid
	fmt.Printf("Updated serial number %s\n", uuid)
	return nil
}

// Manifest cleans .DS_Store and creates manifest.json.
func (p *Package) Manifest() error {
	err := removeDSStore(p.tempDir)
	if err != nil {
		return err
	}
	data, err := manifestDict(p.tempDir)
	if err != nil {
		return err
	}
	err = createManifest(p.tempDir, data)
	if err != nil {
		return err
	}
	return nil
}

// Signature signs `manifest.json`
// and creates signature in DER format.
func (p *Package) Signature(wwdr, cert, key, password string) error {
	err := createSignature(p.tempDir, wwdr, cert, key, password)
	if err != nil {
		return err
	}
	return nil
}

// Zip compresses directory to `.pkpass`.
func (p *Package) Zip() error {
	err := createZip(p.tempDir, p.serialNumber)
	if err != nil {
		return err
	}
	fmt.Println("Compressing the pass")
	return nil
}

// PKPass copies generated `.pkpass` to destination folder.
func (p *Package) PKPass(dst string) error {
	return copyPKPass(p.tempDir, dst, p.serialNumber)
}

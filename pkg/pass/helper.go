package pass

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	// DSStore refers to `.DS_Store`
	DSStore = ".DS_Store"
	// PassJSON refers to `pass.json`
	PassJSON = "pass.json"
	// ManifestJSON refers to `manifest.json`
	ManifestJSON = "manifest.json"
	// Signature refers to `signature`
	Signature = "signature"
)

// CheckDir checks does directory exist.
func CheckDir(dir string) error {
	return exists(dir)
}

func exists(filename string) error {
	_, err := os.Stat(filename)
	if err != nil || os.IsNotExist(err) {
		return err
	}
	return nil
}

func copyDir(src, dst string) error {
	fs, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, fs.Mode())
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		srcfp := filepath.Join(src, f.Name())
		dstfp := filepath.Join(dst, f.Name())

		if f.IsDir() {
			if err = copyDir(srcfp, dstfp); err != nil {
				return err
			}
		} else {
			if err = copyFile(srcfp, dstfp); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}

	fs, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, fs.Mode())
}

func removeDSStore(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			if err = removeDSStore(path); err != nil {
				return err
			}
		}
		if file.Name() == DSStore {
			if err = removeFile(path); err != nil {
				return err
			}
		}
	}
	fmt.Println("Cleaning .DS_Store files")
	return nil
}

func removeFile(src string) error {
	return os.Remove(src)
}

func newUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func updatePassJSON(dir, uuid string) error {
	path := filepath.Join(dir, PassJSON)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return err
	}
	data["serialNumber"] = uuid
	pj, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, pj, 0664)
	if err != nil {
		return err
	}
	return nil
}

func excludedFromManifest(name string) bool {
	return name == ManifestJSON || name == Signature
}

func manifestDict(dir string) ([]byte, error) {
	dict := map[string]string{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	// walk dir and calculate hash for each file inside
	for _, file := range files {
		// skip directories
		if file.IsDir() {
			continue
		}
		// skip excluded files
		if excludedFromManifest(file.Name()) {
			continue
		}
		// calculate hash
		path := filepath.Join(dir, file.Name())
		hash, err := fileHash(path)
		if err != nil {
			return nil, err
		}
		// write result to map
		key := file.Name()
		val := string(hash)
		dict[key] = val
	}
	// marshal map to JSON bytes
	buf, err := json.Marshal(dict)
	if err != nil {
		return nil, err
	}
	// exit
	return buf, nil
}

func fileHash(path string) ([]byte, error) {
	// read file's content
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// create new SHA1 hasher
	h := sha1.New()
	// hash file content
	_, err = h.Write(content)
	if err != nil {
		return nil, err
	}
	// sum hasher
	sum := h.Sum(nil)
	// create file slice
	fh := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(fh, sum)
	// exit
	return fh, nil
}

func createManifest(dir string, content []byte) error {
	// define full path
	path := filepath.Join(dir, ManifestJSON)
	// create `manifest.json` and write data
	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}
	fmt.Println("Generating JSON manifest")
	return nil
}

func createSignature(dir, wwdr, cert, key, password string) error {
	inFile := filepath.Join(dir, ManifestJSON)
	outFile := filepath.Join(dir, Signature)
	cmd := exec.Command(
		"openssl",
		"smime", "-binary",
		"-sign",
		"-certfile", wwdr,
		"-signer", cert,
		"-inkey", key,
		"-in", inFile,
		"-out", outFile,
		"-outform", "DER",
		"-passin", "pass:"+password,
	)
	// execute command
	err := cmd.Run()
	if err != nil {
		return err
	}
	// check is `signature` created
	err = exists(outFile)
	if err != nil {
		return err
	}
	// exit
	fmt.Println("Signing the manifest")
	return nil
}

func createZip(dir, uuid string) error {
	path := filepath.Join(dir, uuid+".pkpass")
	// create new zip file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// create new zip writer
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()
	// get file list
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	// add files to zip
	for _, file := range files {
		filename := filepath.Join(dir, file.Name())
		if filepath.Ext(filename) == ".pkpass" {
			continue
		}
		zipfile, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer zipfile.Close()
		fi, err := zipfile.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		header.Name = file.Name()
		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, zipfile)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyPKPass(srcdir, dstdir, uuid string) error {
	filename := uuid + ".pkpass"
	srcfp := filepath.Join(srcdir, filename)
	dstfp := filepath.Join(dstdir, filename)
	err := exists(srcfp)
	if err != nil {
		return err
	}
	err = copyFile(srcfp, dstfp)
	if err != nil {
		return err
	}
	return nil
}

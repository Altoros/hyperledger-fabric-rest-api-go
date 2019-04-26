package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
)

func PrependPathToTar(tarReared io.Reader, prependPath string) ([]byte, error) {

	gzr, err := gzip.NewReader(tarReared)
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(gzr)

	var codePackage bytes.Buffer
	gw := gzip.NewWriter(&codePackage)
	tw := tar.NewWriter(gw)

LOOP:
	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			break LOOP

		case err != nil:
			return nil, err

		case header == nil:
			continue
		}

		switch header.Typeflag {

		case tar.TypeReg:
			header.Name = prependPath + header.Name

			if err := tw.WriteHeader(header); err != nil {
				return nil, err
			}

			if _, err := io.Copy(tw, tr); err != nil {
				return nil, err
			}
		}
	}

	err = tw.Close()
	if err != nil {
		return nil, err
	}

	err = gw.Close()
	if err != nil {
		return nil, err
	}

	err = gzr.Close()
	if err != nil {
		return nil, err
	}

	return codePackage.Bytes(), nil
}


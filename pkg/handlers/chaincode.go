package handlers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"github.com/labstack/echo/v4"
	"io"
)

// TODO move to some kind of utils file
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


func GetChaincodesInstalledHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	peer, err := c.CurrentPeer()
	if err != nil {
		return err
	}

	jsonString, err := c.Fsc().InstalledChaincodes(peer)
	return GetJsonOutputWrapper(c, jsonString, err)
}

// Get instantiated chaincodes list
func GetChaincodesInstantiatedHandler(ec echo.Context) error {
	c := ec.(*ApiContext)
	jsonString, err := c.Fsc().InstantiatedChaincodes(c.Param("channelId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}

func GetChaincodesInfoHandler(ec echo.Context) error {
	c := ec.(*ApiContext)
	// TODO validate
	jsonString, err := c.Fsc().ChaincodeInfo(c.Param("channelId"), c.Param("chaincodeId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}

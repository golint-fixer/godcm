package dcmimage

import (
	"os"
	"strconv"
	"testing"

	"github.com/grayzone/godcm/core"
)

func gettestdatafolder() string {
	cur, err := os.Getwd()
	if err != nil {
		return ""
	}
	result := cur + "/../test/data/"
	return result
}

func convert2bmp(t *testing.T, bits uint16) {
	cases := []struct {
		in   string
		want bool
	}{
		{"MR-MONO2-8-16x-heart.dcm", false},
		{"US-MONO2-8-8x-execho.dcm", false},
		{"xr_tspine.dcm", false},
		{"xr_chest.dcm", false},
		{"IM0.dcm", false},
		{"image_09-12-2013_4.dcm", false},
		{"CT-MONO2-16-ankle", false},
		{"xr_chicken2.dcm", true},
		{"T23/IM-0001-0001.dcm", true},
		{"IM-0001-0010.dcm", true},
		{"GH195.dcm", true},
		{"GH064.dcm", true},
		{"GH177_D_CLUNIE_CT1_IVRLE_BigEndian_undefined_length.dcm", false},
		{"GH177_D_CLUNIE_CT1_IVRLE_BigEndian_ELE_undefinded_length.dcm", false},
	}
	for _, c := range cases {
		var reader core.DcmReader
		reader.IsReadPixel = true
		reader.IsReadValue = true
		filepath := gettestdatafolder() + c.in
		err := reader.ReadFile(filepath)

		isCompressed, err := reader.IsCompressed()
		if err != nil {
			t.Errorf("WriteBMP() %s", err.Error())
		}

		isBigEndian, err := reader.IsBigEndian()
		if err != nil {
			t.Errorf("WriteBMP() %s", err.Error())
		}

		pixeldata := reader.Dataset.PixelData()

		var img DcmImage

		img.IsCompressed = isCompressed
		if c.want != img.IsCompressed {
			t.Errorf("WriteBMP(%s), isCompressed want %v got %v", c.in, c.want, isCompressed)
		}

		img.IsBigEndian = isBigEndian

		var num interface{}

		num, _ = strconv.ParseUint(reader.Dataset.BitsAllocated(), 10, 16)
		img.BitsAllocated = uint16(num.(uint64))

		num, _ = strconv.ParseUint(reader.Dataset.Columns(), 10, 32)
		img.Columns = uint32(num.(uint64))

		num, _ = strconv.ParseUint(reader.Dataset.Rows(), 10, 32)
		img.Rows = uint32(num.(uint64))

		num, _ = strconv.ParseUint(reader.Dataset.HighBit(), 10, 16)
		img.HighBit = uint16(num.(uint64))

		num, _ = strconv.ParseFloat(reader.Dataset.WindowCenter(), 64)
		img.WindowCenter = num.(float64)

		num, _ = strconv.ParseFloat(reader.Dataset.WindowWidth(), 64)
		img.WindowWidth = num.(float64)

		num, _ = strconv.ParseFloat(reader.Dataset.RescaleIntercept(), 64)
		img.RescaleIntercept = num.(float64)

		num, _ = strconv.ParseFloat(reader.Dataset.RescaleSlope(), 64)
		img.RescaleSlope = num.(float64)

		num, _ = strconv.ParseUint(reader.Dataset.PixelRepresentation(), 10, 16)
		img.PixelRepresentation = uint16(num.(uint64))

		img.PhotometricInterpretation = reader.Dataset.PhotometricInterpretation()

		img.PixelData = pixeldata

		bmpfile := c.in + "_ " + strconv.Itoa(int(bits)) + ".bmp"
		err = img.WriteBMP(bmpfile, bits, 0)
		defer os.Remove(bmpfile)
		if err != nil {
			//		t.Errorf("WriteBMP() %s", err.Error())
		}
	}

}

func TestWrite8BMP(t *testing.T) {

	convert2bmp(t, 8)

}

func TestWrite24BMP(t *testing.T) {
	convert2bmp(t, 24)
}

func TestWrite32BMP(t *testing.T) {
	convert2bmp(t, 32)
}
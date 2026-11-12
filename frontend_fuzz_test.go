package zabbix

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

// PNG decoding crosses the untrusted frontend-response boundary. Dimensions
// must be checked before full allocation, even when HTTP reports image/png.
func FuzzValidateChartImage(f *testing.F) {
	var seed bytes.Buffer
	if err := png.Encode(&seed, image.NewNRGBA(image.Rect(0, 0, 64, 64))); err != nil {
		f.Fatal(err)
	}
	f.Add(seed.Bytes())
	f.Add([]byte("<html>login required</html>"))
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1<<20 {
			t.Skip("outside fuzz input budget")
		}
		if err := validateChartImage(data, 64, 64); err == nil {
			config, err := png.DecodeConfig(bytes.NewReader(data))
			if err != nil || config.Width < 64 || config.Height < 64 || int64(config.Width)*int64(config.Height) > 16<<20 {
				t.Fatalf("accepted unsafe PNG dimensions: %v, %v", config, err)
			}
		}
	})
}

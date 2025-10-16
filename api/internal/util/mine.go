package util

import "bytes"

func SniffMimeForOCR(b []byte) string {
	// JPEG: FF D8
	if len(b) >= 2 && b[0] == 0xFF && b[1] == 0xD8 {
		return "JPEG"
	}
	// PNG
	if len(b) >= 8 &&
		b[0] == 0x89 && b[1] == 0x50 && b[2] == 0x4E && b[3] == 0x47 &&
		b[4] == 0x0D && b[5] == 0x0A && b[6] == 0x1A && b[7] == 0x0A {
		return "PNG"
	}
	// PDF
	if len(b) >= 5 && b[0] == '%' && b[1] == 'P' && b[2] == 'D' && b[3] == 'F' && b[4] == '-' {
		return "PDF"
	}
	return ""
}

func SniffMimeHTTP(b []byte) string {
	if len(b) >= 2 && b[0] == 0xFF && b[1] == 0xD8 {
		return "image/jpeg"
	}
	if len(b) >= 8 &&
		b[0] == 0x89 && b[1] == 0x50 && b[2] == 0x4E && b[3] == 0x47 &&
		b[4] == 0x0D && b[5] == 0x0A && b[6] == 0x1A && b[7] == 0x0A {
		return "image/png"
	}
	return "application/octet-stream"
}

func MakeDataURL(mime, b64 string) string {
	return "data:" + mime + ";base64," + b64
}

// SniffHEICorAVIF пытается распознать контейнеры ISOBMFF (HEIC/HEIF/AVIF).
// Ищет сигнатуру ftyp и совместимые бренды в первых байтах.
func SniffHEICorAVIF(data []byte) string {
	if len(data) < 12 {
		return ""
	}
	// ISO BMFF: bytes 4..7 == 'ftyp'
	if !bytes.Equal(data[4:8], []byte{'f', 't', 'y', 'p'}) {
		return ""
	}
	// Major brand
	major := string(data[8:12])
	if isHeicBrand(major) {
		return "image/heic"
	}
	if isAvifBrand(major) {
		return "image/avif"
	}
	// Просканируем ещё немного совместимые бренды (по 4 байта)
	// ограничимся первыми 64 байтами для простоты
	limit := len(data)
	if limit > 64 {
		limit = 64
	}
	for i := 16; i+4 <= limit; i += 4 { // пропускаем size(4) + 'ftyp'(4) + major(4) + minor(4)
		b := string(data[i : i+4])
		if isHeicBrand(b) {
			return "image/heic"
		}
		if isAvifBrand(b) {
			return "image/avif"
		}
	}
	return ""
}

func isHeicBrand(b string) bool {
	switch b {
	case "heic", "heix", "hevc", "hevx", "mif1", "msf1", "heis", "hevm":
		return true
	}
	return false
}

func isAvifBrand(b string) bool {
	switch b {
	case "avif", "avis":
		return true
	}
	return false
}

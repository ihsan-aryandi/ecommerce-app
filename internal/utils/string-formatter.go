package utils

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatRupiah(value float64) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("Rp%0.0f", value)
}

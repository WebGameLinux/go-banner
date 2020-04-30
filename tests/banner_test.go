package tests

import (
		"fmt"
		"github.com/WebGameLinux/go-banner/fonts"
		"testing"
)

func TestCreateBanner(t *testing.T) {
		var text = "cms"
		str := fonts.DrawText2ImageAscii(text)
		fmt.Println(str)
}

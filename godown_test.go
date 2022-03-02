package godown

import "testing"

func TestDownload(t *testing.T) {
	Download("https://www.google.com/robots.txt", "robots.txt")
}

package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type AnsiCode string

const (
	reset       AnsiCode = "\033[0m"
	Bold        AnsiCode = "\033[1m"
	Italic      AnsiCode = "\033[3m"
	Underline   AnsiCode = "\033[4m"
	Red         AnsiCode = "\033[31m"
	Yellow      AnsiCode = "\033[33m"
	Blue        AnsiCode = "\033[34m"
	ClearScreen AnsiCode = "\033[H\033[2J"
)

type AnsiString string

func Decorate(s string) *AnsiString {
	tmp := AnsiString(s)
	return &tmp
}

func (s *AnsiString) Format(codes ...AnsiCode) *AnsiString {
	var sb strings.Builder
	for _, code := range codes {
		sb.WriteString(string(code))
	}
	sb.WriteString(string(*s))
	if strings.HasSuffix(sb.String(), string(reset)) {
		return Decorate(sb.String())
	}
	sb.WriteString(string(reset))
	return Decorate(sb.String())
}

func (s *AnsiString) Bold() *AnsiString {
	return s.Format(Bold)
}

func (s *AnsiString) Underline() *AnsiString {
	return s.Format(Underline)
}

func (s *AnsiString) Italic() *AnsiString {
	return s.Format(Italic)
}

func Cls() {
	fmt.Println(ClearScreen)
}

func IsImageFile(f string) bool {
	f = strings.ToLower(filepath.Ext(f))
	return f == ".png" ||
		f == ".jpg" ||
		f == ".jpeg"
}

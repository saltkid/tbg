package main

import (
	"path/filepath"
	"strings"
)

type AnsiCode string

const (
	reset     AnsiCode = "\033[0m"
	Bold      AnsiCode = "\033[1m"
	Italic    AnsiCode = "\033[3m"
	Underline AnsiCode = "\033[4m"
	Blue      AnsiCode = "\033[34m"
	Red       AnsiCode = "\033[31m"
	Yellow    AnsiCode = "\033[33m"
)

type Styled string

func Decorate(s string) *Styled {
	tmp := Styled(s)
	return &tmp
}

func (s *Styled) String() string {
	return string(*s)
}

func (s *Styled) Format(codes ...AnsiCode) *Styled {
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

func (s *Styled) Bold() *Styled {
	return s.Format(Bold)
}

func (s *Styled) Underline() *Styled {
	return s.Format(Underline)
}

func (s *Styled) Italic() *Styled {
	return s.Format(Italic)
}

func IsImageFile(f string) bool {
	f = strings.ToLower(filepath.Ext(f))
	return f == ".png" ||
		f == ".jpg" ||
		f == ".jpeg"
}

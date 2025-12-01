package main

import (
	"regexp"
	"time"
)

type Parser interface {
	Parse(line string) (Log, error)
}

// Stategy 1: Regex paser
type RegexParser struct {
	Service string
	Regex   *regexp.Regexp 
}

func NewRegexParser(service string) *RegexParser {
	return &RegexParser{
		Service: service,
		Regex: regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\s+(?:\[(.*?)\]\s+)?\[(.*?)\]\s+(.*)$`),
	}
}

func (p *RegexParser) Parse(line string) (Log, error) {
	l := Log{Service: p.Service, Level: "info", Message: line}

	matches := p.Regex.FindStringSubmatch(line)
	if matches == nil {
		return l, nil 
	}

	if t, err := time.Parse("2006-01-02 15:04:05", matches[1]); err == nil {
		l.Timestamp = t 
	}

	if matches[2] != "" {
		l.Service = matches[2]
	}

	l.Level = matches[3]
	l.Message = matches[4]

	return l, nil 
}


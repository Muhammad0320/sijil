package main

import (
	"encoding/json"
	"regexp"
	"time"
)

type Parser interface {
	Parse(line string) (Log, error)
}

// Stategy 1: Regex parser
type RegexParser struct {
	Service string
	Regex   *regexp.Regexp
}

func NewRegexParser(service string) *RegexParser {
	return &RegexParser{
		Service: service,
		Regex:   regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\s+(?:\[(.*?)\]\s+)?\[(.*?)\]\s+(.*)$`),
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

// -- Strategy 2: JSON Parser ---
type JsonParser struct {
	Service string
}

func NewJsonParser(service string) *JsonParser {
	return &JsonParser{Service: service}
}

func (p *JsonParser) Parse(line string) (Log, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return Log{}, err
	}

	l := Log{
		Service: p.Service,
		Level:   "info",
		Data:    make(map[string]interface{}),
	}

	// Extract standard fields
	if val, ok := raw["service"].(string); ok {
		l.Service = val
		delete(raw, "service")
	}
	if val, ok := raw["level"].(string); ok {
		l.Level = val
		delete(raw, "level")
	}
	if val, ok := raw["message"].(string); ok {
		l.Message = val
		delete(raw, "message")
	}
	if val, ok := raw["timestamp"].(string); ok {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			l.Timestamp = t
		}
		delete(raw, "timestamp")
	}

	// Everything remaining goes into Data
	l.Data = raw

	return l, nil
}
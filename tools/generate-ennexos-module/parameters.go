package main

import (
	"cmp"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"unicode/utf8"

	"github.com/hansmi/modbus-exporter-modules/internal/ennexos"
)

var channelStringRe = regexp.MustCompile(`^(.+\.\w+)\[([0-9]+)\]$`)
var channelPhaseRe = regexp.MustCompile(`^(.+\.phs)([ABC])$`)

func cmpPointer[T cmp.Ordered](a, b *T) int {
	if a == nil {
		if b == nil {
			return 0
		}

		return +1
	}

	if b == nil {
		return -1
	}

	return cmp.Compare(*a, *b)
}

type parameterLabels struct {
	String *int
	Phase  int
}

func (l parameterLabels) toMap() map[string]string {
	result := map[string]string{}

	if s := l.String; s != nil {
		result["string"] = strconv.Itoa(*s)
	}

	if phase := l.Phase; phase != 0 {
		result["phase"] = strconv.Itoa(phase)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

type parameter struct {
	Name   string
	Labels parameterLabels
	Def    ennexos.Parameter
}

func (p parameter) compare(other parameter) int {
	if result := cmp.Compare(p.Name, other.Name); result != 0 {
		return result
	}

	if result := cmpPointer(p.Labels.String, other.Labels.String); result != 0 {
		return result
	}

	return cmp.Compare(p.Labels.Phase, other.Labels.Phase)
}

func preprocess(params []ennexos.Parameter) ([]parameter, error) {
	result := make([]parameter, 0, len(params))

	for _, d := range params {
		p := parameter{
			Name: d.Channel,
			Def:  d,
		}

		if groups := channelStringRe.FindStringSubmatch(d.Channel); len(groups) > 0 {
			if value, err := strconv.Atoi(groups[2]); err != nil {
				return nil, fmt.Errorf("string conversion: %w", err)
			} else {
				p.Name = groups[1]
				p.Labels.String = &value
			}
		} else if groups := channelPhaseRe.FindStringSubmatch(d.Channel); len(groups) > 0 {
			if r, _ := utf8.DecodeRuneInString(groups[2]); r == utf8.RuneError {
				return nil, fmt.Errorf("%w: invalid UTF-8 in %q", os.ErrInvalid, d.Channel)
			} else {
				p.Name = groups[1]
				p.Labels.Phase = 1 + int(r-'A')
			}
		}

		result = append(result, p)
	}

	slices.SortStableFunc(result, (parameter).compare)

	return result, nil
}

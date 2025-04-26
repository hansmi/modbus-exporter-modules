package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"regexp"
	"slices"

	"github.com/goccy/go-yaml"
)

func matchAny(expr []*regexp.Regexp, s string) bool {
	for _, i := range expr {
		if match := i.MatchString(s); match {
			return true
		}
	}

	return false
}

type configUnusedErr struct {
	AddressIncludeUnused []uint32 `json:",omitempty"`
	AddressExcludeUnused []uint32 `json:",omitempty"`

	ChannelIncludeUnused []*regexp.Regexp `json:",omitempty"`
	ChannelExcludeUnused []*regexp.Regexp `json:",omitempty"`
}

var _ error = (*configUnusedErr)(nil)

func (e configUnusedErr) isEmpty() bool {
	return (len(e.AddressIncludeUnused) == 0 &&
		len(e.AddressExcludeUnused) == 0 &&
		len(e.ChannelIncludeUnused) == 0 &&
		len(e.ChannelExcludeUnused) == 0)
}

func (e configUnusedErr) Error() string {
	type plain configUnusedErr

	details, err := json.Marshal((plain)(e))
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("configuration problem: %s", details)
}

type addressFilter struct {
	values []uint32
	unused map[uint32]struct{}
}

func newAddressFilter(values []uint32) *addressFilter {
	result := &addressFilter{
		values: slices.Clone(values),
		unused: make(map[uint32]struct{}, len(values)),
	}

	// Binary search requires sorted values.
	slices.Sort(result.values)

	for _, i := range result.values {
		result.unused[i] = struct{}{}
	}

	return result
}

func (f *addressFilter) getUnused() []uint32 {
	return slices.Sorted(maps.Keys(f.unused))
}

func (f *addressFilter) contains(value uint32) bool {
	_, found := slices.BinarySearch(f.values, value)

	if found {
		delete(f.unused, value)
	}

	return found
}

type regexpFilter struct {
	expr   []*regexp.Regexp
	unused map[*regexp.Regexp]struct{}
}

func newRegexpFilter(expr []*regexp.Regexp) *regexpFilter {
	result := &regexpFilter{
		expr:   slices.Clone(expr),
		unused: make(map[*regexp.Regexp]struct{}, len(expr)),
	}

	for _, i := range result.expr {
		result.unused[i] = struct{}{}
	}

	return result
}

func (f *regexpFilter) getUnused() []*regexp.Regexp {
	return slices.SortedFunc(maps.Keys(f.unused), func(a, b *regexp.Regexp) int {
		return cmp.Compare(a.String(), b.String())
	})
}

func (f *regexpFilter) contains(s string) bool {
	for _, i := range f.expr {
		if i.MatchString(s) {
			delete(f.unused, i)
			return true
		}
	}

	return false
}

type config struct {
	// Addresses to include in module.
	AddressInclude []uint32 `yaml:"address_include"`

	// Addresses to exclude.
	AddressExclude []uint32 `yaml:"address_exclude"`

	// Regular expression matching parameter channels to include in module
	// (after stripping string and phase).
	ChannelIncludeRegexp []*regexp.Regexp `yaml:"channel_include_regexp"`

	// Regular expressions matching parameter channels to exclude.
	ChannelExcludeRegexp []*regexp.Regexp `yaml:"channel_exclude_regexp"`

	// Regular expressions matching parameter channels of "counter" type.
	ChannelCounterTypeRegexp []*regexp.Regexp `yaml:"channel_counter_type_regexp"`
}

func loadConfig(path string) (*config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result config

	if err := yaml.UnmarshalWithOptions(content, &result, yaml.Strict()); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *config) filter(params []parameter) ([]parameter, error) {
	addressInclude := newAddressFilter(c.AddressInclude)
	addressExclude := newAddressFilter(c.AddressExclude)

	channelInclude := newRegexpFilter(c.ChannelIncludeRegexp)
	channelExclude := newRegexpFilter(c.ChannelExcludeRegexp)

	var result []parameter

	for _, p := range params {
		exclude := addressExclude.contains(p.Def.Address)
		exclude = channelExclude.contains(p.Name) || exclude

		include := addressInclude.contains(p.Def.Address)
		include = channelInclude.contains(p.Name) || include

		if include && !exclude {
			result = append(result, p)
		}
	}

	err := &configUnusedErr{
		AddressIncludeUnused: addressInclude.getUnused(),
		AddressExcludeUnused: addressExclude.getUnused(),

		ChannelIncludeUnused: channelInclude.getUnused(),
		ChannelExcludeUnused: channelExclude.getUnused(),
	}

	if !err.isEmpty() {
		return nil, err
	}

	return result, nil
}

func (c *config) metricType(p parameter) string {
	if matchAny(c.ChannelCounterTypeRegexp, p.Name) {
		return "counter"
	}

	return "gauge"
}

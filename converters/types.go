// Package converters defines a type for file conversion functions.
package converters

type Converter func(file string, tpl string, richtext bool) error

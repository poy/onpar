package diff

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

// ShowOpts contains all the options for Show.  They are
// unexported fields and should be controlled using Opt
// functions.
type ShowOpts struct {
	wrappers []func(string) string

	a *ShowOpts
	b *ShowOpts
}

func (o ShowOpts) format(v string) string {
	for _, w := range o.wrappers {
		v = w(v)
	}
	return v
}

type Opt func(ShowOpts) ShowOpts

// WithFormat returns an Opt that wraps up differences
// using a format string.  The format should contain
// one '%s' to add the difference string in.
func WithFormat(format string) Opt {
	return func(o ShowOpts) ShowOpts {
		o.wrappers = append(o.wrappers, func(v string) string {
			return fmt.Sprintf(format, v)
		})
		return o
	}
}

// Style represents display styles (like bold or italic)
// that we can display text as.
type Style int

const (
	Bold Style = 1 + iota
	Faint
	Italic
	Underline
	CrossedOut
)

// WithStyle returns an Opt that wraps up differences
// in a style.
func WithStyle(s Style) Opt {
	switch s {
	case CrossedOut:
		// We aren't matched 1:1 with fatih/color on this attribute.
		return withFatihColor(color.New(color.CrossedOut))
	default:
		return withFatihColor(color.New(color.Attribute(s)))
	}
}

// Color represents colors that we can display text as.
type Color int

const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// withFatihColor is a little helper to standardize Opt types
// that need to wrap up differences using colors from
// github.com/fatih/color.
func withFatihColor(c *color.Color) Opt {
	return func(o ShowOpts) ShowOpts {
		o.wrappers = append(o.wrappers, func(v string) string {
			return c.Sprint(v)
		})
		return o
	}
}

// WithFGColor returns an Opt that wraps up differences
// using a foreground color.
func WithFGColor(c Color) Opt {
	return withFatihColor(color.New(color.Attribute(c + 30)))
}

// WithBGColor returns an Opt that wraps up differences
// using a background color.
func WithBGColor(c Color) Opt {
	return withFatihColor(color.New(color.Attribute(c + 40)))
}

func applyOpts(o *ShowOpts, opts ...Opt) {
	for _, opt := range opts {
		*o = opt(*o)
	}
}

// Actual returns an Opt that only applies other Opt values
// to the actual value.
func Actual(opts ...Opt) Opt {
	return func(o ShowOpts) ShowOpts {
		if o.a == nil {
			o.a = &ShowOpts{}
		}
		applyOpts(o.a, opts...)
		return o
	}
}

// Expected returns an Opt that only applies other Opt values
// to the expected value.
func Expected(opts ...Opt) Opt {
	return func(o ShowOpts) ShowOpts {
		if o.b == nil {
			o.b = &ShowOpts{}
		}
		applyOpts(o.b, opts...)
		return o
	}
}

// Show takes two values and returns a string showing a
// diff of them.  The way that the differing values are
// shown can be controlled by opts.
//
// By default, we use: WithFormat(">%s<"), A(WithFormat("%s!=")).
func Show(actual, expected interface{}, opts ...Opt) string {
	o := ShowOpts{}
	if len(opts) == 0 {
		opts = append(opts, WithFormat(">%s<"), Actual(WithFormat("%s!=")))
	}
	for _, opt := range opts {
		o = opt(o)
	}
	return show(o, reflect.ValueOf(actual), reflect.ValueOf(expected))
}

func show(o ShowOpts, av, bv reflect.Value) string {
	showDiff := func(f string, a, b interface{}) string {
		afmt := fmt.Sprintf(f, a)
		if o.a != nil {
			afmt = o.a.format(afmt)
		}
		bfmt := fmt.Sprintf(f, b)
		if o.b != nil {
			bfmt = o.b.format(bfmt)
		}
		return o.format(afmt + bfmt)
	}

	if !av.IsValid() {
		if !bv.IsValid() {
			return "<nil>"
		}
		return showDiff("%v", "<nil>", bv.Interface())
	}
	if !bv.IsValid() {
		return showDiff("%v", av.Interface(), "<nil>")
	}

	if av.Kind() != bv.Kind() {
		return showDiff("%T", av.Interface(), bv.Interface())
	}

	switch av.Interface().(type) {
	case []rune, []byte, string:
		// we want to find differences in the middle of strings and
		// string-like types, whenever possible.
		if av.Len() != bv.Len() {
			break // let the default logic handle this
		}

		strTyp := reflect.TypeOf("")
		var curra, currb, out string
		for i := 0; i < av.Len(); i++ {
			match := av.Index(i).Interface() == bv.Index(i).Interface()
			if !match {
				curra += av.Index(i).Convert(strTyp).Interface().(string)
				currb += bv.Index(i).Convert(strTyp).Interface().(string)
				continue
			}
			if len(curra) > 0 {
				out += showDiff("%s", curra, currb)
				curra, currb = "", ""
			}
			out += av.Index(i).Convert(strTyp).Interface().(string)
		}
		if len(curra) > 0 {
			out += showDiff("%s", curra, currb)
			curra, currb = "", ""
		}
		return out
	}

	switch av.Kind() {
	case reflect.Ptr, reflect.Interface:
		return show(o, av.Elem(), bv.Elem())
	case reflect.Slice, reflect.Array:
		if av.Len() != bv.Len() {
			// TODO: do a more thorough diff of values
			return showDiff(fmt.Sprintf("%T(len %%d)", av.Interface()), av.Len(), bv.Len())
		}
		var elems []string
		for i := 0; i < av.Len(); i++ {
			elems = append(elems, show(o, av.Index(i), bv.Index(i)))
		}
		return "[ " + strings.Join(elems, ", ") + " ]"
	case reflect.Map:
		var parts []string
		for _, kv := range bv.MapKeys() {
			k := kv.Interface()
			bmv := bv.MapIndex(kv)
			amv := av.MapIndex(kv)
			if !amv.IsValid() {
				parts = append(parts, showDiff("%s", fmt.Sprintf("missing key %v", k), fmt.Sprintf("%v: %v", k, bmv.Interface())))
				continue
			}
			parts = append(parts, fmt.Sprintf("%v: %s", k, show(o, amv, bmv)))
		}
		for _, kv := range av.MapKeys() {
			// We've already compared all keys that exist in both maps; now we're
			// just looking for keys that only exist in a.
			if !bv.MapIndex(kv).IsValid() {
				k := kv.Interface()
				parts = append(parts, showDiff("%s", fmt.Sprintf("extra key %v: %v", k, av.MapIndex(kv).Interface()), fmt.Sprintf("%v: nil", k)))
				continue
			}
		}
		return "{" + strings.Join(parts, ", ") + "}"
	case reflect.Struct:
		if av.Type().Name() != bv.Type().Name() {
			return showDiff("%s", av.Type().Name(), bv.Type().Name()) + "(mismatched types)"
		}
		var parts []string
		for i := 0; i < bv.NumField(); i++ {
			f := bv.Type().Field(i)
			if f.PkgPath != "" {
				// unexported
				continue
			}
			name := f.Name
			bfv := bv.Field(i)
			afv := av.Field(i)
			parts = append(parts, fmt.Sprintf("%s: %s", name, show(o, afv, bfv)))
		}
		return fmt.Sprintf("%T{", av.Interface()) + strings.Join(parts, ", ") + "}"
	default:
		if av.Type().Comparable() {
			a, b := av.Interface(), bv.Interface()
			if a != b {
				return showDiff("%#v", a, b)
			}
			return fmt.Sprintf("%#v", a)
		}
		return o.format(fmt.Sprintf("UNSUPPORTED: could not compare values of type %T", av.Interface()))
	}
}

package modules

import (
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
)

// SetupPath sets up the path module
func SetupPath(vm *goja.Runtime) error {
	path := vm.NewObject()

	// path.join
	path.Set("join", func(call goja.FunctionCall) goja.Value {
		parts := make([]string, len(call.Arguments))
		for i, arg := range call.Arguments {
			parts[i] = arg.String()
		}
		result := filepath.Join(parts...)
		return vm.ToValue(result)
	})

	// path.resolve
	path.Set("resolve", func(call goja.FunctionCall) goja.Value {
		parts := make([]string, len(call.Arguments))
		for i, arg := range call.Arguments {
			parts[i] = arg.String()
		}

		var result string
		if len(parts) > 0 {
			result, _ = filepath.Abs(filepath.Join(parts...))
		} else {
			result, _ = filepath.Abs(".")
		}
		return vm.ToValue(result)
	})

	// path.basename
	path.Set("basename", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue("")
		}

		p := call.Arguments[0].String()
		ext := ""
		if len(call.Arguments) > 1 {
			ext = call.Arguments[1].String()
		}

		base := filepath.Base(p)
		if ext != "" && strings.HasSuffix(base, ext) {
			base = base[:len(base)-len(ext)]
		}

		return vm.ToValue(base)
	})

	// path.dirname
	path.Set("dirname", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(".")
		}

		p := call.Arguments[0].String()
		dir := filepath.Dir(p)
		return vm.ToValue(dir)
	})

	// path.extname
	path.Set("extname", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue("")
		}

		p := call.Arguments[0].String()
		ext := filepath.Ext(p)
		return vm.ToValue(ext)
	})

	// path.parse
	path.Set("parse", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.NewObject()
		}

		p := call.Arguments[0].String()
		result := vm.NewObject()

		dir := filepath.Dir(p)
		base := filepath.Base(p)
		ext := filepath.Ext(p)
		name := base
		if ext != "" {
			name = base[:len(base)-len(ext)]
		}

		result.Set("root", "")
		result.Set("dir", dir)
		result.Set("base", base)
		result.Set("ext", ext)
		result.Set("name", name)

		return result
	})

	// path.format
	path.Set("format", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue("")
		}

		obj := call.Arguments[0].ToObject(vm)
		if obj == nil {
			return vm.ToValue("")
		}

		dir := ""
		base := ""

		if dirVal := obj.Get("dir"); dirVal != nil && !goja.IsUndefined(dirVal) {
			dir = dirVal.String()
		}

		if baseVal := obj.Get("base"); baseVal != nil && !goja.IsUndefined(baseVal) {
			base = baseVal.String()
		} else {
			name := ""
			ext := ""
			if nameVal := obj.Get("name"); nameVal != nil && !goja.IsUndefined(nameVal) {
				name = nameVal.String()
			}
			if extVal := obj.Get("ext"); extVal != nil && !goja.IsUndefined(extVal) {
				ext = extVal.String()
			}
			base = name + ext
		}

		result := filepath.Join(dir, base)
		return vm.ToValue(result)
	})

	// path.isAbsolute
	path.Set("isAbsolute", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(false)
		}

		p := call.Arguments[0].String()
		return vm.ToValue(filepath.IsAbs(p))
	})

	// path.normalize
	path.Set("normalize", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(".")
		}

		p := call.Arguments[0].String()
		normalized := filepath.Clean(p)
		return vm.ToValue(normalized)
	})

	// path.relative
	path.Set("relative", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return vm.ToValue("")
		}

		from := call.Arguments[0].String()
		to := call.Arguments[1].String()
		rel, err := filepath.Rel(from, to)
		if err != nil {
			return vm.ToValue("")
		}
		return vm.ToValue(rel)
	})

	// path.sep
	path.Set("sep", string(filepath.Separator))

	// path.delimiter
	path.Set("delimiter", string(filepath.ListSeparator))

	// Register path module
	return RegisterModule(vm, "path", path)
}

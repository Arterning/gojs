package modules

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dop251/goja"
)

// SetupFS sets up the fs module
func SetupFS(vm *goja.Runtime) error {
	fs := vm.NewObject()

	// fs.readFileSync
	fs.Set("readFileSync", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("readFileSync requires a path argument"))
		}

		path := call.Arguments[0].String()
		encoding := "utf8"
		if len(call.Arguments) > 1 {
			if call.Arguments[1].ExportType().Kind().String() == "string" {
				encoding = call.Arguments[1].String()
			} else if obj := call.Arguments[1].ToObject(vm); obj != nil {
				if enc := obj.Get("encoding"); enc != nil && !goja.IsUndefined(enc) {
					encoding = enc.String()
				}
			}
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			panic(vm.ToValue("Error reading file: " + err.Error()))
		}

		if encoding == "utf8" || encoding == "utf-8" {
			return vm.ToValue(string(data))
		}
		// Return buffer for binary data
		return vm.ToValue(data)
	})

	// fs.writeFileSync
	fs.Set("writeFileSync", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(vm.ToValue("writeFileSync requires path and data arguments"))
		}

		path := call.Arguments[0].String()
		data := call.Arguments[1].String()

		err := ioutil.WriteFile(path, []byte(data), 0644)
		if err != nil {
			panic(vm.ToValue("Error writing file: " + err.Error()))
		}

		return goja.Undefined()
	})

	// fs.existsSync
	fs.Set("existsSync", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(false)
		}

		path := call.Arguments[0].String()
		_, err := os.Stat(path)
		return vm.ToValue(err == nil)
	})

	// fs.mkdirSync
	fs.Set("mkdirSync", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("mkdirSync requires a path argument"))
		}

		path := call.Arguments[0].String()
		recursive := false

		if len(call.Arguments) > 1 {
			if obj := call.Arguments[1].ToObject(vm); obj != nil {
				if rec := obj.Get("recursive"); rec != nil && !goja.IsUndefined(rec) {
					recursive = rec.ToBoolean()
				}
			}
		}

		var err error
		if recursive {
			err = os.MkdirAll(path, 0755)
		} else {
			err = os.Mkdir(path, 0755)
		}

		if err != nil {
			panic(vm.ToValue("Error creating directory: " + err.Error()))
		}

		return goja.Undefined()
	})

	// fs.readdirSync
	fs.Set("readdirSync", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("readdirSync requires a path argument"))
		}

		path := call.Arguments[0].String()
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(vm.ToValue("Error reading directory: " + err.Error()))
		}

		names := make([]string, len(files))
		for i, file := range files {
			names[i] = file.Name()
		}

		return vm.ToValue(names)
	})

	// fs.unlinkSync (delete file)
	fs.Set("unlinkSync", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("unlinkSync requires a path argument"))
		}

		path := call.Arguments[0].String()
		err := os.Remove(path)
		if err != nil {
			panic(vm.ToValue("Error removing file: " + err.Error()))
		}

		return goja.Undefined()
	})

	// fs.statSync
	fs.Set("statSync", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("statSync requires a path argument"))
		}

		path := call.Arguments[0].String()
		info, err := os.Stat(path)
		if err != nil {
			panic(vm.ToValue("Error getting file stats: " + err.Error()))
		}

		stat := vm.NewObject()
		stat.Set("isFile", func(call goja.FunctionCall) goja.Value {
			return vm.ToValue(!info.IsDir())
		})
		stat.Set("isDirectory", func(call goja.FunctionCall) goja.Value {
			return vm.ToValue(info.IsDir())
		})
		stat.Set("size", info.Size())
		stat.Set("mode", int(info.Mode()))
		stat.Set("mtime", info.ModTime().Unix())

		return stat
	})

	// Register fs module
	return RegisterModule(vm, "fs", fs)
}

// RegisterModule registers a module in the require system
func RegisterModule(vm *goja.Runtime, name string, module *goja.Object) error {
	// Get the module cache
	cacheVal := vm.Get("__moduleCache")
	if cacheVal == nil || goja.IsUndefined(cacheVal) {
		// Cache not initialized yet
		return fmt.Errorf("module cache not initialized - call SetupRequire first")
	}

	cache := cacheVal.ToObject(vm)
	if cache == nil {
		return fmt.Errorf("invalid module cache")
	}

	// Create module object with exports
	moduleObj := vm.NewObject()
	moduleObj.Set("exports", module)
	cache.Set(name, moduleObj)

	return nil
}

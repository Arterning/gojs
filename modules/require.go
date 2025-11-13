package modules

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/dop251/goja"
)

// SetupRequire sets up the require function for module loading
func SetupRequire(vm *goja.Runtime, currentDir string) error {
	// Create module cache
	cache := vm.NewObject()

	// Create require function
	require := func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("require() requires a module name"))
		}

		moduleName := call.Arguments[0].String()

		// Check cache first
		cached := cache.Get(moduleName)
		if cached != nil && !goja.IsUndefined(cached) {
			cachedModule := cached.ToObject(vm)
			if cachedModule != nil {
				exports := cachedModule.Get("exports")
				if exports != nil {
					return exports
				}
			}
		}

		// Check if it's a built-in module
		if moduleName == "fs" || moduleName == "path" {
			builtinModule := cache.Get(moduleName)
			if builtinModule != nil && !goja.IsUndefined(builtinModule) {
				moduleObj := builtinModule.ToObject(vm)
				if moduleObj != nil {
					return moduleObj.Get("exports")
				}
			}
			panic(vm.ToValue(fmt.Sprintf("Built-in module '%s' not found", moduleName)))
		}

		// Try to load as file
		var filePath string
		if filepath.IsAbs(moduleName) {
			filePath = moduleName
		} else if moduleName[0] == '.' {
			filePath = filepath.Join(currentDir, moduleName)
		} else {
			// Try node_modules
			filePath = filepath.Join(currentDir, "node_modules", moduleName)
		}

		// Try adding .js extension if not present
		if filepath.Ext(filePath) == "" {
			filePath += ".js"
		}

		// Read the file
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(vm.ToValue(fmt.Sprintf("Cannot find module '%s': %v", moduleName, err)))
		}

		// Create module object
		moduleObj := vm.NewObject()
		exportsObj := vm.NewObject()
		moduleObj.Set("exports", exportsObj)

		// Set module in cache before execution to handle circular dependencies
		cache.Set(moduleName, moduleObj)

		// Create module scope
		moduleDir := filepath.Dir(filePath)
		moduleRequire := func(call goja.FunctionCall) goja.Value {
			return SetupRequireInner(vm, moduleDir)(call)
		}

		// Wrap module code in a function
		wrappedCode := fmt.Sprintf(`(function(exports, require, module, __filename, __dirname) {
%s
})`, string(content))

		// Compile and run
		prg, err := goja.Compile(filePath, wrappedCode, false)
		if err != nil {
			cache.Delete(moduleName)
			panic(vm.ToValue(fmt.Sprintf("Error compiling module '%s': %v", moduleName, err)))
		}

		val, err := vm.RunProgram(prg)
		if err != nil {
			cache.Delete(moduleName)
			panic(vm.ToValue(fmt.Sprintf("Error loading module '%s': %v", moduleName, err)))
		}

		fn, ok := goja.AssertFunction(val)
		if !ok {
			cache.Delete(moduleName)
			panic(vm.ToValue(fmt.Sprintf("Error loading module '%s': not a function", moduleName)))
		}

		// Call the wrapped function
		_, err = fn(goja.Undefined(),
			exportsObj,
			vm.ToValue(moduleRequire),
			moduleObj,
			vm.ToValue(filePath),
			vm.ToValue(moduleDir),
		)

		if err != nil {
			cache.Delete(moduleName)
			panic(vm.ToValue(fmt.Sprintf("Error executing module '%s': %v", moduleName, err)))
		}

		// Return exports
		return moduleObj.Get("exports")
	}

	requireObj := vm.NewObject()
	requireObj.Set("cache", cache)

	vm.Set("require", require)
	vm.GlobalObject().Set("require", require)

	return nil
}

// SetupRequireInner creates a require function for a specific directory context
func SetupRequireInner(vm *goja.Runtime, currentDir string) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		// This is a simplified version - in a real implementation,
		// you'd want to share the cache across all require functions
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("require() requires a module name"))
		}

		moduleName := call.Arguments[0].String()

		// Check for built-in modules
		if moduleName == "fs" || moduleName == "path" {
			require := vm.Get("require")
			if fn, ok := goja.AssertFunction(require); ok {
				val, err := fn(goja.Undefined(), call.Arguments...)
				if err != nil {
					panic(err)
				}
				return val
			}
		}

		// For relative paths, resolve from current directory
		var filePath string
		if moduleName[0] == '.' {
			filePath = filepath.Join(currentDir, moduleName)
		} else {
			filePath = moduleName
		}

		// Delegate to main require
		require := vm.Get("require")
		if fn, ok := goja.AssertFunction(require); ok {
			val, err := fn(goja.Undefined(), vm.ToValue(filePath))
			if err != nil {
				panic(err)
			}
			return val
		}

		panic(vm.ToValue("require function not available"))
	}
}

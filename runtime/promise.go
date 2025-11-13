package runtime

import "github.com/dop251/goja"

// SetupPromise sets up Promise support in the runtime
func SetupPromise(vm *goja.Runtime, loop *EventLoop) error {
	// Promise constructor
	promiseCode := `
(function() {
	const PromiseState = {
		PENDING: 0,
		FULFILLED: 1,
		REJECTED: 2
	};

	function Promise(executor) {
		if (typeof executor !== 'function') {
			throw new TypeError('Promise executor must be a function');
		}

		this._state = PromiseState.PENDING;
		this._value = undefined;
		this._handlers = [];

		try {
			executor(
				(value) => this._resolve(value),
				(reason) => this._reject(reason)
			);
		} catch (error) {
			this._reject(error);
		}
	}

	Promise.prototype._resolve = function(value) {
		if (this._state !== PromiseState.PENDING) return;

		if (value === this) {
			this._reject(new TypeError('Cannot resolve promise with itself'));
			return;
		}

		if (value && (typeof value === 'object' || typeof value === 'function')) {
			let then;
			try {
				then = value.then;
			} catch (error) {
				this._reject(error);
				return;
			}

			if (typeof then === 'function') {
				let called = false;
				try {
					then.call(
						value,
						(result) => {
							if (called) return;
							called = true;
							this._resolve(result);
						},
						(error) => {
							if (called) return;
							called = true;
							this._reject(error);
						}
					);
				} catch (error) {
					if (!called) {
						this._reject(error);
					}
				}
				return;
			}
		}

		this._state = PromiseState.FULFILLED;
		this._value = value;
		this._executeHandlers();
	};

	Promise.prototype._reject = function(reason) {
		if (this._state !== PromiseState.PENDING) return;

		this._state = PromiseState.REJECTED;
		this._value = reason;
		this._executeHandlers();
	};

	Promise.prototype._executeHandlers = function() {
		if (this._state === PromiseState.PENDING) return;

		this._handlers.forEach((handler) => {
			queueMicrotask(() => {
				if (this._state === PromiseState.FULFILLED) {
					if (typeof handler.onFulfilled === 'function') {
						try {
							const result = handler.onFulfilled(this._value);
							handler.resolve(result);
						} catch (error) {
							handler.reject(error);
						}
					} else {
						handler.resolve(this._value);
					}
				} else if (this._state === PromiseState.REJECTED) {
					if (typeof handler.onRejected === 'function') {
						try {
							const result = handler.onRejected(this._value);
							handler.resolve(result);
						} catch (error) {
							handler.reject(error);
						}
					} else {
						handler.reject(this._value);
					}
				}
			});
		});

		this._handlers = [];
	};

	Promise.prototype.then = function(onFulfilled, onRejected) {
		return new Promise((resolve, reject) => {
			this._handlers.push({
				onFulfilled: onFulfilled,
				onRejected: onRejected,
				resolve: resolve,
				reject: reject
			});

			this._executeHandlers();
		});
	};

	Promise.prototype.catch = function(onRejected) {
		return this.then(null, onRejected);
	};

	Promise.prototype.finally = function(onFinally) {
		return this.then(
			(value) => {
				return Promise.resolve(onFinally()).then(() => value);
			},
			(reason) => {
				return Promise.resolve(onFinally()).then(() => {
					throw reason;
				});
			}
		);
	};

	Promise.resolve = function(value) {
		if (value instanceof Promise) {
			return value;
		}
		return new Promise((resolve) => resolve(value));
	};

	Promise.reject = function(reason) {
		return new Promise((resolve, reject) => reject(reason));
	};

	Promise.all = function(promises) {
		return new Promise((resolve, reject) => {
			if (!Array.isArray(promises)) {
				reject(new TypeError('Promise.all expects an array'));
				return;
			}

			if (promises.length === 0) {
				resolve([]);
				return;
			}

			let remaining = promises.length;
			const results = new Array(promises.length);

			promises.forEach((promise, index) => {
				Promise.resolve(promise).then(
					(value) => {
						results[index] = value;
						remaining--;
						if (remaining === 0) {
							resolve(results);
						}
					},
					(error) => {
						reject(error);
					}
				);
			});
		});
	};

	Promise.race = function(promises) {
		return new Promise((resolve, reject) => {
			if (!Array.isArray(promises)) {
				reject(new TypeError('Promise.race expects an array'));
				return;
			}

			promises.forEach((promise) => {
				Promise.resolve(promise).then(resolve, reject);
			});
		});
	};

	Promise.allSettled = function(promises) {
		return new Promise((resolve) => {
			if (!Array.isArray(promises)) {
				resolve([]);
				return;
			}

			if (promises.length === 0) {
				resolve([]);
				return;
			}

			let remaining = promises.length;
			const results = new Array(promises.length);

			promises.forEach((promise, index) => {
				Promise.resolve(promise).then(
					(value) => {
						results[index] = { status: 'fulfilled', value: value };
						remaining--;
						if (remaining === 0) {
							resolve(results);
						}
					},
					(reason) => {
						results[index] = { status: 'rejected', reason: reason };
						remaining--;
						if (remaining === 0) {
							resolve(results);
						}
					}
				);
			});
		});
	};

	return Promise;
})()
	`

	// Run the Promise implementation
	val, err := vm.RunString(promiseCode)
	if err != nil {
		return err
	}

	// Set Promise as a global
	vm.Set("Promise", val)

	return nil
}

// Test basic JavaScript features
console.log("=== GoJS Runtime Test ===");
console.log("");

// Test 1: Basic console.log
console.log("Test 1: Basic console.log");
console.log("Hello, GoJS!");
console.log("");

// Test 2: Variables and types
console.log("Test 2: Variables and types");
const name = "GoJS";
let version = 1.0;
var active = true;
console.log("Name:", name);
console.log("Version:", version);
console.log("Active:", active);
console.log("");

// Test 3: Arrays and objects
console.log("Test 3: Arrays and objects");
const arr = [1, 2, 3, 4, 5];
const obj = { a: 1, b: 2, c: 3 };
console.log("Array:", arr);
console.log("Object:", obj);
console.log("");

// Test 4: Functions
console.log("Test 4: Functions");
function add(a, b) {
    return a + b;
}
console.log("5 + 3 =", add(5, 3));

const multiply = (a, b) => a * b;
console.log("4 * 7 =", multiply(4, 7));
console.log("");

// Test 5: Promise
console.log("Test 5: Promise");
const promise = new Promise((resolve, reject) => {
    resolve("Promise resolved!");
});

promise.then((msg) => {
    console.log("Promise result:", msg);
});
console.log("");

// Test 6: setTimeout
console.log("Test 6: setTimeout");
console.log("Setting timeout for 100ms...");
setTimeout(() => {
    console.log("Timeout executed!");
}, 100);

// Test 7: Microtask (queueMicrotask)
console.log("");
console.log("Test 7: Microtask");
console.log("Before microtask");
queueMicrotask(() => {
    console.log("Microtask executed!");
});
console.log("After queueing microtask");

// Test 8: Async/await (if supported)
console.log("");
console.log("Test 8: Async function");
async function asyncTest() {
    console.log("Async function started");
    const result = await Promise.resolve("Async result!");
    console.log("Async function result:", result);
}
asyncTest();

console.log("");
console.log("=== All tests queued ===");

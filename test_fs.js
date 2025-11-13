// Test Node.js-style modules
console.log("=== Testing fs and path modules ===");
console.log("");

// Test fs module
console.log("Test 1: Writing a file");
const fs = require('fs');
fs.writeFileSync('hello.txt', 'Hello from GoJS!');
console.log("✓ File written: hello.txt");
console.log("");

console.log("Test 2: Reading a file");
const content = fs.readFileSync('hello.txt', 'utf8');
console.log("✓ File content:", content);
console.log("");

console.log("Test 3: Checking file existence");
const exists = fs.existsSync('hello.txt');
console.log("✓ File exists:", exists);
console.log("");

// Test path module
console.log("Test 4: Path operations");
const path = require('path');
const joined = path.join('foo', 'bar', 'baz.js');
console.log("✓ path.join:", joined);

const base = path.basename('/foo/bar/baz.js');
console.log("✓ path.basename:", base);

const dir = path.dirname('/foo/bar/baz.js');
console.log("✓ path.dirname:", dir);

const ext = path.extname('test.js');
console.log("✓ path.extname:", ext);
console.log("");

// Clean up
console.log("Test 5: Deleting the file");
fs.unlinkSync('hello.txt');
console.log("✓ File deleted");
console.log("");

console.log("=== All fs/path tests passed! ===");

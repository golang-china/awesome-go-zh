const fs = require('fs')
const buf = fs.readFileSync('./add.wasm');

async function run() {
	let res = await WebAssembly.instantiate(buf); // HL
	const pkg = res.instance.exports; // HL
	console.log(pkg.add(1, 2)); // HL
}

run();

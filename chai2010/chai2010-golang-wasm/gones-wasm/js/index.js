const go = new Go();
let mod, inst;

WebAssembly.instantiateStreaming(fetch("index.out.wasm"), go.importObject).then((result) => {
	mod = result.module;
	inst = result.instance;
	document.getElementById("runButton").disabled = false;
});

async function run() {
	//sayHi();
	await go.run(inst);
	inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
}

function run2() {
	go.run(inst);
}

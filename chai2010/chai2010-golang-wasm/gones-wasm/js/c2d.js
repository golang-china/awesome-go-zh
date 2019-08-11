function js_TV_SetPixels(p, w, h){
	var canv = document.getElementById('canvas_nes');
	canv.width = w;
	canv.height = h;
	var ctx = canv.getContext("2d");
	//ctx.fillStyle = '#00ffff';
	//ctx.fillRect(0, 0, w, h);
	var imgdata = new ImageData(p, w, h);
	ctx.putImageData(imgdata, 0, 0);
}

function sayHi() {
	console.log("hi22");
}

function jsSetPixels(canvas_id, pixel, width, height){
	//pixel = global.go.loadSlice(pixel);

	console.log("jsSetPixels:", typeof pixel)
	console.log("width:", width)
	console.log("height:", height)

	var canvas = document.getElementById(canvas_id);
	canvas.width = width;
	canvas.height = height;

	var ctx = canvas.getContext("2d");
	var imgdata = new ImageData(new Uint8ClampedArray(pixel), width, height);

	ctx.putImageData(imgdata, 0, 0);
}

function trySetPixels(canvas_id, pixel, width, height){
	console.log("trySetPixels")
	console.log("pixel[0]:", pixel[0])
	console.log("width:", width)
	console.log("height:", height)
}

var frame_count = 0;

function doFrameLoop() {
	var requestId = window.requestAnimationFrame( doFrameLoop );
	console.log(frame_count++);
	goSetupNes();
}

function genxxx(){
	var w = 640;
	var h = 480;
	var buf = new ArrayBuffer(w * h * 4);
	var p = new Uint8ClampedArray(buf);
	for (var y = 0; y < h; y++){
		for (var x = 0; x < w; x++){
			p[(y * w + x) * 4] = 255;
			p[(y * w + x) * 4 + 1] = 0;
			p[(y * w + x) * 4 + 2] = 0;
			p[(y * w + x) * 4 + 3] = 255;
		}
	}
	js_TV_SetPixels(p, w, h);
}

package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"syscall/js"
)

func main() {
	js.Global().Set("checkColorModel", js.FuncOf(checkColorModel))
	<-make(chan struct{})
}

func checkColorModel(this js.Value, args []js.Value) interface{} {
	// 第1引数(Uint8Array)を取得
	uint8Array := args[0]
	// JavaScriptでグローバルに定義されたコールバック関数を取得
	handleImageCheck := js.Global().Get("handleImageCheck")

	// Goのバイトスライスにコピー
	data := make([]byte, uint8Array.Get("length").Int())
	js.CopyBytesToGo(data, uint8Array)

	// バイトデータから画像をデコード
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		// デコードに失敗したら、JavaScriptのコンソールにエラーを出力
		js.Global().Get("console").Call("error", "Failed to decode image:", err.Error())
		// isDecoded: false, isNotGrayscale: false をJSに渡す
		handleImageCheck.Invoke(false, false)
		return nil
	}
	isNotGrayscale := true
	switch img.ColorModel() {
	case color.GrayModel, color.Gray16Model:
		isNotGrayscale = false
	default:
		isNotGrayscale = true
	}
	handleImageCheck.Invoke(true, isNotGrayscale)

	return nil
}

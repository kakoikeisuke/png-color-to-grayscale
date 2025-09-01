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
	js.Global().Set("convertImage", js.FuncOf(convertImage))
	<-make(chan struct{})
}

// グレースケールへの変換
func convertImage(_ js.Value, args []js.Value) interface{} {
	settings := args[0]
	showOutputImage := js.Global().Get("showOutputImage")

	// Uint8Arrayを[]byteにコピー
	uint8Array := settings.Get("fileData")
	data := make([]byte, uint8Array.Get("length").Int())
	js.CopyBytesToGo(data, uint8Array)

	// 設定を取得
	channel := settings.Get("channel").String()
	rWeight := settings.Get("rWeight").Float()
	gWeight := settings.Get("gWeight").Float()
	bWeight := settings.Get("bWeight").Float()
	bit := settings.Get("bit").String()
	invert := settings.Get("invert").Bool()

	// TODO: ここで受け取ったデータを使って画像変換処理を実装します
	js.Global().Get("console").Call("log", "Go received channel:", channel)
	js.Global().Get("console").Call("log", "Go received rWeight:", rWeight)
	js.Global().Get("console").Call("log", "Go received gWeight:", gWeight)
	js.Global().Get("console").Call("log", "Go received bWeight:", bWeight)
	js.Global().Get("console").Call("log", "Go received bit:", bit)
	js.Global().Get("console").Call("log", "Go received invert:", invert)

	// Goの[]byteをJavaScriptのUint8Arrayにコピー
	jsUint8Array := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(jsUint8Array, data)
	showOutputImage.Invoke(true, "", jsUint8Array)

	return nil
}

// カラーモデルの確認
func checkColorModel(_ js.Value, args []js.Value) interface{} {
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

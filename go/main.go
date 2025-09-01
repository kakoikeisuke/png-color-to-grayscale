package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
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

	// 画像をデコード
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		showOutputImage.Invoke(false, "画像のデコードに失敗しました。", js.Null())
		return nil
	}

	// 同じ解像度の画像を用意
	bounds := img.Bounds()
	var convertedImg image.Image
	var setPixel func(x, y int, grayValue float64)

	// bit数に応じて画像とセッターを初期化
	if bit == "bit-8" {
		gray8Img := image.NewGray(bounds)
		setPixel = func(x, y int, grayValue float64) {
			// 0-65535から0-255の範囲に変換
			gray8Value := grayValue / 257.0
			if invert {
				gray8Value = 255.0 - gray8Value
			}
			gray8Img.SetGray(x, y, color.Gray{Y: uint8(gray8Value)})
		}
		convertedImg = gray8Img
	} else {
		gray16Img := image.NewGray16(bounds)
		setPixel = func(x, y int, grayValue float64) {
			if invert {
				grayValue = 65535.0 - grayValue
			}
			gray16Img.SetGray16(x, y, color.Gray16{Y: uint16(grayValue)})
		}
		convertedImg = gray16Img
	}

	// RGB：重みの正規化
	if channel == "rgb-channel" {
		totalWeight := rWeight + gWeight + bWeight
		if totalWeight > 0 {
			rWeight /= totalWeight
			gWeight /= totalWeight
			bWeight /= totalWeight
		}
	}
	// ピクセルごとにグレースケール変換
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			var grayValue float64
			switch channel {
			case "rgb-channel":
				// RGBA()は0-65535の範囲の値を返すため、float64に変換して計算
				grayValue = float64(r)*rWeight + float64(g)*gWeight + float64(b)*bWeight
			case "r-channel":
				grayValue = float64(r)
			case "g-channel":
				grayValue = float64(g)
			case "b-channel":
				grayValue = float64(b)
			case "a-channel":
				grayValue = float64(a)
			}
			setPixel(x, y, grayValue)
		}
	}

	// 変換した画像をPNGにエンコード
	var buf bytes.Buffer
	if err := png.Encode(&buf, convertedImg); err != nil {
		showOutputImage.Invoke(false, "画像のエンコードに失敗しました。", js.Null())
		return nil
	}
	encodedData := buf.Bytes()

	// Goの[]byteをJavaScriptのUint8Arrayにコピー
	jsUint8Array := js.Global().Get("Uint8Array").New(len(encodedData))
	js.CopyBytesToJS(jsUint8Array, encodedData)
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

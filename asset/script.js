let file;
let fileName;
const inputView =document.getElementById('input-view');
const inputFile = document.getElementById('input-file');
inputView.addEventListener('click', () => {
    inputFile.click();
});
inputFile.addEventListener('change', (event) => {
    file = event.target.files[0];
    if (!file) {
        errorMessage("正常に画像が読み込まれませんでした。");
        return;
    }
    if (file.type !== 'image/png') {
        errorMessage("PNG画像を選択してください。");
        inputFile.value = '';
        return;
    }
    const reader = new FileReader();
    reader.onload = () => {
        const uint8Array = new Uint8Array(reader.result);
        window.checkColorModel(uint8Array);
    };
    reader.readAsArrayBuffer(file);
});

// Goから呼び出される
function handleImageCheck(isDecoded, isNotGrayscale) {
    // デコードできたか
    if (!isDecoded) {
        errorMessage("正常に画像が読み込まれませんでした。");
        inputFile.value = '';
        return;
    }
    // グレースケールではないか
    if (!isNotGrayscale) {
        errorMessage("選択された画像はすでにグレースケール画像です。");
        inputFile.value = '';
        return;
    }
    // 画像をimg要素に適用
    inputView.src = URL.createObjectURL(file);
    // ファイルネームを格納
    // ファイルネームが有効かどうかが実質的な変換準備の確認になる
    fileName = file.name;
    
    getConvertSetting();
}

function getConvertSetting() {
}

// エラーテキストを表示
function errorMessage(message) {
    closeErrorMessage();
    const backgroundFilter = document.createElement("div");
    backgroundFilter.className = "background-filter";
    document.body.appendChild(backgroundFilter);
    const errorMessageDiv = document.createElement("div");
    errorMessageDiv.className = "error-message-div";
    backgroundFilter.appendChild(errorMessageDiv);
    const errorMessageText = document.createElement("p");
    errorMessageText.textContent = message;
    errorMessageDiv.appendChild(errorMessageText);
    const closeButton = document.createElement("p");
    closeButton.className = "close-button";
    closeButton.textContent = '閉じる';
    closeButton.addEventListener('click', closeErrorMessage);
    errorMessageDiv.appendChild(closeButton);
    document.addEventListener('keydown', (event) => {
        if (event.key === 'Escape') {
            closeErrorMessage();
        }
    });
}
// エラーテキストを削除
function closeErrorMessage() {
    const backgroundFilter = document.querySelector(".background-filter");
    if (backgroundFilter) {
        backgroundFilter.remove();
    }
}

// WebAssemblyの用意
(async () => {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(fetch('asset/main.wasm'), go.importObject);
    await go.run(result.instance);
})();

// 設定が変更された際に呼び出し
const convertOption = document.getElementById('convert-option');
convertOption.addEventListener('change', (event) => {
    // 加重平均の数値入力の切り替え
    if (event.target.name === 'channel') {
        const isRgbSelected = document.getElementById('rgb-channel').checked;
        document.getElementById('r-weight').disabled = !isRgbSelected;
        document.getElementById('g-weight').disabled = !isRgbSelected;
        document.getElementById('b-weight').disabled = !isRgbSelected;
    }
    // 変換設定が変更されたらgetConvertSettingを呼び出す
    if (fileName) {
        getConvertSetting();
    }
});

// 加重平均のテンプレート
const templateAverage = document.getElementById('template-average');
templateAverage.addEventListener('click', () => {
    document.getElementById('r-weight').value = 1;
    document.getElementById('g-weight').value = 1;
    document.getElementById('b-weight').value = 1;
    if (fileName) {
        getConvertSetting();
    }
})
const templateHdtv = document.getElementById('template-hdtv');
templateHdtv.addEventListener('click', () => {
    document.getElementById('r-weight').value = 0.213;
    document.getElementById('g-weight').value = 0.715;
    document.getElementById('b-weight').value = 0.072;
    if (fileName) {
        getConvertSetting();
    }
})
const templateNtsc = document.getElementById('template-ntsc');
templateNtsc.addEventListener('click', () => {
    document.getElementById('r-weight').value = 0.299;
    document.getElementById('g-weight').value = 0.587;
    document.getElementById('b-weight').value = 0.114;
    if (fileName) {
        getConvertSetting();
    }
})
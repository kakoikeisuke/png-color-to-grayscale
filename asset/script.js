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
(async () => {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(fetch('asset/main.wasm'), go.importObject);
    await go.run(result.instance);
})();

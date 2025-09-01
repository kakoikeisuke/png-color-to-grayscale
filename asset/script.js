// 左の画像ビューをクリックすると, 画像を選択するinputに飛ぶ
const inputView =document.getElementById('input-view');
const inputFile = document.getElementById('input-file');
inputView.addEventListener('click', () => {
    inputFile.click();
});

// 画像が入力されると, 画像がpng画像か確認
inputFile.addEventListener('change', (event) => {
    const file = event.target.files[0];
    if (!file) {
        errorMessage("正常に画像が読み込まれませんでした。");
        return;
    }
    if (file.type !== 'image/png') {
        errorMessage("PNG画像を選択してください。");
        inputFile.value = '';
        return;
    }
    errorMessage("正常に読み込まれました。");
});

// エラーがあった際の表示
function errorMessage(message) {
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
function closeErrorMessage() {
    const backgroundFilter = document.querySelector(".background-filter");
    if (backgroundFilter) {
        backgroundFilter.remove();
    }
}
// 左の画像ビューをクリックすると, 画像を選択するinputに飛ぶ
const inputView =document.getElementById('input-view');
const inputFile = document.getElementById('input-file');
inputView.addEventListener('click', () => {
    inputFile.click();
});
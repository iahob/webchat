// 获取上载输入和按钮
var uploadInput = document.getElementById('uploadInput');
  // 监听上传按钮的点击事件

// 监听上载输入的改变事件
uploadInput.addEventListener('change', function(event) {
	var file = event.target.files[0];
	console.log(file)
});
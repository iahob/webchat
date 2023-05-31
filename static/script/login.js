document.getElementById("loginForm").addEventListener("submit",
	function(event) {
	event.preventDefault(); // 阻止表单的默认提交行为

	// 获取表单元素
	var form = event.target;
	var usernameInput = form.elements.username;
	var passwordInput = form.elements.password;

	// 获取用户输入的值
	var username = usernameInput.value;
	var password = passwordInput.value;

	// 发送登录请求到后端
	var data = {
		name: username,
		pwd: password
	};
	var url = "http://" + window.location.host + "/auth"
	// 使用fetch发送POST请求
	fetch(url, {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify(data)
	})
		.then(response => response.json())
		.then(responseData => {
			// 登录成功，处理后端返回的数据
			// 在这里可以进行页面跳转或其他操作
			if (responseData.code === 100) {
				localStorage.setItem("token",responseData.token)
				console.log(responseData);
				window.location.href = "http://"+ window.location.host + "/index.html"
			} else {
				console.log(responseData);
			}

		})
		.catch(error => {
			// 处理错误情况
			console.error("Login error:", error);
		});
});
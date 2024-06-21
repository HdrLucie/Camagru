document.getElementById("signUp").onclick = async function () {
	console.log("\n\nSIGN UP FUNCTION\n\n")
	var tmpUser = document.getElementById("Username")
	var tmpEmail = document.getElementById("Email")
	var tmpPassword = document.getElementById("Password")

	var login = tmpUser.value
	var email = tmpEmail.value
	var password = tmpPassword.value

	console.log("Login : " + login + " Email : " + email + " Mdp : " + password)
	if (login == "" || email == "" || password == "") {
		alert("Wrong username, email or password")
		// return
	}
	try {
		const response = await fetch("/signUp", {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({
				"email": email,
				"username": login,
				"password": password
			})
		});
		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}

		const data = await response.json();
		if (response.status === 201) {
			alert("User created successfully");
		} else {
			alert(`Error creating user: ${data.message}`);
		}
	} catch (error) {
		console.error("Error: ", error);
		alert(`Error creating user: ${error.message}`);
	}
}

document.getElementById("login").onclick = async function () {
	var user = document.getElementById("UsernameLogin")
	var mdp = document.getElementById("PasswordLogin")

	var login = user.value
	var password = mdp.value
	try {
		const response = await fetch("/login", {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({
				Username: login,
				Password: password
			})
		});
		const data = await response.json();

		const token = data.token;
		localStorage.setItem('token', token);

		console.log('Token:', token);
		if (response.status === 200) {
			alert("Successfully connected");
			window.location.href = data.redirectPath
		} else {
			alert(`Error connection: ${data.message}`);
		}
	} catch (error) {
		console.error("Error:", error);
		alert("Wrong username or password");
	}
}

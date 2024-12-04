console.log('La page entière est complètement chargée');

function checkPassword(password) {
	let regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,25}$/;
	let result = regex.test(password);
	return result
}

const resetButton = document.getElementById("forgetPassword");
if (resetButton) {
	console.log("Button found!");
	resetButton.onclick = function (event) {
		event.preventDefault();
		console.log("Reset password function");
		window.location.href = "/forgetPassword";
	};
} else {
	console.log("Button not found!");
}

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
	}
	var result = checkPassword(password);
	if (result == false) {
		alert("Password must contain between 8 and 25 characters, an uppercase letter, a lowercase letter, a number and a non-alphanumeric character (!@#$%&*...).")
		return
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
			console.log("\nAccount created successfully! Please check your email to verify your account.")
			window.location.href = data.redirectPath
		} else if (response.status === 409) {
			alert("Username or email already in use.")
		} else {
			alert(`Error creating user: ${data.message}`);
		}
	} catch (error) {
		console.error("Error: ", error);
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
				"Username": login,
				"Password": password
			})
		});
		const data = await response.json();
		const token = data.token;
		localStorage.setItem('token', token);
		console.log('Token:', token);
		if (response.status === 200) {
			window.location.href = data.redirectPath
		} else if (response.status === 401) {
			alert("Unable to log in. Please verify your credentials and try again.")
		} else if (response.status === 403) {
			alert("Account verification required. Please check your email to complete the verification process.")
		} else {
			alert(`Error connection: ${data.message}`);
		}
	} catch (error) {
		console.error("Error:", error);
		alert("Wrong username or password");
	}
}

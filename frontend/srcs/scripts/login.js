function checkPassword(password) {
	let regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,25}$/;
	let result = regex.test(password);
	return result
}

const resetButton = document.getElementById("forgetPassword");
if (resetButton) {
	resetButton.onclick = function (event) {
		event.preventDefault();
		window.location.href = "/forgetPassword";
	};
} 

document.getElementById("signUp").onclick = async function () {
	var login = document.getElementById("Username").value;
	var email = document.getElementById("Email").value;
	var password = document.getElementById("Password").value;

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
			throw new Error(`Response status: ${response.status}`);
		}
		const data = await response.json();
		if (data.success === false) {
			if (data.reason === "conflict") {
				alert("Username or email already in use.");
			} else {
				alert(data.message);
			}
		} else {
			window.location.href = data.redirectPath;
		}
	} catch (error) {
		alert("Username or email already in use.")
	}
}


// Login function
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
		if (!response.ok) {
			throw new Error(`Response status: ${response.status}`);
		}
		const data = await response.json();
		if (data.success === false) {
			if (data.reason === "Unauthorized") {
				alert("Unable to log in. Please verify your credentials and try again.");
				window.location.href("/connection");
			} else if (data.reason === "Forbidden") {
				alert("Account verification required. Please check your email to complete the verification process.");
			} else {
				alert(`Error connection : ${data.message}`)
			}
		} else {
			console.log(data.redirectPath);
			localStorage.setItem('token', data.token);
			window.location.href = data.redirectPath;
		}
	} catch (error) {
	}
}

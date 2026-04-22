import { check_token } from './check-token.js';

document.addEventListener('DOMContentLoaded', async () => {
	const r = await check_token();
	if (r == false) {
		window.location.href = '/';
	}
});

window.addEventListener('load', function () {
	
	document.getElementById("sendLink").onclick = async function () {
		var tmpUser = document.getElementById("Username")
		var tmpEmail = document.getElementById("Email")

		var login = tmpUser.value
		var email = tmpEmail.value

		if (login == "" || email == "") {
			alert("Please enter your email and your username")
		}
		// try {
			const response = await fetch("/sendResetLink", {
				method: "POST",
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify({
					"email": email,
					"username": login,
				})
			});
			// if (!response.ok) {
			// 	throw new Error(`HTTP error! status: ${response.status}`);
			// }

			// const data = await response.json();
			// if (response.status === 201) {
			// 	alert(`Account created successfully! Please check your email to verify your account.`);
			// 	window.location.href = data.redirectPath
			// }
			// else if (response.status === 409) {
			// 	alert("Username or email already in use.")
			// } else {
			// 	alert(`Error creating user: ${data.message}`);
			// }
		// } catch (error) {
		// 	console.error("Error: ", error);
		// 	alert(`HEEEREEEE Error creating user: ${error.message}`);
		// }
	}
});

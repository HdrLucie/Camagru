// TO DO LIST Fonction pour changer capturer les changements.
// Vérifier la validité des changements : si username - email OK. 
// Fonction pour récupérer le nouvel avatar.
// Fonction pour sauvegarder les nouvelles informations/refresh de la page avec les nouvelles info. 

let notifyState = initialNotifyState();


document.addEventListener('DOMContentLoaded', () => {
	checkToken();
	initialNotifyState();
});

function initialNotifyState() {
	const user = getUser();
	console.log(user.Notify);
	return user.notify;
}

function checkToken() {
	const token = localStorage.getItem('token');
	if (!token) {
		window.location.href = '/';
	}
	return token
}

function checkNotifyState() {
	const notifyState = document.querySelector('.checkbox-notify')
	
	return notifyState.checked;
}

function syncNotificationCheckbox(enabled) {
	const checkbox = document.querySelector('.checkbox-notify');
	if (checkbox) {
		checkbox.checked = enabled;

		checkbox.dispatchEvent(new Event('change'));
	}
	notifyState = enabled;
	console.log("État final checkbox:", enabled);
}

async function getModifications() {
	var tmpUser = document.getElementById("Username")
	var tmpEmail = document.getElementById("Email")
	var notifyState = checkNotifyState();
	var token = checkToken();

	var login = tmpUser.value
	var email = tmpEmail.value
	console.log(login, email);
		// try {
		const response = await fetch("/setUserDatas", {
			method: "POST",
			headers: {
                "Authorization": `Bearer ${token}`,
				"Content-Type": "application/json"
			},
			body: JSON.stringify({
				"email": email,
				"username": login,
				"notifyState": notifyState,
			})
		});
	// 	if (!response.ok) {
	// 		throw new Error(`HTTP error! status: ${response.status}`);
	// 	}
	//
	// 	const data = await response.json();
	// 	if (response.status === 201) {
	// 		console.log("\nAccount created successfully! Please check your email to verify your account.")
	// 		window.location.href = data.redirectPath
	// 	} else if (response.status === 409) {
	// 		alert("Username or email already in use.")
	// 	} else {
	// 		alert(`Error creating user: ${data.message}`);
	// 	}
	// } catch (error) {
	// 	console.error("Error: ", error);
	// }
}

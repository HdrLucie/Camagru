import { check_token } from './check-token.js';

document.addEventListener('DOMContentLoaded', () => {
	const r = check_token();
	if (r == false) {
		window.location.href = '/';
	}
	initialNotifyState();
});

let notifyState = initialNotifyState();

function initialNotifyState() {
	const user = getUser();
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
}

function emailIsValid (email) {
	if (email == "")
		return true
	return /\S+@\S+\.\S+/.test(email)
}

async function getUser() {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch("/getUser", {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
        });
        const userData = await response.json();
		return userData;
    } catch (error) {
        console.error("Erreur:", error);
        return null;
    }
}

document.getElementById('sendBtn').addEventListener('click', async () => {
	var tmpUser = document.getElementById("Username")
	var tmpEmail = document.getElementById("Email")
	var notifyState = checkNotifyState();
	var token = checkToken();
	var login = tmpUser.value
	let email = tmpEmail.value;
    if (!emailIsValid(email)) {
        alert("Wrong email");

        const user = await getUser();
        if (!user) {
            console.error("Impossible de récupérer l'utilisateur");
            return;
        }
        email = user.email;
    }
	try {
		const response = await fetch("/setUserDatas", {
			method: "POST",
			headers: {
                "Authorization": `Bearer ${token}`,
				"Content-Type": "application/json"
			},
			body: JSON.stringify({
				"email": email,
				"username": login,
				"notifyState": !notifyState,
			})
		});
		if (response.status === 200) {
			window.location.href = "/profile"
		} else if (response.status === 409) {
			alert("Username or email already in use.")
		}
	} catch (error) {
		console.error("Error: ", error);
	}
});

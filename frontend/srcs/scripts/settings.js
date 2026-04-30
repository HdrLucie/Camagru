import { check_token } from './check-token.js';
import { getUser } from './get-user.js';

document.addEventListener('DOMContentLoaded', async () => {
	const r = await check_token();
	if (r == false) {
		window.location.href = '/';
	}
    loadUserData();
});

function redirectionPage(path) {
    window.location.href = path;
}

const resetButton = document.getElementById("forgetPassword");
if (resetButton) {
	resetButton.onclick = function (event) {
		event.preventDefault();
		window.location.href = "/forgetPassword";
	};
}

async function loadUserData() {
    const userData = await getUser();
    if (!userData) return;
	document.querySelectorAll('[user-data]:not(img)').forEach(element => {
        const field = element.getAttribute('user-data');
        element.textContent = userData[field];
    });
	const emailInput = document.getElementById('Email');
	if (emailInput && userData.email) {
		emailInput.placeholder = userData.email;
	}
	const usernameInput = document.getElementById('Username');
	if (usernameInput && userData.username) {
		usernameInput.placeholder = userData.username;
	}
	if (userData.notify == false) {
		document.getElementById('checkedBox').checked = true;
	}
}

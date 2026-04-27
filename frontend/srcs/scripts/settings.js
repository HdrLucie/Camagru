import { check_token } from './check-token.js';

document.addEventListener('DOMContentLoaded', async () => {
	const r = await check_token();
	if (r == false) {
		window.location.href = '/';
	}
    loadUserData();
});

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

function redirectionPage(path) {
    window.location.href = path;
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


document.addEventListener('DOMContentLoaded', () => {
	checkToken();
    loadUserData();
});

async function checkToken() {
    const token = localStorage.getItem('token');
    if (!token) {
        // alert('No token found. Please login.');
        window.location.href = '/';
    }
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

function redirectionPage(path) {
    window.location.href = path;
}

async function loadUserData() {
    const userData = await getUser();
    if (!userData) return;
	console.log(userData);
	const avatarDiv = document.getElementById('avatarId');
	if (avatarDiv) {
        const img = document.createElement('img');
        img.src = "/stickers/" + userData.avatar;
		img.alt = userData.avatar;
        avatarDiv.appendChild(img);
	}
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


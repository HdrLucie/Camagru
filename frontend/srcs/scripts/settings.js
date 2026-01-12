document.addEventListener('DOMContentLoaded', () => {
	checkToken();
    loadUserData();
});

async function checkToken() {
    const token = localStorage.getItem('token');
    console.log("Function check token")
    console.log(token)
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
		console.log(userData);
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
}

document.addEventListener('DOMContentLoaded', function() {
    const photoButtons = document.querySelectorAll('[data-tab="photos"]');
    const settingsButtons = document.querySelectorAll('[data-tab="settings"]');
    const stickersButtons = document.querySelectorAll('[data-tab="stickers"]');
    const image_input = document.getElementById( 'image_input' );
    const photosList = document.getElementById('photos-list');
    const settings = document.getElementById('settings');
    const stickers = document.getElementById('stickers');
    function showContent(contentToShow) {
        [photosList, settings, stickers].forEach(content => {
            if (content) content.classList.remove('active');
        });
        if (contentToShow) contentToShow.classList.add('active');
    }
    photoButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            showContent(photosList);
        });
    });
    settingsButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            showContent(settings);
        });
    });
    stickersButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            showContent(stickers);
        });
    });
    image_input.addEventListener('click', function() {
	  const file_reader = new FileReader();
	  file_reader.addEventListener("load", () => {
	    const uploaded_image = file_reader.result;
	    document.querySelector("#display_image").style.backgroundImage = `url(${uploaded_image})`;
	  });
	  file_reader.readAsDataURL(this.files[0]);
	});

    showContent(photosList);
});

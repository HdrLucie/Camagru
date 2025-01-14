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
}

document.addEventListener('DOMContentLoaded', function() {
    const showPhotos = document.getElementById('show-photos');
    const showSettings = document.getElementById('show-settings');
    const showStickers = document.getElementById('show-stickers');
    const photosList = document.getElementById('photos-list');
    const settings = document.getElementById('settings');
    const stickers = document.getElementById('stickers');
	const image_input = document.getElementById( 'image_input' );

	function showContent(contentToShow) {
        [photosList, settings, stickers].forEach(content => {
            content.classList.remove('active');
        });
        contentToShow.classList.add('active');
    }

    showPhotos.addEventListener('click', function(e) {
        e.preventDefault();
        showContent(photosList);
    });

    showSettings.addEventListener('click', function(e) {
        e.preventDefault();
        showContent(settings);
    });

    showStickers.addEventListener('click', function(e) {
        e.preventDefault();
        showContent(stickers);
    });
    // La liste des photos est affichée par défaut
    showContent(photosList);

	image_input.addEventListener('click', function() {
	  const file_reader = new FileReader();
	  file_reader.addEventListener("load", () => {
	    const uploaded_image = file_reader.result;
	    document.querySelector("#display_image").style.backgroundImage = `url(${uploaded_image})`;
	  });
	  file_reader.readAsDataURL(this.files[0]);
	});
});

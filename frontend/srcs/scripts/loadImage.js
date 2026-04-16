window.addEventListener('load', function() {
	document.querySelector('input[type="file"]').addEventListener('change', async function() {
		if (this.files && this.files[0]) {
			const file = this.files[0];
			const allowedTypes = ['image/jpeg', 'image/png', 'image/gif'];
			if (!allowedTypes.includes(file.type)) {
				alert('Only JPG, PNG or GIF format');
				return;
			}
			const resizeImg = document.getElementById('uploadedPhoto');
			resizeImg.src = URL.createObjectURL(file);
		}
	});
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


async function sendImg() {
    const token = localStorage.getItem('token');
	const stickers = [];
    const fileInput = document.getElementById('uploadPhoto');
    if (!fileInput.files || fileInput.files.length === 0) {
        alert('Veuillez sélectionner une photo avant d\'envoyer !');
        return;
    }

    const elements = document.querySelectorAll('.placed-sticker');
    if (!elements || elements.length === 0) {
        alert('Veuillez placer un sticker sur votre photo avant d\'envoyer !');
        return;
    }

    const user = await getUser();

    const formData = new FormData();
    formData.append('image', fileInput.files[0]);
    formData.append('id', user.id);
	elements.forEach(sticker => {
		stickers.push({
			path: sticker.src,
			posX: Math.floor(sticker.dataset.relativeX),
			posY: Math.floor(sticker.dataset.relativeY),
			id: parseInt(sticker.id) 
		});
	});
	formData.append('stickers', JSON.stringify(stickers));
	formData.append('timestamp', new Date().toISOString());
	const response = await fetch("/sendImage", {
		method: "POST",
		headers: {
			"Authorization": `Bearer ${token}`
		},
		body: formData
	});
	window.location.href = `/gallery/1`;
}

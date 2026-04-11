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

    const fileInput = document.getElementById('uploadPhoto');
    if (!fileInput.files || fileInput.files.length === 0) {
        alert('Veuillez sélectionner une photo avant d\'envoyer !');
        return;
    }

    const sticker = document.getElementsByClassName('placed-sticker');
    if (!sticker || sticker.length === 0) {
        alert('Veuillez placer un sticker sur votre photo avant d\'envoyer !');
        return;
    }

    const user = await getUser();

    const formData = new FormData();
    formData.append('image', fileInput.files[0]);
    formData.append('id', user.id);
    formData.append('imageId', sticker[0].id);
    const relativeX = sticker[0].dataset.relativeX;
    const relativeY = sticker[0].dataset.relativeY;
    formData.append('stickerPath', sticker[0].src);
    formData.append('posX', JSON.stringify(Math.floor(relativeX)));
    formData.append('posY', JSON.stringify(Math.floor(relativeY)));
    formData.append('timestamp', new Date().toISOString());

    const response = await fetch("/sendImage", {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${token}`
        },
        body: formData
    });
}

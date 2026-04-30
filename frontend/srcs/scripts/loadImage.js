import { check_token } from './check-token.js';
import { getUser } from './get-user.js';

document.addEventListener('DOMContentLoaded', async () => {
	const r = await check_token();
	if (r == false) {
		window.location.href = '/';
	}
});

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
			document.querySelector('.image-upload-wrap').style.display = 'none';
			document.querySelector('.file-upload-content').style.display = 'block';
		}
	});
});

document.getElementById('removeBtn').addEventListener('click', removeUpload);

function removeUpload() {
    document.querySelector('.image-upload-wrap').style.display = 'block';
    document.querySelector('.file-upload-content').style.display = 'none';
    document.getElementById('uploadedPhoto').src = '#';
    document.querySelectorAll('.placed-sticker').forEach(s => s.remove());

}

document.getElementById('sendImg').addEventListener('click', async () => {
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
	if (!response.ok) {
		throw new Error(`Response status: ${response.status}`);
	}
	window.location.href = `/gallery/1`;
});

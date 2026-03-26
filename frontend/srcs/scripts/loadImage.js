document.addEventListener('DOMContentLoaded', () => {
	checkToken();
	displayStickers();
});

async function checkToken() {
	const token = localStorage.getItem('token');
	if (!token) {
		// alert('No token found. Please login.');
		window.location.href = '/';
	}
}

window.addEventListener('load', function() {
  document.querySelector('input[type="file"]').addEventListener('change', function() {
      if (this.files && this.files[0]) {
          var img = document.getElementById('uploadedPhoto');
          img.onload = () => {
              URL.revokeObjectURL(img.src);  // no longer needed, free memory
          }

          img.src = URL.createObjectURL(this.files[0]); // set src to blob url
      }
  });
});

const draggableItem = document.getElementById("draggableItem");

async function getStickers() {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch("/getStickers", {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
        });
        const stickers = await response.json();
        return stickers;
    } catch (error) {
        console.error("Erreur:", error);
        return null;
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



async function sendPictures() {
	var path;
	const token = localStorage.getItem('token');
	const sticker = document.getElementsByClassName('placed-sticker');
	if (!sticker || sticker.length === 0) {
		alert('Veuillez placer un sticker sur votre photo avant d\'envoyer !');
		return;
	}
	const user = await getUser();
    const uploadedImg = document.getElementById('uploadPhoto');
    const imgBlob = uploadedImg.files[0];
    if (!imgBlob) {
        alert('Veuillez choisir une image avant d\'envoyer !');
        return;
    }
	console.log(imgBlob);
	const formData = new FormData();
	formData.append('id', user.id);
	formData.append("image", imgBlob, 'photo.jpg');
	// formData.append('imageId', sticker[0].id);
	formData.append('imageId', sticker[0].dataset.stickerId);
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
	if (!response.ok) {
    const errorText = await response.text();
    console.error('Erreur serveur:', errorText);
    alert('Erreur: ' + errorText);
    return;
}
}

async function getStickerNameById(stickerId) {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch(`/getSticker/${stickerId}`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            }
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const sticker = await response.json();
        return sticker.name;
    } catch (error) {
        console.error("Error while getting sticker:", error);
        return null;
    }
}

async function createStickerOnImage(stickerId, x, y) {
    const existingSticker = document.querySelector('.placed-sticker');
    if (existingSticker) {
		existingSticker.remove();
	}
    const name = await getStickerNameById(stickerId);
    const stickerElement = document.createElement('img');
    stickerElement.src = "/stickers/" + name;
    stickerElement.className = 'placed-sticker';
	stickerElement.dataset.stickerId = stickerId;
    const container = document.getElementById('uploadedImageContainer');

    const percentX = (x / container.offsetWidth) * 100;
    const percentY = (y / container.offsetHeight) * 100;

    stickerElement.style.position = 'absolute';
    stickerElement.style.left = percentX + '%';
    stickerElement.style.top = percentY + '%';
    stickerElement.style.width = '80px';
    stickerElement.style.height = '80px';
    stickerElement.style.transform = 'translate(-50%, -50%)';
    stickerElement.dataset.relativeX = percentX;
    stickerElement.dataset.relativeY = percentY;

    container.appendChild(stickerElement);
}

function handleDrop(event) {
	event.preventDefault();
	const id = event
	.dataTransfer
	.getData('text/plain');
	
	const container = document.getElementById('uploadedImageContainer');
    const rect = container.getBoundingClientRect();
    const dropX = event.clientX - rect.left;
    const dropY = event.clientY - rect.top;
    createStickerOnImage(id, dropX, dropY);
}

function handleDragOver(event) {
  event.preventDefault();
}

function handleDragStart(event) {
	event
    .dataTransfer
    .setData('text/plain', event.target.id);
}

async function displayStickers() {
    const stickers = await getStickers();
    const container = document.getElementById('stickersContainer');
    container.innerHTML = '';
    if (stickers && stickers.length > 0) {
        stickers.forEach(sticker => {
            const img = document.createElement('img');
            img.src = "/stickers/" + sticker.name;
			img.alt = sticker.name;
			img.id = sticker.id;
			img.draggable = true;
			img.className = 'sticker-item';
			img.addEventListener('dragstart', handleDragStart);
            container.appendChild(img);
        });
    } else {
        const message = document.createElement('p');
        message.textContent = 'No sticker available';
        container.appendChild(message);
    }
}

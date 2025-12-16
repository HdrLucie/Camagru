document.addEventListener('DOMContentLoaded', () => {
    displayAvatars();
});

async function getAvatars() {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch("/getAvatars", {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
        });
        const avatars = await response.json();
        return avatars;
    } catch (error) {
        console.error("Erreur:", error);
        return null;
    }
}

async function displayStickers() {
    const avatars = await getStickers();
    const container = document.getElementById('avatarsContainer');
    container.innerHTML = '';
    if (avatars && avatars.length > 0) {
        avatars.forEach(avatar => {
            const img = document.createElement('img');
            img.src = "/avatars/" + avatar.name;
			img.alt = avatar.name;
			img.id = avatar.id;
			img.draggable = true;
			img.className = 'avatar-item';
			img.addEventListener('dragstart', handleDragStart);
            container.appendChild(img);
        });
    } else {
        const message = document.createElement('p');
        message.textContent = 'No avatar available';
        container.appendChild(message);
    }
}

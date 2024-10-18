document.addEventListener('DOMContentLoaded', () => {
    displayStickers();
});

async function getStickers() {
    const token = localStorage.getItem('token');
    console.log("getStickers function")    
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

async function displayStickers() {
    const stickers = await getStickers();
    const container = document.getElementById('stickersContainer');
    
    // Nettoyer le conteneur
    container.innerHTML = '';
    
    if (stickers && stickers.length > 0) {
        stickers.forEach(sticker => {
            const img = document.createElement('img');
            img.src = "/stickers/" + sticker.name;
            img.alt = sticker.name;
            container.appendChild(img);
        });
    } else {
        const message = document.createElement('p');
        message.textContent = 'Aucun sticker disponible';
        container.appendChild(message);
    }
}
import { check_token } from './check-token.js';

document.getElementById("logout").onclick = async function () {
	const r = await check_token();
	if (r == false) {
		return ;
	}
	const token = localStorage.getItem('token')
	try {
		const response = await fetch("/logout", {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"Authorization": `Bearer ${token}`,
			},
		});
        const responseText = await response.text();
        let data;
        if (responseText) {
            data = JSON.parse(responseText);
        }
		if (!response.ok) {
			localStorage.clear()
			window.location.href = '/'
		}
		if (response.status === 200) {
			localStorage.clear()
			window.location.href = '/'
		} else {
			alert(`Error logout message: ${data.message}`);
		}
	} catch (error) {
		console.error("Error: ", error);
		alert(`Error logout: ${error.message}`);
	}
}

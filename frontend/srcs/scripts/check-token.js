export async function check_token() {
	const token = localStorage.getItem('token');
	if (token == null) {
		return false;
	}
	try {
		const response = await fetch("/getToken", {
			method: "GET",
			headers: {
				"Authorization": `Bearer ${token}`,
			},
		});
		if (!response.ok) {
			throw new Error(`Response status: ${response.status}`);
		}
		if (response.status === 401) {
			localStorage.removeItem('token');
			return false;
		}
		return true;
	} catch (error) {
		console.error(error);
	}
}

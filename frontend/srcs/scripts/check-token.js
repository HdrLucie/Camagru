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
		if (response.status === 401) {
			localStorage.removeItem('token');
			return false;
		}
		return true;
	} catch (error) {
		console.error(error);
	}
}

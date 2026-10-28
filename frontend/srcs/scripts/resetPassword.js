
document.addEventListener('DOMContentLoaded', async () => {
});

function checkPassword(password) {
	let regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,25}$/;
	let result = regex.test(password);
	return result
}

window.addEventListener('load', function () {
	
    function getEmail(str, char) {
        const parts = str.split(char);
        if (parts.length > 1) {
            return parts.slice(1).join(char);
        }
        return '';
    }

	document.getElementById("resetPassword").onclick = async function () {
		var firstPswd = document.getElementById("firstPassword")
		var secondPswd = document.getElementById("secondPassword")
        var pathName = window.location.search

        var email = getEmail(pathName, '=')
		var pswd1 = firstPswd.value
		var pswd2 = secondPswd.value

		if (pswd1 == "" || pswd2 == "" || pswd1 != pswd2) {
			alert("Please make sure both passwords are identical.")
            return 
		}
		if (checkPassword(pswd1) == false || checkPassword(pswd2) == false) {
			alert("Please make sure the format is right.")
            return 
		}

		try {
			const response = await fetch("/newPassword", {
				method: "POST",
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify({
					"email": email,
					"password": pswd1,
				})
			});
			const data = await response.json();
			if (response.status === 200) {
				alert(`Password successfully changed`);
				window.location.href = data.redirectPath
			}
		} catch (error) {
			alert(`Oops, it seems we don’t know each other! Feel free to create an account to join us!`);
		}
	}
});

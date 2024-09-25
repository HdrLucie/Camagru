window.addEventListener('load', function () {
	console.log('La page entière est complètement chargée');
	
    function getEmail(str, char) {
        console.log("GetEmail")
        const parts = str.split(char);
        if (parts.length > 1) {
            return parts.slice(1).join(char);
        }
        return '';
    }

	document.getElementById("resetPassword").onclick = async function () {
		console.log("\n\nRESET PASSWORD\n\n")
		var firstPswd = document.getElementById("firstPassword")
		var secondPswd = document.getElementById("secondPassword")
        var pathName = window.location.search

        var email = getEmail(pathName, '=')
		var pswd1 = firstPswd.value
		var pswd2 = secondPswd.value

        console.log(pswd1 + pswd2 + email)
		if (pswd1 == "" || pswd2 == "" || pswd1 != pswd2) {
			alert("Please make sure both passwords are identical.")
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
			// if (!response.ok) {
			// 	throw new Error(`HTTP error! status: ${response.status}`);
			// }

			const data = await response.json();
			if (response.status === 200) {
				alert(`Password successfully changed`);
				window.location.href = data.redirectPath
			}
			else if (response.status === 404) {
				alert("Error : User not found")
			} else if (response.status === 400) {
				alert("Error : Impossible to change password");
			}
		} catch (error) {
			console.error("Error: ", error);
			alert(`Oops, it seems we don’t know each other! Feel free to create an account to join us!`);
		}
	}
});
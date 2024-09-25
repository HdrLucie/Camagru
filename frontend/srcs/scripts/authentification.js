window.onload = verifyAccount;

function getToken(str, char) {
    console.log("GetToken")
    const parts = str.split(char);
    if (parts.length > 1) {
        return parts.slice(1).join(char);
    }
    return '';
}

async function verifyAccount() {
    var pathName = window.location.search
    var token = getToken(pathName, '=')
    console.log(pathName)
    console.log(token)
    console.log("Verify Account")
    try {
        const response = await fetch("/verifyAccount", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "token": token,
            })
        });
        const responseText = await response.text();
        let data;
        if (responseText) {
            data = JSON.parse(responseText);
        }
        console.log(response.status)
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        if (response.status === 201) {
            alert("Profile confirmed");
            localStorage.clear()
            window.location.href = data.redirectPath
        } else {
            alert(`Error logout message: ${data.message}`);
        }
    } catch (error) {
        console.error("Error: ", error);
        alert(`Error logout: ${error.message}`);
    }
}
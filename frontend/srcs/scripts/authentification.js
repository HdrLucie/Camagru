window.onload = verifyAccount;

function getToken(str, char) {
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

    const response = await fetch("/verifyAuth", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "token": token,
        })
    });
}
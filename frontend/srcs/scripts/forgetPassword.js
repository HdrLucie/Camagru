document.addEventListener('DOMContentLoaded', async () => {
    document.getElementById("sendLink").addEventListener("click", async function () {
        const login = document.getElementById("Username").value;
        const email = document.getElementById("Email").value;

        if (!login || !email) {
            alert("Please enter your email and your username");
            return;
        }

        try {
			console.log("send reset link");
            const response = await fetch("/sendResetLink", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, username: login })
            });

            const data = await response.json();

            if (response.ok) {
                alert("A reset link has been sent to your email.");
            } else {
                alert(`Error: ${data.message}`);
            }
        } catch (error) {
            console.error("Error: ", error);
            alert("Something went wrong, please try again.");
        }
    });
});

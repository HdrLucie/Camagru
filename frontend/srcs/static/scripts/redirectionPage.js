document.addEventListener("DOMContentLoaded", function() {
    console.log("Je passe bien par cette fonction.\n")
    button = document.getElementById("navigateButton");
    button.addEventListener("click", function() {
        window.location.href = "/register";
    });
});

// TO DO LIST Fonction pour changer capturer les changements.
// Vérifier la validité des changements : si username - email OK. 
// Fonction pour récupérer le nouvel avatar.
// Fonction pour sauvegarder les nouvelles informations/refresh de la page avec les nouvelles info. 

function getModifications() {
	var tmpUser = document.getElementById("Username")
	var tmpEmail = document.getElementById("Email")

	var login = tmpUser.value
	var email = tmpEmail.value

	console.log(login, email);	
}

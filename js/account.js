window.addEventListener("load", (_) => {
	if (!(
		sessionStorage.getItem("sessionId")
		&& sessionStorage.getItem("userId")
		&& sessionStorage.getItem("name")
	)) {
		alert("you're not logged in - log in to play")
		logout()
	}
})

function logout() {
	sessionStorage.removeItem("sessionId")
	sessionStorage.removeItem("userID")
	sessionStorage.removeItem("name")
	window.open(`http://${window.location.host}`, "_self")
}

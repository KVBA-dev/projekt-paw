function create() {
	fetch("/create", {
		method: "POST",
	}).then(r => {
		if (!r.ok) {
			throw "something went wrong upon game creation"
		}
		return r.text()
	}).then(id => {
		console.log(`game created with id ${id}`)
		sessionStorage.setItem("gameid", id)
		window.open(`http://${window.location.host}/game?id=${id}`, "_self")
	})
}

function join(gameid) {
	sessionStorage.setItem("gameid", gameid)
	window.open(`http://${window.location.host}/game?id=${gameid}`, "_self")
}

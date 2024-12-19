function loginAnonymous(id_username) {
	const username = document.getElementById(id_username).value;
	let id, name;
	fetch("/loginanon", {
		method: "POST",
		body: {
			"username": username
		}
	}).then(r => {
		if (r.status >= 400) {
			throw "got an error upon login";
		}
		return r.json();
	}).then(r => {
		id = r.id;
		name = r.name;
		localStorage.setItem("id", id);
		localStorage.setItem("name", username);
		window.open("http://localhost:8080/list", "_self");
	});
}

function loginAnonymous(id_username) {
	const username = document.getElementById(id_username).value;
	let id, name;
	// NOTE: lmao
	fetch(`/loginanon?username=${username}`, {
		method: "POST",
	}).then(r => {
		if (r.status >= 400) {
			throw "got an error upon login";
		}
		return r.json();
	}).then(r => {
		id = r.id;
		name = r.name;
		sessionStorage.setItem("id", id);
		sessionStorage.setItem("name", name);
		window.open(`http://${window.location.host}/list`, "_self");
	});
}

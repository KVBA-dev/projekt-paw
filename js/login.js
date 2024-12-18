async function loginAnonymous(id_username) {
	const username = document.getElementById(id_username).value;
	let id = await fetch("/loginanon", {
		method: "POST",
		body: {
			"username": username
		}
	}).then(r => {
		console.log(r.status)
		if (r.status >= 400) {
			console.log("coś się zepsuło :/");
			throw "got an error upon login";
		}
		return r.json();
	}).then(r => r.id);
	console.log("id set");
	localStorage.setItem("id", id);
	console.log("id in local storage");
	window.open("http://localhost:8080/list", "_self");
}

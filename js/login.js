function loginAnonymous(id_username) {
	const username = document.getElementById(id_username).value;
	// why, javascript? why?
	if (!username) {
		document.getElementById("login-error").innerText = "Something went wrong"
		return
	}
	fetch(`/loginanon?username=${username}`, {
		method: "POST",
	}).then(r => {
		if (r.status >= 400) {
			document.getElementById("login-error").innerText = "Something went wrong"
			throw "got an error upon login";
		}
		return r.json();
	}).then(r => {
		sessionStorage.setItem("sessionId", r.session_id);
		sessionStorage.setItem("name", username);
		sessionStorage.setItem("userId", -1);
		window.open(`http://${window.location.host}/list`, "_self");
	}).catch(_ => { })
}

function login(id_login, id_password) {
	const username = document.getElementById(id_login).value
	const password = document.getElementById(id_password).value

	if (!username || !password) {
		document.getElementById("login-error").innerText = "Something went wrong"
		return
	}

	if (!/[A-Za-z0-9-_]+/.test(username)) {
		document.getElementById("login-error").innerText = "Something went wrong"
		return
	}

	fetch(`/login`, {
		method: "POST",
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			username: username,
			password: password
		})
	}).then(r => {
		if (r.status >= 400) {
			document.getElementById("login-error").innerText = "Something went wrong"
			throw "got an error upon login";
		}
		return r.json()
	}).then(r => {
		sessionStorage.setItem("sessionId", r.session_id)
		sessionStorage.setItem("userId", r.user_id)
		sessionStorage.setItem("name", r.username)
		window.open(`http://${window.location.host}/list`, "_self");
	}).catch(_ => { })
}

function register(id_login, id_password) {
	const username = document.getElementById(id_login).value
	const password = document.getElementById(id_password).value

	if (!username || !password) {
		document.getElementById("login-error").innerText = "Something went wrong"
		return
	}

	if (!/[A-Za-z0-9-_]+/.test(username)) {
		document.getElementById("login-error").innerText = "Something went wrong"
		return
	}


	fetch(`/register`, {
		method: "POST",
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			username: username,
			password: password
		})
	}).then(r => {
		if (r.status >= 400) {
			document.getElementById("login-error").innerText = "Something went wrong"
			throw "got an error upon login";
		}
		return r.json()
	}).then(r => {
		sessionStorage.setItem("sessionId", r.session_id)
		sessionStorage.setItem("userId", r.user_id)
		sessionStorage.setItem("name", r.username)
		window.open(`http://${window.location.host}/list`, "_self");
	}).catch(_ => { })
}

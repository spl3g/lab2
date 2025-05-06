// AUTH
let username = "";
let password = "";

let refresh_token = "";
let connected = false;

const btnLogin = document.getElementById("btnLogin")
const btnLogout = document.getElementById("btnLogout")

const headers = new Headers();
headers.append("Content-Type", "application/x-www-form-urlencoded");

const getToken = async () => {
	console.log("get token called");

	const params = new URLSearchParams();
	params.append("client_id", "frontend");

	const requestOptions = {
		method: "POST",
		headers: headers,
		body: params,
	};

	if (refresh_token == "") {
		params.append("username", username);
		params.append("password", password);
		params.append("grant_type", "password");
	} else {
		params.append("refresh_token", refresh_token);
		params.append("grant_type", "refresh_token");
	}

	return fetch(
		"http://127.0.0.1:8090/realms/sirius/protocol/openid-connect/token",
		requestOptions
	)
		.then((response) => response.json())
		.then((result) => {
			console.log("token accepted");
			refresh_token = result.refresh_token;

			return result.access_token;
		})
		.catch((error) => console.error(error));
};

// CHAT
const output = document.getElementById("output")

const btnSubscribe = document.getElementById("btnSubscribe")
const btnSend = document.getElementById("btnSend")

const message = document.getElementById("message")
const channel = document.getElementById("channel")

// CENTRIFUGE
const client = new Centrifuge(
	"ws://127.0.0.1:8080/centrifugo/connection/websocket",
	{
		getToken: getToken,
		debug: true,
	}
);

client.on("connected", () => {
	password = null;

	document.getElementById("password").value = "";

	connected = true;
});

client.on("disconnected", () => {
	refresh_token = "";
	client.setToken("");

	connected = false;
});

// HANDLERS
btnLogin.addEventListener("click", () => {
	if (connected) {
		return;
	}

	username = document.getElementById("username").value;
	password = document.getElementById("password").value;

	if (username != "" && password != "") {
		client.connect();
	} else {
		console.error("empty username or password");
	}
});

btnLogout.addEventListener("click", () => {
	if (!connected) {
		return;
	}

	console.log("disconnecting");

	const params = new URLSearchParams();
	params.append("client_id", "frontend");

	const requestOptions = {
		method: "POST",
		headers: headers,
		body: params,
	};

	if (refresh_token != "") {
		params.append("refresh_token", refresh_token);
		params.append("grant_type", "refresh_token");

		fetch(
			"http://127.0.0.1:8090/realms/sirius/protocol/openid-connect/logout",
			requestOptions
		)
			.then((response) => {
				console.log("session closed");
			})
			.catch((error) => console.error(error))
			.finally(() => {
				client.disconnect();
			});
	} else {
		console.error("no refresh token");
	}
});

btnSubscribe.addEventListener("click", () => {
	// if (!connected) {
	// 	return;
	// }

	console.log("Subscribe: " + channel.value)
	const sub = client.newSubscription(channel.value);

	sub.on("publication", (msg) =>
	{
		let line = `<p><strong>${msg.data.from}:</strong> ${msg.data.message}</p>`
		output.innerHTML += line
	})


	sub.on('subscribed', function(ctx) {
		console.log('subscribed');
	});

	sub.on('unsubscribed', function(ctx) {
		console.log('unsubscribed');
	});

	sub.subscribe();

});


btnSend.addEventListener("click", () => {
	// if (!connected) {
	// 	return;
	// }

	if (message.value != "")
	{
		client.publish(channel.value, {
			from: username,
			message: message.value
		})

		message.value = ""
	}

	console.log("Send: " + channel)
});



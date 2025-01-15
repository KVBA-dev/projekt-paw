const delim = '»'
const gameid = sessionStorage.getItem("gameid")

//  WARN: at least it works
let socket

function game() {

	socket = new WebSocket(`ws://${window.location.host}/gamews?game=${gameid}`)
	socket.addEventListener("open", () => {
		const name = sessionStorage.getItem("name")
		const id = sessionStorage.getItem("sessionId")
		socket.send(`d${delim}${name}${delim}${id}`)
	})

	$('.btn-flag').click(buttonCheck)

	$('#selected-container, #yes-no, #question-input-container').hide()

	socket.addEventListener("message", (e) => {
		console.log(e.data)
		const msg = e.data.split(delim).slice(1)
		switch (e.data[0]) {
			case "w":
				toggleFlagButtons(false)
				break
			case "s":
				// NOTE: select flag
				toggleFlagButtons(true)
				$('.btn-flag').removeClass('unchecked')
				$('.btn-flag').unbind('click')
				$('.btn-flag').click(function() {
					const id = $(this).attr('data-id')
					socket.send(`s${delim}${id}`)

					toggleFlagButtons(false)
					$('#selected').html($(this).html())
					$('#selected-container').show()

					$('.btn-flag').unbind('click')
					$('.btn-flag').click(buttonCheck)
				})
				break
			case "o":
				// NOTE: start game
				toggleFlagButtons(true)
				$('.btn-flag').toggleClass("checked")
				break
			case "t":
				// NOTE: start turn
				toggleFlagButtons(true)
				$('#question-input-container').show()
				$('#question-input').text("")

				break
			case "a":
				// NOTE: receive question
				$('#yes-no').show()
				$('#question').text(`Opponent asked: ${msg[0]}`)
				break
			case "r":
				// NOTE: receive answer
				break
			case "q":
				// NOTE: game over
				toggleFlagButtons(false)
				log("Game over")
				break
			case "l":
				// NOTE: log message
				log(msg)
				break
			case "x":
				// NOTE: game error - exit game
				alert(msg)
				socket.close()
			default:
				break
		}
	})

	socket.addEventListener("close", (_) => {
		alert("game closed")
		history.back()
	})

	window.addEventListener("beforeunload", (_) => {
		socket.close()
	})
}

function toggleFlagButtons(enabled) {
	$(".btn-flag").prop('disabled', !enabled)
}

function buttonCheck() {
	$(this).toggleClass('unchecked')
	$(this).toggleClass('checked')
}

function log(msg) {
	$('#gamelog').append(`${msg}\n`)
	const text = $('#gamelog').text()
	$('#gamelog').text(text.split('\n').slice(-7).join('\n'))
}

function sendQuestion() {
	const text = $('#question-input').val()
	socket.send(`a${delim}${text}`)
	$('#question-input-container').hide()
}

function reply(r) {
	alert(typeof r)
	let s = "n"
	// i hate javascript
	if (r) {
		s = "y"
	}
	socket.send(`r${delim}${s}`)
	$('#yes-no').hide()
}

function prepareGuess() {
	$('.btn-flag').unbind('click')
	$('.btn-flag').click(function() {
		const name = $(this).children('span').text()
		const id = $(this).attr('data-id')

		socket.send(`g${delim}${id}`)

		$('.btn-flag').unbind('click')
		$('.btn-flag').click(buttonCheck)
	})
	$('#question-input-container').hide()
}


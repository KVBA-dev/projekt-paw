const delim = '»'
const gameid = sessionStorage.getItem("gameid")
const socket = new WebSocket(`ws://${window.location.host}/gamews?game=${gameid}`)

function game() {
	socket.addEventListener("open", () => {
		const name = sessionStorage.getItem("name")
		const id = sessionStorage.getItem("id")
		socket.send(`d${delim}${name}${delim}${id}`)
	})

	$('.btn-flag').click(buttonCheck)

	$('#selected-container, #yes-no, #question-input-container').hide()

	socket.addEventListener("message", (e) => {
		const msg = e.data.split(delim).slice(1)
		switch (e.data[0]) {
			case "w":
				toggleFlagButtons(false)
				break
			case "s":
				// NOTE: select character
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
				break
			case "t":
				// NOTE: start turn
				toggleFlagButtons(true)
				$('#question-input-container').show()

				break
			case "a":
				// NOTE: receive question
				log(msg[0])
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
}

function log(msg) {
	$('#gamelog').append(`${msg}<br/>`)
}

function sendQuestion() {
	const text = $('#question-input').val()
	socket.send(`a${delim}${text}`)
	$('#question-input-container').hide()
}

function reply(r) {
	let s = "n"
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


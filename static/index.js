const ta = document.querySelector("#regta")
const searchta = document.querySelector("#searchta")
const rec = new Recorder()

rec.track(ta)
rec.track(searchta)

let recs = []

const app = new Vue({
	"el": "#app",
	"data": {
		"state": "main",
		"name": "",
		"bestMatchName": "",
		"searchRes": [
			{
				"id": "",
				"name": "",
				"propMatch": 0,
				"confidence": 0,
				"dist": 0
			}
		]
	},
	"methods": {
		"resetRec": function () {
			if (this.state === "register") {
				document.querySelector("#registerdiv").style.display = "block"
			} else if (this.state === "search") {
				document.querySelector("#searchdiv").style.display = "block"
			}

			recs = []

			rec.reset()
		},
		"next": async function () {
			const exp = rec.exportV1()

			console.log(exp)

			recs.push(exp)

			rec.reset()

			regta.value = ""

			if (recs.length >= 3) {
				// Done. Onto name.
				document.querySelector("#registerdiv").style.display = "none"
				this.state = "getName"
			}
		},
		"submitName": async function () {
			if (this.name.length < 1) {
				alert("Please provide a name.")
			}

			this.state = "loading"

			const res = await fetch(window.location.protocol + "//" + window.location.host + "/upload", {
				"method": "POST",
				"body": JSON.stringify({recs, "name": this.name})
			})

			recs = []
			this.name = ""

			if (res.status === 201) {
				this.state = "doneUpload"
			} else {
				this.state = "main"
				alert("Error: " + await res.text())
			}
		},
		"submitSearch": async function () {
			document.querySelector("#searchdiv").style.display = "none"
			searchta.value = ""

			this.state = "loading"

			const exp = rec.exportV1()
			rec.reset()

			console.log(exp)

			const res = await fetch(window.location.protocol + "//" + window.location.host + "/search", {
				"method": "POST",
				"body": JSON.stringify({"rec": exp})
			})

			if (res.status === 200) {
				// Load data to app

				this.searchRes.splice(0, this.searchRes.length)

				const parsed = await res.json()

				parsed.results.forEach((r) => {
					this.searchRes.push(r)
				})

				// Find best match

				const bestMatch = parsed.results.sort((a, b) => a.dist < b.dist ? -1 : 1)[0]

				this.bestMatchName = bestMatch.name

				// Done

				this.state = "doneSearch"
			} else {
				this.state = "main"
				alert("Error: " + await res.text())
			}
		}
	}
})

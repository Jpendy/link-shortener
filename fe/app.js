const form = document.querySelector('form')
const a = document.getElementById('short-link-a-tag')

form.addEventListener('submit', (e) => {
    e.preventDefault()
    const formData = new FormData(form)
    const fullLink = formData.get('full-link')

    return fetch("/shorten", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ fullLink })
    })
        .then(res => res.json())
        .then(({ shortLink }) => {
            console.log(shortLink)
            a.textContent = shortLink;
            a.href = shortLink
        })

})
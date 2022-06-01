const form = document.querySelector('form')
const input = document.querySelector('input')
const a = document.querySelector('a')
const errorBox = document.querySelector('#error-message')

form.addEventListener('submit', (e) => {
    e.preventDefault()
    const formData = new FormData(form)
    const fullLink = formData.get('full-link').trim()

    if (!fullLink.startsWith('http://') && !fullLink.startsWith('https://')) {
        errorBox.textContent = 'Please prefix link with "http://" or "https://"'
        return;
    }

    return fetch("/shorten", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ fullLink })
    })
        .then(res => res.json())
        .then(({ shortLink }) => {
            const href = shortLink.split('/')[1]
            a.textContent = shortLink;
            a.href = href
        })

})

input.addEventListener('input', () => {
    a.textContent = '';
    errorBox.textContent = '';
})
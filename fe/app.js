const form = document.querySelector('form')
const input = document.querySelector('input')
const a = document.querySelector('a')
const errorBox = document.querySelector('#error-message')
const linkBox = document.querySelector('#link-box')
const copyButton = document.querySelector('#copy-button')
const tooltip = document.querySelector('#tooltiptext')

form.addEventListener('submit', (e) => {
    e.preventDefault()
    const formData = new FormData(form)
    const fullLink = formData.get('full-link').trim()


    const lowerLink = fullLink.toLowerCase()
    if (!lowerLink.startsWith('http://') && !lowerLink.startsWith('https://')) {
        errorBox.textContent = 'Please prefix link with "http://" or "https://"'
        return;
    }

    return fetch("/shorten", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ fullLink: lowerLink })
    })
        .then(res => res.json())
        .then(({ shortLink }) => {
            tooltip.textContent = 'Copy to clipboard'
            copyButton.classList.remove('hidden')
            const href = shortLink.split('/')[1]
            a.textContent = 'www.' + shortLink;
            a.href = href
        })

})

copyButton.addEventListener('click', () => {
    navigator.clipboard.writeText(a.textContent)
    tooltip.textContent = 'Copied!'
})

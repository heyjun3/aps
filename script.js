async function scrollBottom(count=100, intervalSec=2000) {
    for(let i = 0; i < count; i++) {
        window.scrollTo(0, document.documentElement.scrollHeight);
        await new Promise((resolve, reject) => setTimeout(resolve, intervalSec))
    }
}

async function openTab(url, intervalSec = 500) {
    const tab = window.open(url, '_blank');
    await new Promise((resolve, reject) => setTimeout(resolve, intervalSec))
    tab.close()
    await new Promise((resolve, reject) => setTimeout(resolve, intervalSec))
}

async function main() {
    const elems = document.querySelectorAll('.link-area.needsclick.append-anchor')
    for (const elem of elems) {
        await openTab(elem.href)
    }
}

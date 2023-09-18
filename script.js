async function scrollBottom(count=100, intervalSec=500) {
    for(let i = 0; i < count; i++) {
        window.scrollTo(0, document.documentElement.scrollHeight);
        await new Promise((resolve, reject) => setTimeout(resolve, intervalSec))
    }
}

async function openTab(url, intervalSec = 250) {
    const tab = window.open(url, '_blank');
    await new Promise((resolve, reject) => setTimeout(resolve, intervalSec))
    return tab
}

async function main() {
    const elems = document.querySelectorAll('.link-area.needsclick.append-anchor')
    const tabs = []
    for (const elem of elems) {
        const tab = await openTab(elem.href)
        if (tab) {
            tabs.push(tab)
        }
        if (tabs.length > 10) {
            tabs.map((tab) => tab.close())
            tabs.splice(0)
        }
    }
    tabs.map((tab) => tab.close());
}

main();

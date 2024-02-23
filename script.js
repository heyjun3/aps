function sleep(intervalSec) {
  return new Promise((resolve) => setTimeout(resolve, intervalSec));
}

async function scrollBottom(count = 100, intervalSec = 2000) {
  for (let i = 0; i < count; i++) {
    window.scrollTo(0, document.documentElement.scrollHeight);
    await sleep(intervalSec);
  }
}

async function openTab(url, intervalSec = 500) {
  const tab = window.open(url, "_blank");
  await sleep(intervalSec);
  tab.close();
  await sleep(intervalSec);
}

async function main() {
  await scrollBottom();
  const elems = document.querySelectorAll(
    ".link-area.needsclick.append-anchor"
  );
  for (const elem of elems) {
    await openTab(elem.href);
  }
}

main();

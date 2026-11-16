// Measures a relation-map page built with the relation-map second-opinion preset (REQ-737).
// One input: the page URL. Prints one JSON document on stdout.
//
// Run (zsh, from any folder, with a local server on the bundle folder):
//   NODE_PATH="$(npm root -g)" node measure-relation-map.js http://127.0.0.1:8737/<page>.html > out.json
// Optional: CAPTURE_DIR=<folder> CAPTURE_PREFIX=<name> also saves the maintainer's PNG captures.
// SHORT_HEADER_PX (default 87.4, REQ-738's header) gives firstView.frameTopYWithShortHeader: the frame's
// top y at load with the page's own content shifted under a header of that height instead of the one
// the page has now (frameTopY - pageTopY + SHORT_HEADER_PX, where pageTopY is where .page starts).
//
// Light theme, at 390 x 844 and 1440 x 1000. From each starting view (the first view and "Show
// everything") it clicks every shown box and resets to that view between clicks. Title size is
// the title's font size times the map's on-screen scale (REQ-737 D-05); "in view" means the
// title's rectangle lies inside the map frame and inside the screen after the page's own scroll.
// The page contract it reads: #map-frame, #map-world, .box[data-box] with .box-title, .box-rect,
// .rel-line, .rel-word, #btn-main, #btn-all, #map-tools, #map-panel, .box-text, .rel-table rows.
"use strict";
const { chromium } = require("playwright");
const path = require("path");

const pageUrl = process.argv[2];
if (!pageUrl) { console.error("usage: node measure-relation-map.js <page url>"); process.exit(2); }
const captureDir = process.env.CAPTURE_DIR || "";
const capturePrefix = process.env.CAPTURE_PREFIX || "page";
const READABLE_PX = 11.5; // D-05: target 12 px, 11.5 px or more counts
const SHORT_HEADER_PX = parseFloat(process.env.SHORT_HEADER_PX || "87.4");
const SIZES = [{ name: "phone", width: 390, height: 844 }, { name: "desktop", width: 1440, height: 1000 }];

// Runs inside the page: the geometry of one box, its title and the frame.
const PAGE_HELPERS = `
window.__measure = {
  frameRect() {
    const f = document.getElementById("map-frame"), r = f.getBoundingClientRect();
    return { x: r.left + f.clientLeft, y: r.top + f.clientTop, w: f.clientWidth, h: f.clientHeight };
  },
  inside(a, b, tol) { tol = tol || 0.5; return a.x >= b.x - tol && a.y >= b.y - tol && a.x + a.w <= b.x + b.w + tol && a.y + a.h <= b.y + b.h + tol; },
  rectOf(el) { const r = el.getBoundingClientRect(); return { x: r.left, y: r.top, w: r.width, h: r.height }; },
  screen() { return { x: 0, y: 0, w: window.innerWidth, h: window.innerHeight }; },
  box(id) {
    const g = document.getElementById("map-box-" + id), title = g.querySelector(".box-title");
    const scale = title.closest("foreignObject").getScreenCTM().a;
    const titleRect = this.rectOf(title), boxRect = this.rectOf(g.querySelector(".box-rect")), frame = this.frameRect();
    return {
      id, scale: +scale.toFixed(4),
      titlePx: +(parseFloat(getComputedStyle(title).fontSize) * scale).toFixed(2),
      titleInView: this.inside(titleRect, frame) && this.inside(titleRect, this.screen()),
      boxInFrame: this.inside(boxRect, frame),
      hidden: g.classList.contains("is-hidden"),
    };
  },
  shownIds() { return [...document.querySelectorAll(".box[data-box]")].filter(g => !g.classList.contains("is-hidden")).map(g => g.getAttribute("data-box")); },
  // A visible point of the box (inside the frame and the screen) that hits the box itself, or null.
  clickPoint(id) {
    const g = document.getElementById("map-box-" + id), r = this.rectOf(g.querySelector(".box-title")), f = this.frameRect(), s = this.screen();
    const x0 = Math.max(r.x, f.x, s.x), y0 = Math.max(r.y, f.y, s.y), x1 = Math.min(r.x + r.w, f.x + f.w, s.w), y1 = Math.min(r.y + r.h, f.y + f.h, s.h);
    if (x1 - x0 < 6 || y1 - y0 < 6) return null;
    const x = (x0 + x1) / 2, y = (y0 + y1) / 2, hit = document.elementFromPoint(x, y);
    return hit && hit.closest && hit.closest(".box") === g ? { x, y } : null;
  },
  // Pointer events on the box itself, for a shown box the screen does not show (as a script, not a finger).
  syntheticClick(id) {
    const target = document.getElementById("map-box-" + id).querySelector(".box-rect"), r = target.getBoundingClientRect();
    const opts = { bubbles: true, cancelable: true, pointerId: 7, button: 0, isPrimary: true, clientX: r.left + 5, clientY: r.top + 5 };
    target.dispatchEvent(new PointerEvent("pointerdown", opts));
    target.dispatchEvent(new PointerEvent("pointerup", opts));
  },
  // Arrowheads of the drawn lines, as screen rectangles around the line's end point: a marker is
  // 7 line widths long and wide, plus 6 px for a word's halo, to each side of the end point.
  arrowheads() {
    return [...document.querySelectorAll(".rel-line[data-rel]")].filter(l => !l.classList.contains("is-hidden") && getComputedStyle(l).display !== "none").map(l => {
      const ctm = l.getScreenCTM(), p = l.getPointAtLength(l.getTotalLength()), x = ctm.a * p.x + ctm.e, y = ctm.d * p.y + ctm.f;
      const sw = parseFloat(getComputedStyle(l).strokeWidth), len = 7 * sw + 6, half = 3.5 * sw + 6;
      return { x: x - len, y: y - half, w: 2 * len, h: 2 * half };
    });
  },
  // Relation words: in view, on a shown box, on an arrowhead, cut at the frame edge, or broken inside a label.
  words() {
    const f = this.frameRect(), shownBoxes = [...document.querySelectorAll(".box[data-box]:not(.is-hidden) .box-rect")].map(e => this.rectOf(e));
    const heads = this.arrowheads();
    const out = { displayed: 0, inView: 0, onBoxes: 0, onArrowheads: 0, cutAtFrame: 0, broken: 0, brokenIds: [] };
    document.querySelectorAll(".rel-word").forEach(w => {
      if (getComputedStyle(w).display === "none") return;
      const r = this.rectOf(w); if (r.w === 0) return;
      out.displayed++;
      if (this.inside(r, f)) out.inView++; else out.cutAtFrame++;
      if (shownBoxes.some(b => r.x < b.x + b.w - 1 && b.x < r.x + r.w - 1 && r.y < b.y + b.h - 1 && b.y < r.y + r.h - 1)) out.onBoxes++;
      if (heads.some(b => r.x < b.x + b.w && b.x < r.x + r.w && r.y < b.y + b.h && b.y < r.y + r.h)) out.onArrowheads++;
      const n = w.getAttribute("data-rel"), row = document.querySelector("#rel-row-" + n + " strong");
      const label = [...w.querySelectorAll("tspan")].map(t => t.textContent).join(" ");
      if (row && label !== row.textContent) { out.broken++; out.brokenIds.push(n); }
    });
    return out;
  },
  // Box titles never narrower than their longest word: no title overflows its box and none may break a word.
  narrowLabels() {
    let n = 0;
    document.querySelectorAll(".box-title").forEach(t => {
      const cs = getComputedStyle(t);
      if (t.scrollWidth > t.clientWidth + 1 || cs.wordBreak === "break-all" || cs.hyphens === "auto" || cs.overflowWrap === "anywhere") n++;
    });
    return n;
  },
  sideways() { return document.documentElement.scrollWidth - window.innerWidth; },
};`;

async function settle(page) { await page.evaluate(() => new Promise(r => requestAnimationFrame(() => requestAnimationFrame(r)))); }

// Starting view: the page's own button, then the map tool row at the top of the screen.
async function startView(page, button) {
  await page.evaluate(b => { document.getElementById(b).click(); document.getElementById("map-tools").scrollIntoView({ block: "start" }); }, button);
  await settle(page);
}

async function clickBox(page, id) {
  const point = await page.evaluate(i => window.__measure.clickPoint(i), id);
  if (point) await page.mouse.click(point.x, point.y);
  else await page.evaluate(i => window.__measure.syntheticClick(i), id);
  await settle(page);
  return point ? "mouse" : "synthetic";
}

async function selectedState(page) {
  return page.evaluate(() => {
    const m = window.__measure, sel = document.querySelector(".box.is-selected");
    const related = [...document.querySelectorAll(".box.is-related")].map(g => m.box(g.getAttribute("data-box")));
    return { selected: sel && sel.getAttribute("data-box"), clicked: sel && m.box(sel.getAttribute("data-box")), related, words: m.words(), sideways: m.sideways() };
  });
}

// One click from a starting view, then a one-step reach test for every related box whose title is not in view.
async function measureClick(page, button, id, isPhone) {
  await startView(page, button);
  const method = await clickBox(page, id);
  const state = await selectedState(page);
  const row = {
    box: id, method, selectedOk: state.selected === id,
    scale: state.clicked && state.clicked.scale,
    clickedTitlePx: state.clicked && state.clicked.titlePx,
    clickedInView: !!(state.clicked && state.clicked.titleInView),
    clickedBoxInFrame: !!(state.clicked && state.clicked.boxInFrame),
    related: state.related.map(r => ({ box: r.id, titlePx: r.titlePx, titleInView: r.titleInView, boxInFrame: r.boxInFrame })),
    wordsInView: state.words.inView, wordsOnBoxes: state.words.onBoxes, wordsOnArrowheads: state.words.onArrowheads, wordsCut: state.words.cutAtFrame, wordsBroken: state.words.broken,
    sideways: state.sideways, outOfView: [],
  };
  for (const r of state.related.filter(x => !x.titleInView)) {
    // Reset and click again, so each reach starts from the same click view.
    await startView(page, button);
    await clickBox(page, id);
    const step = page.locator('#map-panel button[data-go="' + r.id + '"]').first();
    const named = (await step.count()) > 0;
    let after = null;
    if (named) {
      const label = (await step.textContent()).trim();
      await step.click();
      await settle(page);
      after = await page.evaluate(i => window.__measure.box(i), r.id);
      after.label = label;
    }
    row.outOfView.push({
      box: r.id, stepNamed: named, stepLabel: after && after.label,
      stepWorks: !!(after && after.titleInView && after.titlePx >= READABLE_PX),
      afterTitlePx: after && after.titlePx, afterInView: after && after.titleInView,
    });
  }
  return row;
}

async function measureSize(browser, size) {
  const context = await browser.newContext({ viewport: { width: size.width, height: size.height }, colorScheme: "light", deviceScaleFactor: 2 });
  const page = await context.newPage();
  await page.addInitScript(PAGE_HELPERS);
  await page.goto(pageUrl, { waitUntil: "load" });
  await settle(page);
  const isPhone = size.width < 600;
  const result = { size: size.name, viewport: size.width + "x" + size.height, views: {} };

  // First view, as loaded, then with the map tool row scrolled to the top of the screen.
  const tops = await page.evaluate(() => ({ frame: document.getElementById("map-frame").getBoundingClientRect().top + window.scrollY, page: document.querySelector(".page").getBoundingClientRect().top + window.scrollY }));
  await page.evaluate(() => document.getElementById("map-tools").scrollIntoView({ block: "start" }));
  await settle(page);
  result.firstView = await page.evaluate(() => {
    const m = window.__measure, ids = m.shownIds();
    const boxesInView = ids.map(i => m.box(i)).filter(b => b.titleInView).map(b => b.id);
    const header = document.querySelector('body > aside[aria-label="Second-opinion direction"]');
    return {
      shownBoxes: ids.length, titlesInViewNoPan: boxesInView, titlePx: ids.length ? m.box(ids[0]).titlePx : null,
      words: m.words(), countLine: (document.getElementById("map-count") || {}).textContent || "",
      headerBlockHeight: header ? Math.round(header.getBoundingClientRect().height) : 0,
      sideways: m.sideways(), narrowLabels: m.narrowLabels(),
    };
  });
  result.firstView.frameTopY = Math.round(tops.frame);
  result.firstView.pageTopY = Math.round(tops.page);
  result.firstView.frameTopYWithShortHeader = Math.round(tops.frame - tops.page + SHORT_HEADER_PX);
  if (captureDir) await page.screenshot({ path: path.join(captureDir, capturePrefix + "-" + size.width + "-first-view.png") });

  for (const [viewName, button] of [["firstView", "btn-main"], ["showEverything", "btn-all"]]) {
    await startView(page, button);
    const ids = await page.evaluate(() => window.__measure.shownIds());
    const startWords = await page.evaluate(() => window.__measure.words());
    const clicks = [];
    for (const id of ids) clicks.push(await measureClick(page, button, id, isPhone));
    result.views[viewName] = { shownBoxes: ids.length, startWords, clicks };
  }

  // Tab focus with every box shown and the map zoomed in: each focused box fully in the frame, frame scroll 0.
  await startView(page, "btn-all");
  for (let i = 0; i < 8; i++) await page.click("#btn-zoom-in");
  await page.focus("#map-frame");
  const tab = { reached: 0, inFrame: 0, frameScrollNonZero: 0, scale: null };
  for (let i = 0; i < 60; i++) {
    await page.keyboard.press("Tab");
    await settle(page);
    const t = await page.evaluate(() => {
      const a = document.activeElement, f = document.getElementById("map-frame");
      if (!a || !a.closest || !a.closest(".box")) return null;
      const b = window.__measure.box(a.closest(".box").getAttribute("data-box"));
      return { boxInFrame: b.boxInFrame, scale: b.scale, scroll: f.scrollLeft + f.scrollTop };
    });
    if (!t) break;
    tab.reached++; if (t.boxInFrame) tab.inFrame++; if (t.scroll !== 0) tab.frameScrollNonZero++; tab.scale = t.scale;
  }
  result.tabFocus = tab;

  // The list under the map: one row per relation, rows within the screen on a phone.
  result.list = await page.evaluate(() => {
    const rows = [...document.querySelectorAll(".rel-table tbody tr")], lines = [...document.querySelectorAll(".rel-line[data-rel]")];
    const named = new Set(rows.map(r => { const c = r.querySelectorAll("td"); return [c[0].querySelector("a").textContent, c[1].textContent.trim(), c[2].querySelector("a").textContent].join("|"); }));
    const missing = lines.filter(l => { const t = l.querySelector("title").textContent, n = l.getAttribute("data-rel"), row = document.getElementById("rel-row-" + n); if (!row) return true; const c = row.querySelectorAll("td"); return t !== c[0].querySelector("a").textContent + " " + c[1].textContent.trim() + " " + c[2].querySelector("a").textContent; });
    return {
      relations: lines.length, rows: rows.length, distinctRows: named.size, relationsNotNamed: missing.length,
      rowsWiderThanScreen: rows.filter(r => r.getBoundingClientRect().width > window.innerWidth || r.getBoundingClientRect().right > window.innerWidth + 0.5).length,
      boxesAsText: document.querySelectorAll(".box-text").length, mapBoxes: document.querySelectorAll(".box[data-box]").length,
    };
  });
  result.sidewaysAtEnd = await page.evaluate(() => window.__measure.sideways());

  if (captureDir) await captureSet(page, size, isPhone);
  await context.close();

  // No JavaScript.
  const noJsContext = await browser.newContext({ viewport: { width: size.width, height: size.height }, colorScheme: "light", javaScriptEnabled: false });
  const noJs = await noJsContext.newPage();
  await noJs.goto(pageUrl, { waitUntil: "load" });
  result.noJavaScript = await noJs.evaluate(() => {
    const rectOf = el => { const r = el.getBoundingClientRect(); return { x: r.left, y: r.top, w: r.width, h: r.height }; };
    const tools = document.getElementById("map-tools"), frame = document.getElementById("map-frame");
    const boxes = [...document.querySelectorAll(".box-rect")].map(rectOf);
    let onLine = 0, onBoxes = 0, total = 0;
    document.querySelectorAll(".rel-word").forEach(w => {
      total++;
      const r = rectOf(w), line = document.getElementById("line-r" + w.getAttribute("data-rel")), ctm = line.getScreenCTM(), len = line.getTotalLength();
      let hit = false;
      for (let d = 0; d <= len && !hit; d += 4) { const p = line.getPointAtLength(d), x = ctm.a * p.x + ctm.e, y = ctm.d * p.y + ctm.f; hit = x >= r.x - 4 && x <= r.x + r.w + 4 && y >= r.y - 4 && y <= r.y + r.h + 4; }
      if (hit) onLine++;
      if (boxes.some(b => r.x < b.x + b.w - 1 && b.x < r.x + r.w - 1 && r.y < b.y + b.h - 1 && b.y < r.y + r.h - 1)) onBoxes++;
    });
    return {
      toolsHidden: getComputedStyle(tools).display === "none", words: total, wordsOnTheirLine: onLine, wordsOnBoxes: onBoxes,
      wordsAsSvgText: document.querySelectorAll(".rel-word text").length,
      mapScrolls: frame.scrollHeight > frame.clientHeight || frame.scrollWidth > frame.clientWidth,
      boxesAsText: document.querySelectorAll(".box-text").length, tableRows: document.querySelectorAll(".rel-table tbody tr").length,
      sideways: document.documentElement.scrollWidth - window.innerWidth,
    };
  });
  await noJsContext.close();
  result.summary = summarize(result, isPhone);
  return result;
}

// The maintainer's captures: viewport shots of what a reader sees.
async function captureSet(page, size, isPhone) {
  const shot = async name => page.screenshot({ path: path.join(captureDir, capturePrefix + "-" + size.width + "-" + name + ".png") });
  await startView(page, "btn-main"); await clickBox(page, "O2"); await shot("click-O2");
  if (!isPhone) return;
  await startView(page, "btn-main"); await clickBox(page, "O1"); await shot("click-O1");
  await startView(page, "btn-main"); await clickBox(page, "O2");
  const outOfView = await page.evaluate(() => [...document.querySelectorAll(".box.is-related")].map(g => window.__measure.box(g.getAttribute("data-box"))).filter(b => !b.titleInView).map(b => b.id));
  const step = outOfView.length ? page.locator('#map-panel button[data-go="' + outOfView[0] + '"]').first() : null;
  if (step && (await step.count())) { await step.click(); await settle(page); await shot("click-O2-then-reach-" + outOfView[0]); }
  await startView(page, "btn-all"); await shot("show-everything");
  await clickBox(page, "S6"); await shot("show-everything-click-S6");
}

function summarize(r, isPhone) {
  const clicks = [...r.views.firstView.clicks, ...r.views.showEverything.clicks];
  const related = clicks.flatMap(c => c.related);
  const relatedInView = related.filter(x => x.titleInView);
  const outOfView = clicks.flatMap(c => c.outOfView);
  const min = a => a.length ? Math.min(...a) : null;
  const s = {
    clicks: clicks.length, mouseClicks: clicks.filter(c => c.method === "mouse").length, selectFailed: clicks.filter(c => !c.selectedOk).length,
    scaleRange: [min(clicks.map(c => c.scale)), Math.max(...clicks.map(c => c.scale))],
    smallestClickedTitlePx: min(clicks.map(c => c.clickedTitlePx)),
    clickedTitleNotInView: clicks.filter(c => !c.clickedInView).length,
    relatedShown: related.length, relatedTitlesInView: relatedInView.length,
    smallestRelatedTitleInViewPx: min(relatedInView.map(x => x.titlePx)),
    relatedOutOfView: outOfView.length, outOfViewWithNamedStep: outOfView.filter(x => x.stepNamed).length, outOfViewStepWorks: outOfView.filter(x => x.stepWorks).length,
    clickedOrRelatedBoxOutOfFrame: clicks.filter(c => !c.clickedBoxInFrame).length + related.filter(x => !x.boxInFrame).length,
    maxSideways: Math.max(r.firstView.sideways, r.sidewaysAtEnd, ...clicks.map(c => c.sideways)),
    wordsOnBoxes: r.firstView.words.onBoxes + clicks.reduce((n, c) => n + c.wordsOnBoxes, 0),
    wordsBroken: r.firstView.words.broken + clicks.reduce((n, c) => n + c.wordsBroken, 0),
    wordsOnArrowheads: r.firstView.words.onArrowheads + clicks.reduce((n, c) => n + c.wordsOnArrowheads, 0),
    wordsCutAtFrame: r.firstView.words.cutAtFrame + clicks.reduce((n, c) => n + c.wordsCut, 0),
    frameTopYWithShortHeader: r.firstView.frameTopYWithShortHeader,
    narrowLabels: r.firstView.narrowLabels,
  };
  s.green = isPhone
    ? s.selectFailed === 0 && s.smallestClickedTitlePx >= READABLE_PX && s.clickedTitleNotInView === 0 && (s.smallestRelatedTitleInViewPx === null || s.smallestRelatedTitleInViewPx >= READABLE_PX) &&
      s.outOfViewStepWorks === s.relatedOutOfView && s.maxSideways <= 0 && r.list.relationsNotNamed === 0 && r.list.rows === r.list.relations && r.list.rowsWiderThanScreen === 0 &&
      s.wordsCutAtFrame === 0 && s.wordsOnArrowheads === 0 && s.frameTopYWithShortHeader <= 422
    : s.selectFailed === 0 && s.clickedOrRelatedBoxOutOfFrame === 0 && s.maxSideways <= 0 && r.list.relationsNotNamed === 0;
  return s;
}

(async () => {
  const browser = await chromium.launch();
  const out = { url: pageUrl, measuredAt: new Date().toISOString(), browser: browser.version(), readablePx: READABLE_PX, sizes: {} };
  try { for (const size of SIZES) out.sizes[size.name] = await measureSize(browser, size); }
  finally { await browser.close(); }
  process.stdout.write(JSON.stringify(out, null, 1) + "\n");
})().catch(e => { console.error(e); process.exit(1); });

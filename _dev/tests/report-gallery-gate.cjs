#!/usr/bin/env node
const assert = require('node:assert/strict');
const { createHash } = require('node:crypto');
const fs = require('node:fs/promises');
const path = require('node:path');
const { fileURLToPath, pathToFileURL } = require('node:url');
const { parseArgs } = require('node:util');
const { chromium } = require(process.env.GALLERY_PLAYWRIGHT || 'playwright');

const root = path.resolve(__dirname, '../..');
const { values: options } = parseArgs({ options: {
  gallery: { type: 'string', default: path.join(root, 'docs/report-gallery/index.html') },
  manifest: { type: 'string' },
  out: { type: 'string' },
  selectors: { type: 'string' },
  channel: { type: 'string', default: 'chrome' },
} });

const pageUrl = pathToFileURL(path.resolve(options.gallery)).href;
const repoUrl = (file) => pathToFileURL(path.resolve(root, file)).href;
const hashFile = async (url) =>
  createHash('sha256').update(await fs.readFile(fileURLToPath(url))).digest('hex');
const sameUrl = (actual, expected) =>
  assert.equal(actual, expected, "Expected " + expected + ", reached " + actual);
const normalizeName = (name) => name.trim().toLowerCase().replace(/\s+/g, '-');
const hashPattern = /^[a-f0-9]{64}$/;

async function verifyCapture(entry, report, image, dimensions, viewport) {
  assert(entry, 'No trusted capture entry for ' + report);
  sameUrl(report, repoUrl(entry.report));
  sameUrl(image, repoUrl(entry.image));
  assert.deepEqual(dimensions, [viewport.width, viewport.height]);
  assert.equal(await hashFile(report), entry.reportSha256, 'Report SHA-256 differs');
  assert.equal(await hashFile(image), entry.imageSha256, 'Preview SHA-256 differs');
}

async function verifyLinks(page, probe) {
  const urls = await page.locator('a[href], area[href]').evaluateAll(
    (links) => [...new Set(links.map((link) => link.href))]);
  const errors = [];
  for (const href of urls) {
    try {
      const url = new URL(href);
      const fragment = decodeURIComponent(url.hash.slice(1));
      url.hash = '';
      if (url.protocol === 'file:') {
        assert((await fs.stat(fileURLToPath(url))).isFile(), 'Expected a file');
      } else {
        assert(['http:', 'https:'].includes(url.protocol), 'Unsupported link protocol');
      }
      if (url.protocol !== 'file:' || fragment) {
        const destination = url.href === pageUrl ? page : probe;
        if (destination === probe) {
          const response = await probe.goto(url.href, { waitUntil: 'domcontentloaded' });
          assert(!response || response.ok(), 'HTTP ' + response?.status());
        }
        if (fragment) {
          assert(await destination.evaluate((id) =>
            Boolean(document.getElementById(id) ||
              [...document.getElementsByName(id)].some((node) => node.matches('a'))), fragment),
          'Missing anchor #' + fragment);
        }
      }
    } catch (error) {
      errors.push(href + ': ' + error.message);
    }
  }
  return errors;
}

async function inspectControls(controls) {
  const errors = [];
  for (const control of await controls.all()) {
    const label = await control.evaluate((node) =>
      node.getAttribute('aria-label') || node.textContent.trim() || node.outerHTML.slice(0, 120));
    try {
      await control.scrollIntoViewIfNeeded();
      const problem = await control.evaluate((node) => {
        if (!node.checkVisibility({ checkOpacity: true, checkVisibilityCSS: true })) {
          return 'invisible';
        }
        const box = node.getBoundingClientRect();
        if (box.left < -1 || box.right > innerWidth + 1 || box.height <= 0 || box.width <= 0) {
          return 'outside the horizontal viewport or empty';
        }
        // Sample the interior as well as the center to catch partial clipping or overlap.
        for (const column of [0.15, 0.5, 0.85]) {
          for (const row of [0.15, 0.5, 0.85]) {
            const x = box.left + box.width * column;
            const y = box.top + box.height * row;
            if (!node.contains(document.elementFromPoint(x, y))) return 'clipped or obstructed';
          }
        }
        return null;
      });
      assert(!problem, problem);
      await control.click({ trial: true });
    } catch (error) {
      errors.push(label + ': ' + error.message.split('\n')[0]);
    }
  }
  return errors;
}

async function findControlOverlaps(controls) {
  return controls.evaluateAll((nodes) => {
    const label = (node) => node.getAttribute('aria-label') ||
      node.textContent.trim() || node.localName;
    const visible = nodes.filter((node) =>
      node.checkVisibility({ checkOpacity: true, checkVisibilityCSS: true }));
    const boxes = visible.map((node) => Array.from(node.getClientRects()));
    const errors = [];
    for (const [index, node] of visible.entries()) {
      for (let otherIndex = index + 1; otherIndex < visible.length; otherIndex++) {
        const other = visible[otherIndex];
        if (node.contains(other) || other.contains(node)) continue;
        const overlaps = boxes[index].some((a) => boxes[otherIndex].some((b) =>
          Math.min(a.right, b.right) - Math.max(a.left, b.left) > 1 &&
          Math.min(a.bottom, b.bottom) - Math.max(a.top, b.top) > 1));
        if (overlaps) errors.push(label(node) + ' overlaps ' + label(other));
      }
    }
    return errors;
  });
}

async function verifyFragment(page, fragment, cardSelector, index) {
  await page.goto('about:blank');
  await page.goto(pageUrl + '#' + encodeURIComponent(fragment), { waitUntil: 'load' });
  const card = page.locator(cardSelector).nth(index);
  const result = await card.evaluate(async (node, id) => {
    // Smooth scrolling and late layout can outlive the load event.
    await document.fonts.ready;
    const deadline = performance.now() + 1500;
    do {
      const target = document.getElementById(id);
      if (!target || (target !== node && !node.contains(target))) return false;
      const rect = node.getBoundingClientRect();
      const targetRect = target.getBoundingClientRect();
      const scroller = document.scrollingElement;
      const maxScroll = scroller.scrollHeight - scroller.clientHeight;
      const visible = rect.bottom > 0 && rect.top < innerHeight &&
        targetRect.bottom > 0 && targetRect.top < innerHeight;
      const aligned = Math.abs(targetRect.top -
        (parseFloat(getComputedStyle(target).scrollMarginTop) || 0) -
        (parseFloat(getComputedStyle(scroller).scrollPaddingTop) || 0)) <= 2;
      if (visible &&
          (aligned || scrollY >= maxScroll - 2 || scrollY <= 2 && targetRect.top <= 2)) {
        return true;
      }
      await new Promise(requestAnimationFrame);
    } while (performance.now() < deadline);
    return false;
  }, fragment);
  assert(result, 'Bare fragment #' + fragment + ' did not bring its card into view');
}

async function verifyPreview(page, preview, expected) {
  let popup;
  const rememberPopup = (opened) => { popup = opened; };
  page.on('popup', rememberPopup);
  try {
    await preview.click();
    const deadline = Date.now() + 2000;
    while ((popup || page).url() !== expected && Date.now() < deadline) {
      await new Promise((resolve) => setTimeout(resolve, 50));
    }
    const destination = popup || page;
    sameUrl(destination.url(), expected);
    await destination.waitForLoadState('domcontentloaded');
  } finally {
    page.off('popup', rememberPopup);
    await popup?.close();
  }
}

async function runView(browser, width, theme, selectors, manifest, output) {
  const context = await browser.newContext({
    viewport: { width, height: 900 }, colorScheme: theme, reducedMotion: 'reduce',
    permissions: ['clipboard-read', 'clipboard-write'],
  });
  const failures = Object.fromEntries(Array.from({ length: 7 }, (_, index) =>
    ['G' + (index + 1), []]));
  const record = async (gate, label, check) => {
    try {
      await check();
    } catch (error) {
      failures[gate].push(label + ': ' + error.message.split('\n')[0]);
    }
  };
  const page = await context.newPage();
  const probe = await context.newPage();
  page.setDefaultTimeout(2200);
  probe.setDefaultNavigationTimeout(8000);

  try {
    await page.goto(pageUrl, { waitUntil: 'load' });
    const cards = page.locator(selectors.card);
    const count = await cards.count();
    await record('G1', 'Card count', () => assert.equal(count, 10));
    await record('G1', 'Visible card count', async () =>
      assert.equal(await cards.filter({ visible: true }).count(), 10));
    await page.locator('img').evaluateAll((images) => {
      for (const image of images) image.loading = 'eager';
    });
    for (const image of await page.locator('img').all()) {
      const source = await image.evaluate((node) => node.currentSrc || node.src);
      await record('G6', source || 'Image without source', async () => {
        await image.evaluate((node) => Promise.race([
          node.decode(),
          new Promise((_, reject) => setTimeout(() => reject(Error('Image decode timed out')), 5000)),
        ]));
      });
    }
    failures.G6.push(...await verifyLinks(page, probe));

    await page.screenshot({ path: path.join(output, width + '-' + theme + '.png'),
      fullPage: true, animations: 'disabled' });

    const names = new Set();
    const reports = new Set();
    const cardsToExercise = [];
    for (let index = 0; index < count; index++) {
      const card = cards.nth(index);
      const nameControl = card.locator(selectors.name);
      const preview = card.locator(selectors.preview);
      const image = preview.locator(selectors.image);
      const label = 'Card ' + (index + 1);
      let name;
      let report;
      await record('G3', label, async () => {
        name = (await nameControl.innerText()).trim();
        assert(/^[a-z]+(?:-[a-z]+){1,2}$/.test(name) && name.length <= 24,
          'Invalid name ' + JSON.stringify(name));
        assert(!names.has(name), 'Duplicate name ' + name);
        names.add(name);
      });
      await record('G2', name || label, async () => {
        report = await (selectors.report ? card.locator(selectors.report) : preview)
          .evaluate((node) => node.href);
        const capture = manifest?.captures.find((entry) => repoUrl(entry.report) === report);
        const rendered = await image.evaluate((node) => ({
          url: node.currentSrc, dimensions: [node.naturalWidth, node.naturalHeight],
        }));
        await verifyCapture(capture, report, rendered.url, rendered.dimensions, manifest?.viewport);
        assert(!reports.has(report), 'Duplicate report ' + report);
        reports.add(report);
      });
      await record('G7', name || label, async () => {
        assert(await card.isVisible(), 'Card hidden');
        assert(await image.isVisible(), 'Preview image hidden');
        const problems = await inspectControls(nameControl.or(preview));
        assert(!problems.length, problems.join(' | '));
      });
      cardsToExercise.push({ index, name, report });
    }
    // Optional disclosures may contain additional interactive controls.
    for (const summary of await page.locator('details:not([open]) > summary:first-child').all()) {
      if (await summary.isVisible()) await summary.click().catch((error) =>
        failures.G7.push(error.message.split('\n')[0]));
    }
    const controls = page.locator(selectors.controls).filter({ visible: true });
    failures.G7.push(...await findControlOverlaps(controls));
    failures.G7.push(...await inspectControls(controls));
    await record('G7', 'Horizontal overflow', async () => assert(
      await page.evaluate(() => document.scrollingElement.scrollWidth <= innerWidth + 1),
      'Document is wider than viewport'));

    for (const { index, name, report } of cardsToExercise) {
      const label = name || 'Card ' + (index + 1);
      await page.bringToFront();
      await page.goto(pageUrl);
      const card = page.locator(selectors.card).nth(index);
      let copied;
      await record('G4', label, async () => {
        const sentinel = 'gallery-gate-clipboard-' + index;
        await page.evaluate((text) => navigator.clipboard.writeText(text), sentinel);
        await card.locator(selectors.name).click();
        const deadline = Date.now() + 1000;
        do {
          copied = await page.evaluate(() => navigator.clipboard.readText());
          if (copied !== sentinel) break;
          await new Promise((resolve) => setTimeout(resolve, 40));
        } while (Date.now() < deadline);
        assert.notEqual(copied, sentinel, 'Clipboard stayed at sentinel');
        assert.equal(copied, name, 'Clipboard differs from visible name');
        assert(/^[a-z]+(?:-[a-z]+){1,2}$/.test(copied) && copied.length <= 24,
          'Copied text is not a bare name: ' + JSON.stringify(copied));
      });
      const fragment = copied && !copied.startsWith('gallery-gate-clipboard-') ? copied : name;
      await record('G5', label, async () => {
        assert(fragment, 'No name available');
        await verifyFragment(page, fragment, selectors.card, index);
      });
      if (fragment && normalizeName(fragment) !== fragment) {
        await record('G5', label + ' (bare slug)', () =>
          verifyFragment(page, normalizeName(fragment), selectors.card, index));
      }
      await record('G5', label + ' (preview)', async () => {
        assert(report, 'No expected report URL');
        await verifyPreview(page, card.locator(selectors.preview), report);
      });
    }
  } finally {
    await context.close();
  }
  return failures;
}

async function main() {
  assert(options.out, 'Supply --out DIRECTORY');
  const output = path.resolve(options.out);
  const selectors = {
    card: '[data-gallery-card]',
    name: '[data-gallery-name]',
    preview: '[data-gallery-preview]',
    image: 'img',
    controls: 'a[href], button, summary, input:not([type=hidden]), select, textarea, [role=button], [role=link]',
    ...options.selectors && JSON.parse(await fs.readFile(options.selectors, 'utf8')),
  };
  let manifest;
  if (options.manifest) {
    const raw = JSON.parse(await fs.readFile(options.manifest, 'utf8'));
    manifest = raw.cards ? {
      captures: raw.cards,
      viewport: { width: raw.width, height: raw.height },
      scrollY: raw.scrollY,
    } : raw;
    assert.equal(manifest.captures.length, 10);
    assert.equal(new Set(manifest.captures.map((entry) => entry.report)).size, 10);
    assert.equal(new Set(manifest.captures.map((entry) => entry.image)).size, 10);
    assert(manifest.scrollY === 0 && Number.isInteger(manifest.viewport.width) &&
      manifest.viewport.width > 0 && Number.isInteger(manifest.viewport.height) &&
      manifest.viewport.height > 0, 'Expected fixed top-viewport capture settings');
    for (const entry of manifest.captures) {
      assert(entry.report.startsWith('ai-reports/') && entry.report.endsWith('/index.html') &&
        path.posix.normalize(entry.report) === entry.report,
      'Expected canonical ai-reports/.../index.html source');
      assert(hashPattern.test(entry.reportSha256) && hashPattern.test(entry.imageSha256),
        'Expected both SHA-256 hashes');
    }
  }

  await fs.mkdir(output, { recursive: true });
  const browser = await chromium.launch({ channel: options.channel });
  const results = {};
  try {
    for (const width of [1440, 390]) {
      for (const theme of ['light', 'dark']) {
        const view = width + '-' + theme;
        const failures = await runView(browser, width, theme, selectors, manifest, output);
        results[view] = failures;
        for (const [gate, errors] of Object.entries(failures)) {
          process.stdout.write(view + ' ' + gate + ' ' + (errors.length ? 'FAIL' : 'PASS') + '\n');
          for (const error of errors) process.stdout.write('  ' + error + '\n');
          if (errors.length) process.exitCode = 1;
        }
      }
    }
  } finally {
    await browser.close();
    await fs.writeFile(path.join(output, 'gates.json'), JSON.stringify(results, null, 2) + '\n');
  }
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 2;
});

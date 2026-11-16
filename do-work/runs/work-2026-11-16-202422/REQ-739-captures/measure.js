async (page) => {
  // REQ-739 measurement: header side margin, page-rule precedence, map top, sideways scroll, fold alignment, kit checks.
  const base = 'http://127.0.0.1:61370/';
  const pages = [
    ['P1 kit new', 'kit/p1.html'], ['P1 kit before(0.324.2)', 'kit/b1-p1.html'],
    ['P2 v1c01 new', 'tstack/index-v1c01-codex-gpt-6.1-terra-custom-matrix-to-implementation.html'], ['P2 v1c01 before(0.324.2)', 'tstack/b1-v1c01.html'],
    ['P3 plain new', 'plain/p3.html'], ['P3 plain before(0.324.2)', 'plain/b1-p3.html'],
    ['P4 v1c03 new', 'tstack/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html'], ['P4 v1c03 published(0.324.0)', 'tstack/b0-v1c03.html'], ['P4 v1c03 before(0.324.2)', 'tstack/b1-v1c03.html'],
    ['P5 v1c04 new', 'tstack/index-v1c04-claude-claude-opus-5-5-relation-map-phone-first-map.html'], ['P5 v1c04 published(=0.324.2)', 'tstack/b1-v1c04.html'],
  ];
  await page.emulateMedia({ colorScheme: 'light' });
  const out = [];
  for (const [width, height] of [[390, 844], [1440, 1000]]) {
    await page.setViewportSize({ width, height });
    for (const [name, path] of pages) {
      await page.goto(base + path, { waitUntil: 'load' });
      await page.waitForTimeout(150);
      const m = await page.evaluate((isKit) => {
        const r = (n) => Math.round(n * 10) / 10;
        const aside = document.querySelector('body>aside[aria-label="Second-opinion direction"]');
        const cs = getComputedStyle(aside);
        const walker = document.createTreeWalker(aside, NodeFilter.SHOW_TEXT, { acceptNode: (n) => n.textContent.trim() ? 1 : 3 });
        const first = walker.nextNode();
        const range = document.createRange(); range.selectNodeContents(first);
        const textLeft = range.getClientRects()[0].left;
        const box = aside.getBoundingClientRect();
        const vw = document.documentElement.clientWidth;
        const contentRight = box.right - parseFloat(cs.borderRightWidth) - parseFloat(cs.paddingRight);
        const lines = (el) => { const g = document.createRange(); g.selectNodeContents(el); return new Set([...g.getClientRects()].map((x) => Math.round(x.top))).size; };
        const iteration = aside.querySelector('p');
        const details = aside.querySelector('details');
        const summary = aside.querySelector('summary');
        const frame = document.getElementById('map-frame');
        const scroll = () => [document.documentElement.scrollWidth, document.documentElement.clientWidth];
        const closed = scroll();
        details.open = true; const opened = scroll(); details.open = false;
        const res = {
          textLeft: r(textLeft), rightGap: r(vw - contentRight),
          headerHeight: r(box.height + parseFloat(cs.marginTop) + parseFloat(cs.marginBottom)),
          margin: cs.margin, padding: cs.padding,
          mapTop: frame ? r(frame.getBoundingClientRect().top + scrollY) : null,
          scrollClosed: closed.join('/'), scrollOpen: opened.join('/'),
          iterationLeft: r(iteration.getBoundingClientRect().left), summaryLeft: r(summary.getBoundingClientRect().left),
          iterationLines: lines(iteration), summaryLines: lines(summary),
        };
        if (isKit) {
          const ds = getComputedStyle(details), ss = getComputedStyle(summary);
          res.kitCard = { detailsBorder: ds.borderTopWidth, detailsRadius: ds.borderTopLeftRadius, summaryWeight: ss.fontWeight, summarySize: ss.fontSize, bodySize: getComputedStyle(document.body).fontSize };
          const button = document.getElementById('open-all');
          const section = button.closest('section');
          button.click();
          const kitDetails = [...section.querySelectorAll('details')];
          res.kitCard.openAllOpensKit = kitDetails.length > 0 && kitDetails.every((d) => d.open);
          res.kitCard.headerFoldStaysClosed = !details.open;
        }
        return res;
      }, name.startsWith('P1'));
      out.push(`${width} | ${name} | ` + JSON.stringify(m));
    }
  }
  return out.join('\n');
}

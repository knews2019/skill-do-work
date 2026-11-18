# Original gallery

Captured from main at e26f4962 before REQ-749 changed the gallery. Each file includes all ten cards and the sections below them. The browser verified that all ten card images decoded before capture.

From the repository root, run `python3 -m http.server 17499 --bind 127.0.0.1` in a separate terminal. Run the following against the original revision to reproduce the desktop light image. For the other views, use 390 × 844 for phone and set `colorScheme` to `dark` for dark mode.

```sh
playwright-cli -s=req749 open http://127.0.0.1:17499/docs/report-gallery/
playwright-cli -s=req749 resize 1440 900
playwright-cli -s=req749 run-code "async page => { await page.emulateMedia({colorScheme:'light'}); await page.locator('img[src]').evaluateAll(async images => { if(images.length!==10) throw Error('Expected ten previews'); for(const img of images) img.loading='eager'; await Promise.all(images.map(img=>img.decode())); }); }"
playwright-cli -s=req749 screenshot --full-page --filename=original-desktop-light.png
playwright-cli -s=req749 close
```

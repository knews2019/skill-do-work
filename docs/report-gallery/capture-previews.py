"""Regenerate the gallery previews from their source report pages."""
from hashlib import sha256
from tempfile import TemporaryDirectory
from functools import partial
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer
import json
from pathlib import Path
import subprocess
from threading import Thread
from uuid import uuid4

gallery = Path(__file__).resolve().parent
root = gallery.parents[1]
manifest = json.loads((gallery / "preview-manifest.json").read_text())
session = f"-s=report-previews-{uuid4().hex}"


def playwright(*args):
    result = subprocess.run(
        ["playwright-cli", session, *args], cwd=root, text=True,
        stdout=subprocess.PIPE, stderr=subprocess.STDOUT, check=True,
    )
    if "### Error" in result.stdout:
        raise RuntimeError(result.stdout)


if __name__ == "__main__":
    handler = partial(SimpleHTTPRequestHandler, directory=root)
    with TemporaryDirectory(dir=gallery) as staging, ThreadingHTTPServer(("127.0.0.1", 0), handler) as server:
        Thread(target=server.serve_forever, daemon=True).start()
        origin = f"http://127.0.0.1:{server.server_port}"
        staged_cards = [
            {**card, "stagedImage": str(Path(staging) / Path(card["image"]).name)}
            for card in manifest["cards"]
        ]
        source_hashes = {
            card["report"]: sha256((root / card["report"]).read_bytes()).hexdigest()
            for card in staged_cards
        }
        cards = json.dumps(staged_cards)
        try:
            playwright("open", "about:blank")
            playwright("resize", str(manifest["width"]), str(manifest["height"]))
            playwright("run-code", f"""async page => {{
                await page.emulateMedia({{colorScheme: 'light', reducedMotion: 'reduce'}});
                for (const card of {cards}) {{
                    const response = await page.goto({json.dumps(origin)} + '/' + card.report);
                    if (!response.ok()) throw Error(card.report + ': ' + response.status());
                    await page.evaluate(async () => {{
                        await document.fonts.ready;
                        window.scrollTo(0, 0);
                        await new Promise(resolve => requestAnimationFrame(() => requestAnimationFrame(resolve)));
                    }});
                    if (await page.evaluate(() => window.scrollY) !== 0) throw Error(card.report + ': unexpected scroll');
                    await page.screenshot({{path: card.stagedImage, scale: 'css', animations: 'disabled'}});
                }}
            }}""")
        finally:
            playwright("close")

        for card in staged_cards:
            report_hash = sha256((root / card["report"]).read_bytes()).hexdigest()
            if report_hash != source_hashes[card["report"]]:
                raise RuntimeError(f"Report changed during capture: {card['report']}")
            image = Path(card["stagedImage"]).read_bytes()
            if not image:
                raise RuntimeError(f"Empty preview: {card['image']}")
            card["reportSha256"] = report_hash
            card["imageSha256"] = sha256(image).hexdigest()
        for card in staged_cards:
            Path(card.pop("stagedImage")).replace(root / card["image"])
        manifest["cards"] = staged_cards
        manifest["scrollY"] = 0
        staged_manifest = Path(staging) / "preview-manifest.json"
        staged_manifest.write_text(json.dumps(manifest, indent=2) + "\n")
        staged_manifest.replace(gallery / staged_manifest.name)

"""Regenerate the gallery previews from their source report pages."""
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
    with ThreadingHTTPServer(("127.0.0.1", 0), handler) as server:
        Thread(target=server.serve_forever, daemon=True).start()
        origin = f"http://127.0.0.1:{server.server_port}"
        cards = json.dumps(manifest["cards"])
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
                    await page.screenshot({{path: card.image, scale: 'css', animations: 'disabled'}});
                }}
            }}""")
        finally:
            playwright("close")

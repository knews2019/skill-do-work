#!/usr/bin/env python3
"""Generate a local REQ-752 report. See INPUT.md for the JSON contract."""

import argparse
import json
import os
from html import escape
from pathlib import Path
from urllib.parse import quote


def html(value):
    return escape(str(value), quote=True)


def read_json(path):
    return json.loads(path.read_text(encoding="utf-8"))


def local_file(directory, value):
    path = directory / value
    if not path.is_file():
        raise ValueError(f"Local file not found: {path}")
    return path.absolute()


def require_count(items, expected, label):
    if len(items) != expected:
        raise ValueError(f"{label} needs exactly {expected} entries")


def table(headings, rows):
    header = "".join(f'<th scope="col">{html(label)}</th>' for label in headings)
    body = "".join(
        '<tr><th scope="row">' + row[0] + '</th>'
        + "".join(f"<td>{cell}</td>" for cell in row[1:]) + "</tr>"
        for row in rows
    )
    return f"<table><thead><tr>{header}</tr></thead><tbody>{body}</tbody></table>"


def render(data, source_dir, output_dir):
    def href(value, directory=source_dir):
        path = local_file(directory, value)
        return quote(os.path.relpath(path, output_dir))

    manifest = local_file(source_dir, data["manifest"])
    reports = read_json(manifest)
    require_count(reports, 10, "Manifest")
    report_rows = [
        (html(item["name"]), html(item["kind"]),
         f'<a href="{href(item["path"], manifest.parent)}">{html(item["path"])}</a>')
        for item in reports
    ]

    comparison = data["comparison"]
    baseline, final = comparison["baseline"], comparison["final"]
    counts = data["counts"]
    if any(type(counts[key]) is not int or counts[key] < 0
           for key in ("rounds", "kept", "reverted")):
        raise ValueError("Candidate counts must be nonnegative integers")
    if counts["rounds"] != counts["kept"] + counts["reverted"]:
        raise ValueError("Candidate rounds must equal kept plus reverted")

    galleries = []
    for stage in ("original", "baseline", "final"):
        gallery = data["galleries"][stage]
        figures = []
        for width in (1440, 390):
            for theme in ("light", "dark"):
                variant = f"{width}-{theme}"
                capture = gallery["captures"][variant]
                label = f"{stage.title()} gallery · {width} px · {theme}"
                image_url = href(capture["path"])
                figures.append(
                    f'<figure><a href="{image_url}"><img loading="lazy" src="{image_url}"'
                    f' alt="{html(label)}"></a><figcaption>{html(label)}. '
                    f'{html(capture["provenance"])}</figcaption></figure>'
                )
        galleries.append(
            f'<section><h3>{stage.title()}</h3><div class="gallery">'
            + "".join(figures) + "</div></section>"
        )

    rubric_labels = ("Previews", "Names", "Choice", "Range", "Attention")
    rubric_rows = [
        (f"R{index} · {label}", html(baseline["rubric"][f"R{index}"]),
         html(final["rubric"][f"R{index}"]))
        for index, label in enumerate(rubric_labels, 1)
    ]

    choices = []
    for label, endpoint in (("Baseline", baseline), ("Final", final)):
        require_count(endpoint["choices"], 4, f"{label} first choices")
        rows = []
        for choice in endpoint["choices"]:
            match = choice["match"]
            if type(match) is not bool:
                raise ValueError("First-choice match must be a JSON boolean")
            rows.append((html(choice["need"]), html(choice["first_choice"]),
                         "Match" if match else "Mismatch"))
        choices.append(f"<h3>{label}</h3><p>{html(endpoint['weakest_point'])}</p>"
                       + table(("Need", "First choice", "Result"), rows))

    evidence = [
        ("Fresh blind evaluation", comparison["evidence"]),
        ("Climb log and historical round scores", data["climb_log"]),
        *((item["label"], item["path"]) for item in data["evidence"]),
    ]
    timing = data.get("timing")
    timing_copy = "No timed test with new readers was supplied."
    if timing:
        timing_copy = (f'{html(timing["summary"])} '
                       f'<a href="{href(timing["evidence"])}">Timing evidence</a>.')

    decisions = data["open_decisions"]
    decision_items = [
        f"<li><strong>D{number} · {html(item['question'])}</strong>"
        f"<dl><dt>Value</dt><dd>{html(item['value'])}</dd>"
        f"<dt>Risk</dt><dd>{html(item['risk'])}</dd>"
        f"<dt>Recommendation</dt><dd>{html(item['recommendation'])}</dd>"
        f"<dt>Default</dt><dd>{html(item['default'])}</dd></dl></li>"
        for number, item in enumerate(decisions, 1)
    ]
    decision_copy = ("<ul>" + "".join(decision_items) + "</ul>"
                     if decision_items else "<p>None remain.</p>")
    evidence_links = "".join(
        f'<li><a href="{href(path)}">{html(label)}</a></li>' for label, path in evidence
    )

    return f"""<!doctype html>
<html lang="en">
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Gallery climb final comparison</title>
<style>
:root {{ color-scheme: light dark; }}
body {{ font: 1rem/1.5 system-ui, sans-serif; margin: auto; max-width: 74rem;
        padding: 1rem; overflow-wrap: anywhere; }}
main > section {{ margin-block: 2.5rem; }}
table {{ width: 100%; border-collapse: collapse; table-layout: fixed; }}
th, td {{ border-bottom: 1px solid GrayText; text-align: left;
          vertical-align: top; padding: .45rem; }}
.gallery {{ display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 1rem; }}
figure {{ margin: 0; min-width: 0; }}
img {{ display: block; width: 100%; height: auto; }}
figcaption {{ margin-top: .35rem; }}
dt {{ font-weight: 600; }}
@media (max-width: 42rem) {{ .gallery {{ grid-template-columns: 1fr; }} }}
</style>
<main>
<h1>Gallery climb final comparison</h1>
<p><strong>{html(data['verdict'])} · Fresh blind comparison: baseline {html(baseline['score'])},
final {html(final['score'])} · {counts['rounds']} candidate rounds,
{counts['kept']} kept, {counts['reverted']} reverted</strong></p>
<p>One fresh blind evaluator ({html(comparison['judge'])}) scored both endpoints.
The climb log records historical scores. Candidate counts exclude Round 0 and the final comparison.</p>
<p><strong>Stopping rule.</strong> {html(data['stopping_rule'])}</p>
<p><strong>Ten-second selection goal.</strong> {timing_copy}</p>
<section><h2>Gallery screenshots</h2>{''.join(galleries)}</section>
<section><h2>Reports in the manifest</h2>{table(('Name', 'Kind', 'Path'), report_rows)}</section>
<section><h2>Fresh blind rubric</h2>{table(('Criterion', 'Baseline', 'Final'), rubric_rows)}</section>
<section><h2>First choices and checked matches</h2>{''.join(choices)}</section>
<section><h2>Evidence</h2><ul>{evidence_links}</ul></section>
<section><h2>Open decisions</h2>{decision_copy}</section>
</main>
</html>
"""


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("input", type=Path, help="report data JSON")
    parser.add_argument("output", type=Path, help="new HTML file in an existing directory")
    args = parser.parse_args()
    try:
        source = args.input.resolve()
        output = args.output.absolute()
        content = render(read_json(source), source.parent, output.parent.resolve())
        with output.open("x", encoding="utf-8") as handle:
            handle.write(content)
    except (OSError, ValueError, KeyError, TypeError) as error:
        parser.exit(1, f"Cannot generate report: {error}\n")


if __name__ == "__main__":
    main()

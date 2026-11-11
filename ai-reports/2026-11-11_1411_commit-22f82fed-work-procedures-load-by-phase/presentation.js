// ai-report presentation kit: reader-driven behavior. Copy beside index.html with presentation.css.
// Every report reads fully without this file. It only adds focus for glossary words, stepping,
// "open all details" and deck mode. Nothing moves until the reader clicks or presses a key.

// Glossary words: <span class="term" data-def="plain definition">word</span>. CSS shows the definition
// on hover or focus. Focus is what makes a tap work on a touch screen.
for (const term of document.querySelectorAll(".term[data-def]")) {
  if (!term.hasAttribute("tabindex")) term.tabIndex = 0;
  term.setAttribute("aria-description", term.dataset.def);
}

// Reader-driven stepping. Markup:
//   <section data-stepper> ... [data-step] children in order ...
//     elements with data-at="1 3" light up at those steps (1-based)
//     [data-step-prev] [data-step-next] buttons, [data-step-label] text
// Arrow keys step while focus is inside the stepper. Clicking a step selects it.
for (const stepper of document.querySelectorAll("[data-stepper]")) {
  const steps = [...stepper.querySelectorAll("[data-step]")];
  if (steps.length === 0) continue;
  const label = stepper.querySelector("[data-step-label]");
  let active = 1;
  const show = (index) => {
    active = Math.min(Math.max(index, 1), steps.length);
    steps.forEach((step, position) => step.classList.toggle("is-active", position + 1 === active));
    for (const part of stepper.querySelectorAll("[data-at]")) {
      part.classList.toggle("is-lit", part.dataset.at.split(/\s+/).includes(String(active)));
    }
    if (label) label.textContent = `Step ${active} of ${steps.length}`;
  };
  // The badge pop in a terminal replay plays only after the reader has used this stepper.
  const touch = () => { stepper.dataset.touched = ""; };
  stepper.classList.add("stepper-on");
  stepper.querySelector("[data-step-prev]")?.addEventListener("click", () => { touch(); show(active - 1); });
  stepper.querySelector("[data-step-next]")?.addEventListener("click", () => { touch(); show(active + 1); });
  steps.forEach((step, position) => step.addEventListener("click", () => { touch(); show(position + 1); }));
  stepper.addEventListener("keydown", (event) => {
    if (event.key === "ArrowRight" || event.key === "ArrowDown") { touch(); show(active + 1); event.preventDefault(); }
    if (event.key === "ArrowLeft" || event.key === "ArrowUp") { touch(); show(active - 1); event.preventDefault(); }
  });
  if (!stepper.hasAttribute("tabindex")) stepper.tabIndex = -1;
  show(1);
}

// "Open all details" button for depth 3: <button id="open-all"> inside the section that holds the <details>.
const openAll = document.getElementById("open-all");
openAll?.addEventListener("click", () => {
  const details = [...openAll.closest("section").querySelectorAll("details")];
  const open = !details.every((item) => item.open);
  details.forEach((item) => { item.open = open; });
  openAll.textContent = open ? "Fold all details" : "Open all details";
});

// Deck mode: each .page > .band is one slide. D or the #deck-toggle button switches it, Escape leaves it.
const deckToggle = document.getElementById("deck-toggle");
const deckHint = document.getElementById("deck-hint");
const deckSlides = [...document.querySelectorAll(".page > .band")];
const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)");
// While a smooth scroll is still moving, count from the slide it is heading to, so quick presses add up.
let headingTo = null;
window.addEventListener("scrollend", () => { headingTo = null; });

function currentSlideIndex() {
  const middle = window.innerHeight / 2;
  let index = 0;
  deckSlides.forEach((slide, position) => { if (slide.getBoundingClientRect().top <= middle) index = position; });
  return index;
}
function deckOn() { return document.documentElement.classList.contains("deck"); }
function showHint() {
  if (deckHint) deckHint.textContent = `Slide ${currentSlideIndex() + 1} of ${deckSlides.length}. Arrow keys move. Press D to leave.`;
}
function setDeck(on) {
  if (deckSlides.length === 0) return;
  const anchor = deckSlides[currentSlideIndex()];
  document.documentElement.classList.toggle("deck", on);
  deckToggle?.setAttribute("aria-pressed", String(on));
  anchor.scrollIntoView({ block: "start" });
  if (on) showHint();
}
function moveSlides(offset) {
  const from = headingTo ?? currentSlideIndex();
  headingTo = Math.min(Math.max(from + offset, 0), deckSlides.length - 1);
  deckSlides[headingTo].scrollIntoView({ behavior: reduceMotion.matches ? "auto" : "smooth", block: "start" });
}

deckToggle?.addEventListener("click", () => setDeck(!deckOn()));
window.addEventListener("scroll", () => { if (deckOn()) showHint(); }, { passive: true });
document.addEventListener("keydown", (event) => {
  // A stepper that handled the key has already called preventDefault.
  if (event.altKey || event.ctrlKey || event.metaKey || event.defaultPrevented) return;
  if (event.target instanceof Element && event.target.closest("input, textarea, select, [contenteditable]")) return;
  if (event.key === "d" || event.key === "D") { setDeck(!deckOn()); return; }
  if (!deckOn()) return;
  if (event.key === "Escape") { setDeck(false); return; }
  if (event.target instanceof Element && event.target.closest("button, summary")) return;
  const forward = ["ArrowRight", "ArrowDown", "PageDown"].includes(event.key) || (event.key === " " && !event.shiftKey);
  const backward = ["ArrowLeft", "ArrowUp", "PageUp"].includes(event.key) || (event.key === " " && event.shiftKey);
  if (!forward && !backward) return;
  event.preventDefault();
  moveSlides(forward ? 1 : -1);
});

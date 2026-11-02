// Presentation-style gallery: shared behavior. Every page works without this file; it only adds
// glossary definitions and reader-driven stepping.

// One definition per word, written for a reader who has never seen the code.
const GLOSSARY = {
  backend: "The other AI tool the launcher runs to write a second opinion: codex, gemini or agy.",
  launcher: "The do-work-cli command that starts the backend, collects its output and checks the result.",
  cap: "The most backend output the launcher keeps in memory: 16 MB (16,777,216 bytes). Anything past it should be cut and reported.",
  buffer: "A piece of memory that collects bytes as they arrive. Here: secondOpinionBuffer, one for stdout and one for stderr.",
  embedding: "A Go feature: put a type inside a struct without giving it a field name. The struct then gets all of that type's methods for free.",
  promoted: "A method the outer struct gets from an embedded type. You did not write it, but the struct has it.",
  readfrom: "A method on bytes.Buffer that keeps reading a stream into memory until the stream ends. It does not know about our 16 MB limit.",
  write: "Our own method on the buffer. It checks the 16 MB limit, cuts extra bytes and marks the buffer as truncated.",
  iocopy: "Go's standard function that moves bytes from a reader to a writer. If the writer has ReadFrom, io.Copy hands it the whole stream and skips Write.",
  pipe: "A channel the operating system gives two programs: the backend writes into one end, the launcher reads from the other.",
  truncated: "A flag set when output was cut at the limit. The launcher turns it into the error \"backend output exceeded capture limit\".",
  namedfield: "A struct field with a name, like captured bytes.Buffer. Its methods stay on the field, so the struct does not get them.",
  pin: "A test that fails if this exact bug comes back. It names the failure it guards against.",
};

function fillGlossary() {
  for (const term of document.querySelectorAll(".term[data-term]")) {
    const definition = GLOSSARY[term.dataset.term];
    if (!definition) continue;
    if (!term.dataset.def) term.dataset.def = definition;
    if (!term.hasAttribute("tabindex")) term.tabIndex = 0;
    term.setAttribute("aria-description", definition);
  }
  // Pages may ask for a visible glossary list: <dl data-glossary-list="readfrom write iocopy">.
  for (const list of document.querySelectorAll("[data-glossary-list]")) {
    const keys = list.dataset.glossaryList.trim().split(/\s+/);
    for (const key of keys) {
      const label = document.querySelector(`.term[data-term="${key}"]`)?.textContent ?? key;
      const term = document.createElement("dt");
      term.textContent = label;
      const definition = document.createElement("dd");
      definition.textContent = GLOSSARY[key] ?? "";
      list.append(term, definition);
    }
  }
}

// Reader-driven stepping. Markup:
//   <section data-stepper> ... [data-step] children in order ...
//     elements with data-at="1 3" light up at those steps (1-based)
//     [data-step-prev] [data-step-next] buttons, [data-step-label] text
// Arrow keys step while focus is inside the stepper. Clicking a step selects it.
function setupSteppers() {
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
      stepper.dataset.active = String(active);
    };
    stepper.classList.add("stepper-on");
    stepper.querySelector("[data-step-prev]")?.addEventListener("click", () => show(active - 1));
    stepper.querySelector("[data-step-next]")?.addEventListener("click", () => show(active + 1));
    steps.forEach((step, position) => step.addEventListener("click", () => show(position + 1)));
    stepper.addEventListener("keydown", (event) => {
      if (event.key === "ArrowRight" || event.key === "ArrowDown") { show(active + 1); event.preventDefault(); }
      if (event.key === "ArrowLeft" || event.key === "ArrowUp") { show(active - 1); event.preventDefault(); }
    });
    if (!stepper.hasAttribute("tabindex")) stepper.tabIndex = -1;
    show(1);
  }
}

fillGlossary();
setupSteppers();

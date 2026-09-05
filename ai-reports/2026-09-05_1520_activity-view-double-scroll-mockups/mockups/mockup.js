// Renders the real Activity rows into the mockup table the way board-activity.js
// does: filter to the selected window, newest first, then one <tr> per stamp.
// The board's "now" is the data's generatedAt so the relative ages match the
// board the rows were copied from.
(function () {
  var data = window.mockupBoardData;
  var now = new Date(data.generatedAt).getTime();
  var months = ["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"];
  function pad(n) { return n < 10 ? "0" + n : String(n); }
  function whenText(iso) {
    var d = new Date(iso);
    return months[d.getUTCMonth()] + " " + d.getUTCDate() + ", " + pad(d.getUTCHours()) + ":" + pad(d.getUTCMinutes()) + " UTC";
  }
  function agoText(iso) {
    var mins = Math.max(0, Math.round((now - new Date(iso).getTime()) / 60000));
    if (mins < 60) return mins + "min ago";
    var hours = Math.round(mins / 60);
    if (hours < 48) return hours + "h ago";
    return Math.round(hours / 24) + "d ago";
  }
  function cell(tag, text, header, cls) {
    var el = document.createElement(tag);
    if (tag === "th") el.setAttribute("scope", "row"); else el.setAttribute("headers", header);
    if (cls) el.className = cls;
    el.textContent = text;
    return el;
  }
  function render(hours) {
    var rows = data.activity.filter(function (r) { return now - new Date(r.stampAt).getTime() <= hours * 3600000; });
    rows.sort(function (a, b) { return a.stampAt < b.stampAt ? 1 : a.stampAt > b.stampAt ? -1 : 0; });
    var ids = {};
    rows.forEach(function (r) { ids[r.id] = true; });
    document.getElementById("activity-summary").textContent =
      rows.length + " transitions across " + Object.keys(ids).length + " REQs in the last " + (hours >= 48 ? (hours / 24) + " days" : hours + " hours");
    var body = document.getElementById("activity-table-body");
    body.textContent = "";
    rows.forEach(function (r) {
      var req = data.requests[r.id] || {};
      var tr = document.createElement("tr");
      tr.setAttribute("data-activity-request", r.id);
      tr.appendChild(cell("th", r.id));
      tr.appendChild(cell("td", req.title || "", "activity-table-column-title"));
      tr.appendChild(cell("td", req.status || "", "activity-table-column-status"));
      tr.appendChild(cell("td", r.transition || "", "activity-table-column-transition"));
      var when = cell("td", whenText(r.stampAt) + " ", "activity-table-column-when");
      var small = document.createElement("small"); small.textContent = "· " + agoText(r.stampAt); when.appendChild(small);
      tr.appendChild(when);
      tr.appendChild(cell("td", r.stampField || "", "activity-table-column-stamp"));
      body.appendChild(tr);
    });
    document.querySelectorAll("[data-activity-window]").forEach(function (b) {
      var on = Number(b.getAttribute("data-activity-window")) === hours;
      b.classList.toggle("is-active", on); b.setAttribute("aria-pressed", String(on));
    });
  }
  document.querySelectorAll("[data-activity-window]").forEach(function (b) {
    b.addEventListener("click", function () { render(Number(b.getAttribute("data-activity-window"))); });
  });
  var initial = Number(new URLSearchParams(location.search).get("hours")) || 24;
  render(initial);
})();

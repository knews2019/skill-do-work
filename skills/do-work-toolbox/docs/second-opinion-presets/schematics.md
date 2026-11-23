Make schematics: draw the mechanism as labelled technical drawings in inline SVG. Show components as boxes with their inputs and outputs, data and control flow as arrows, and the file path on each part, with a legend. Give one schematic of the whole system and one detailed schematic for each important part, and keep the prose to short captions that say what each drawing shows. Every fact still comes from the report or the repository.

1. Number the arrows in the order things happen, and give the same numbers as a step list under the drawing, one plain sentence per step. An arrow outside that order, such as a part that reads a file later, gets a label and no number.
2. Every arrow has a label in words. Write it level, never turned sideways, and place it beside its line so it covers no line, box or arrowhead.
3. When a part accepts a fixed set of requests or commands, name each one on the arrow that carries it. Requests that happen at different steps get separate arrows.
4. Fill each box with a light tint for its kind of part, and color each arrow for its kind of flow. The kinds come from the system you draw: program, access rule and stored data are examples of kinds of part, and request, data and control are examples of kinds of flow.
5. One color has one meaning on the whole page, a box color is never reused for an unrelated arrow, and red is kept for what is refused, failed or not covered. Use about five colors at most.
6. Color is never the only signal: a stored file or folder also has a dashed outline, control flow also has a dashed line, and the legend names every fill and every line.
7. A labelled frame, not a color, shows where a part runs, such as a machine or a process.
8. The drawing sits on a plain white sheet in the light theme and on a plain dark grey sheet in the dark theme. The sheet, the frames, the lines and the text use neutral greys with no color tint, and a frame has an outline and no fill, so only the parts and the arrows carry color.
9. Text keeps readable contrast on every fill, in the light theme and in the dark theme.
10. When you can, open the page in a browser at desktop and phone width and fix every check that fails.

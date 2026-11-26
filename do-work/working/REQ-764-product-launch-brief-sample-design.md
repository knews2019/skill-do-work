---
id: REQ-764
title: 'Add Product Launch Brief to Sample Designs'
status: claimed
created_at: 2026-11-26T19:38:31Z
user_request: UR-194
domain: ui-design
prime_files: []
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: false
maintenance: false
claimed_at: 2026-11-26T19:44:59Z
---
# Add Product Launch Brief to Sample Designs

## What

Add a reusable report named **Product Launch Brief** to the report gallery’s **Sample designs** section. Let readers try a product example, then follow that same example through the people, work, and decisions needed to launch it. Give the sample the same preview, open-report, and reuse options as its gallery neighbors.

## Context

The design reference is the password-protected German proposal for a daily spot-the-difference game. It follows a collectible card from its illustrated front and article excerpt through the player’s collection to the daily editorial and design workflow. Use original public-facing imagery and copy for the gallery sample. [Design reference](https://zingy-treacle-7f4c65.netlify.app/)

The discussion suggested a compact version with one flippable card, a collection view showing that card, a production workflow, and three launch questions. A thumbnail showing both faces would make the interaction discoverable. Suggested gallery question: “What will readers experience, and what would it take for our team to run it?” These are design suggestions for the builder.

## Red-Green Proof

**RED prompt/case:** Find Product Launch Brief under Sample designs and open the report.
**Why RED now:** The gallery has no sample by that name.
**GREEN when:** The named entry opens a self-contained public sample. A reader can inspect the product example, recognize it in the operating workflow, and see which decisions remain before launch. The gallery provides a working preview and the usual design reuse prompt.
**Validation:** Inferred during capture

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read the applicable agent rules and record the approach.
- [ ] **[APPLY]:** Build the sample and its gallery entry.
- [ ] **[UNIFY]:** Review the diff and verify the gallery entry, report, and reuse controls.

*Source: [Original request](../user-requests/UR-194/input.md)*

## Addendum (2026-11-26)

The user approved the design guidance discussed after capture. [Follow-up request](../user-requests/UR-195/input.md)

### Worked example

- Place the launch premise beside the product example on wide screens. Stack them in reading order on phones, with the example in reach before the operating plan.
- Use a readable editorial type hierarchy, with prose around 65 characters per line. Repeat the example’s title and imagery in the collection and workflow so readers recognize the same item.
- Give the collectible a distinct outline. Present the workflow as ordered steps with the responsible person beside each one. Set the launch questions as a short list. Use one consistent accent for controls and focus.
- If using the proposed two-sided card, give it a labeled button operable by touch and keyboard. Keep the flip brief and user-triggered. Switch faces immediately when reduced motion is enabled. Show both faces together in the static gallery preview.
- Keep the example self-contained and use original local assets. Label the assumptions behind any illustrative figures.

### Reusable guidance for “Use this design”

- Let the reader try one representative part of the product.
- Follow that same example through the people and steps needed to deliver it.
- End with the launch decisions those steps expose.

### Verification

Open the sample at desktop and phone widths. Follow the same item from the demo into the workflow and identify a launch question that arises from one of its steps. For a two-sided card, check both directions by keyboard and touch, verify visible keyboard focus, and verify that reduced motion switches faces immediately. Check that the gallery preview communicates the interaction before the reader opens the report.

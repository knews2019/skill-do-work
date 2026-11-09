---
id: UR-164
title: 'Closing second-opinion suggestion picks the best learning method for the subject'
created_at: 2026-11-09T12:02:46Z
requests: [REQ-712]
word_count: 89
---
## Summary
When an AI report is made and its closing band suggests a second opinion on that same report, the suggestion should pick the method that would best help the user understand or deliver the subject. The user gave ten "accelerated AI prompts" (ELI5, 80/20, sticky notes, Feynman technique, mind map, quiz me, visual study guide, first principles, Socrates, schematics) as a curated starting list. In capture the user added: the pick sits beside the existing layout suggestion, and the model may choose a better method outside the list.

## Full Verbatim Input
> ```
> Given these accelerated AI prompts, when you make an AI report and then suggest a new second opinion on the same report, make sure to pick the best prompt that would have help me understand or deliver that subject. 
> 
> """"
> 
> ELI5
> 80/20
> Make-sticky-notes
> Feynman-technique
> Make-mindmap
> Quizme
> Make-Visual-studyguide
> Firstprinciples
> Socrates
> Make-Schematics
> 
> 
> FEYNMAN-TECHNIQUE help me understand the 
> Make sticky notes
> Make mind map
> QUIZME
> MAKE-VISUAL-STUDYGUIDE FIRSTPRINCIPLES
> SOCRATES - teach by asking questions
> MAKE-SCHEMATIC
> 
> 
> Codes to facilitate learning in Ai bots
> ELI5
> 80/20
> Make-sticky-notes
> Feynman-technique
> Make-mindmap
> Quizme
> Make-Visual-studyguide
> Firstprinciples
> Socrates
> Make-Schematics
> 
> """"
> ```

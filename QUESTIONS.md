# QUESTIONS.md — RAG Evaluation Questions

Sample questions for evaluating the hybrid agentic RAG pipeline over the Chronicles of Aethelgard corpus. Each question tests different retrieval strengths across the three search modes: **text** (keyword/BM25), **semantic** (vector), and **rff** (hybrid fusion).

The corpus contains 188 documents, 2,303 chunks spanning characters, locations, events, factions, and lore — all cross-referenced with 3–5 natural entity links per entry.

---

## 1. Exact Entity Lookup (text mode)

These questions target a specific named entity. BM25 should dominate retrieval.

1. Who is Throndor the Wise and what did he found?
2. What is the Arcane College of Thaldor?
3. Where is the Stonehold Royal Armory located?
4. Who leads the Ironhand Guild?
5. What role does Shaman Vessska play in Lizardfolk society?
6. What is the Moonwell and where is it located?
7. Who is Mayor Ironscale and what city does he govern?
8. What is the Shadowfen Collective?
9. What happened during the Lizardfolk Unification?
10. What is the Order of the Radiant Dawn?

## 2. Descriptive / Conceptual (semantic mode)

These questions describe a concept without naming it directly. Vector search should retrieve the relevant chunk based on meaning.

11. What volcanic landmark is sacred to fire elementals and has an eternal flame?
12. Which organization smuggled tiefling refugees to safety during a persecution?
13. What sunken underwater city holds ancient Lizardfolk magical archives?
14. Which bridge magically records conversations spoken upon it?
15. What organization trains falcons and hawks for military reconnaissance?
16. Which landmark serves as a memorial to victims of religious persecution?
17. What facility produces the finest enchanted dwarven weapons and armor?
18. Which guild regulates the production of candles, soaps, and magical lighting?
19. What organization maps ley lines and warns travelers about dangerous magical zones?
20. Which prison facility houses dangerous rogue mages and uses abjuration specialists?

## 3. Cross-Reference / Multi-Hop (hybrid RFF)

These questions require connecting information across multiple documents. The answer may not be in a single chunk — the agent must retrieve and combine results.

21. Which characters fought on both sides of the Siege of Silverhaven, and what were their roles?
22. What is the relationship between Lady Celestia Moonwhisper and the Elven Withdrawal? How did she change that policy?
23. How did Grimjaw Ironhand's collaboration with Throndor the Wise influence dwarven weapon production in Stonehold?
24. Trace the trade route from Silverhaven to Crystal Bay — which guilds and organizations manage each leg?
25. What events led from the Elven Withdrawal to the Treaty of Silverhaven? Name every major event in between.
26. Which factions have headquarters or major operations in Thaldor, and how are they connected?
27. How does the Moonwell's healing magic relate to the Moonwhisper Estate and Lady Celestia's diplomatic work?
28. What role did Poppy Greenfields play during the Siege of Silverhaven, and which organizations did she coordinate with?
29. How did Chieftain Grommash Bloodaxe transform from warlord to diplomat, and what treaty did he sign?
30. What is the connection between the Pirate Wars, Admiral Hornbreaker, and the Crystal Bay Free Port declaration?

## 4. Category Aggregation (broad retrieval)

These questions ask about a whole category. Tests breadth of retrieval and deduplication.

31. List all the guilds in Silverhaven and their primary trade or craft.
32. What are all the major events that occurred during the Age of Shattered Skies?
33. Name every character who serves as a leader or ruler of a city, organization, or kingdom.
34. What religious traditions are practiced in Aethelgard, and which deities are worshipped?
35. What are all the landmarks located in the Southern Expanse?
36. Which factions operate in Crystal Bay, and what do they do there?
37. What are all the magical schools or traditions in the corpus?
38. List every dwarven organization, guild, or institution mentioned in the corpus.
39. Which characters have ties to both the Eastern Kingdoms and the Southern Expanse?
40. What events involved the Orcish people or Grommashar?

## 5. Negation / Distractor (precision testing)

These questions include tempting but incorrect matches. Tests whether retrieval avoids near-misses and returns the correct chunk.

41. Who is the gnomish leader of the Arcane College? (No gnomish leader exists — Archchancellor Marcus Thornfield is human)
42. What naval battle did Queen Elara personally command? (She coordinated logistics during the Siege of Silverhaven; Sarah Ironside commanded the defense)
43. Which dwarf founded the Arcane College of Thaldor? (Throndor the Wise, an elf, founded it — Grimjaw Ironhand collaborated but did not found it)
44. What is the capital of the Western Reaches? (Answer is Silvermoon, not Silverhaven which is the capital of the Eastern Kingdoms)
45. Who unified the Lizardfolk tribes? (Chief Ssithik with Shaman Vessska's support, not Admiral Hornbreaker)
46. Did Queen Elara sign the Treaty of Thronehold? (No — the Treaty of Silverhaven; Thronehold is a different city)
47. What species is Malakor Shadowdancer? (Tiefling-Dark Elf hybrid, not a full tiefling or full dark elf)
48. Which organization does Pyre the Unquenchable lead? (He doesn't lead an organization — he's an independent smith at Everflame Peak)
49. Is the Stormpeak Earthquake the same event as the Sundering? (No — they are separate events centuries apart)
50. Does the Forest Wardens organization operate in the Eastern Kingdoms? (No — they protect the Whispering Woods in the Western Reaches)

## 6. Fuzzy / Natural Language (semantic robustness)

These questions are phrased casually or imprecisely. Tests whether semantic search handles paraphrase and colloquial language.

51. Where do fire genasi go to forge magical weapons?
52. Who's in charge of the lizard people and what do they believe?
53. What's the deal with the bridge that spies on people?
54. How did the orcs go from being raiders to signing a peace treaty?
55. What's the dwarven approach to blessing weapons and armor?
56. Is there anywhere tieflings can go to remember their persecuted ancestors?
57. Who runs the criminal underground in Silverhaven?
58. What's the fancy elves-only school for magic?
59. Where would you find the biggest multi-faith religious building?
60. Which dragonborn leader handles diplomacy instead of just fighting?

## 7. Temporal / Causal (reasoning over events)

These questions require understanding chronological relationships and cause-effect chains. Multiple chunks often needed.

61. What caused the War of Shattered Skies, and what ended it?
62. How did the Stormpeak Earthquake affect dwarven mining and trade routes?
63. Why did the elves withdraw from the world, and why did they come back?
64. What was the sequence of events from the Elven Withdrawal to the founding of Silverhaven?
65. How did the Pirate Wars lead to the Crystal Bay Free Port declaration?
66. What role did the Great Persecution play in shaping tiefling organizations?
67. How did the Catfolk Trade Expansion change commerce across Aethelgard?
68. What was the Minotaur Naval Alliance, and how did it connect to anti-piracy efforts?
69. Why was the Battle of Thunder Pass strategically important in the War of Shattered Skies?
70. How did the Dwarven Surface Trade Expansion change the relationship between Stonehold and Silverhaven?

## 8. Geographic / Spatial (location reasoning)

These questions involve spatial relationships, regions, and travel between places.

71. What region separates the Eastern Kingdoms from the Western Reaches?
72. Where is Everflame Peak relative to Drakkenheim, and what connects them?
73. Which cities sit on the River Aethel, and how does the river influence trade?
74. What geographic features make Thunder Pass a strategic military location?
75. How does one travel from Silverhaven to Crystal Bay by sea, and what landmarks would they pass?
76. What are the Whispering Woods and why are access restrictions enforced there?
77. Which regions of Aethelgard are most affected by the Dragon Wastes?
78. What is the significance of the Northern Wastes, and who ventures there?
79. How does the Shattered Isles' geography reflect its formation during the Sundering?
80. What role does the Eastern Plains play in feeding the Eastern Kingdoms?

## 9. Relational / Comparative (entity comparison)

These questions ask about connections, similarities, or differences between entities. Tests multi-document retrieval and synthesis.

81. Compare the leadership styles of Queen Elara and King Thorin Stoneforge.
82. How do the Healers' Guild and the Apothecaries' Guild differ in their approach to medicine?
83. What is the relationship between the Bards' College and the Sages' Consortium?
84. How do the Forest Wardens differ from the Ranger Alliance?
85. Compare the Crystal Bay Naval Yard and the Stonehold Royal Armory — what does each produce?
86. What distinguishes the Arcane College Security Force from the Prison Wardens?
87. How does the Silverhaven Merchant Consortium relate to the Trade Federation of Crystal Bay?
88. Compare the religious practices of Moradin's followers and Pelor's followers.
89. What are the differences between necromancy ethics and healing magic in the corpus?
90. How do the Dragonborn Ancestral Halls and the Stonehold Treasury serve similar cultural purposes for different species?

## 10. Edge Cases / Stress Tests

These questions test boundaries — rare terms, deeply nested references, or ambiguous queries.

91. What is the significance of the number 342 in Aethelgard's history?
92. Which character has the longest lifespan mentioned in the corpus?
93. What material are the Silverhaven Gates enchanted with, and who enchanted them?
94. How many different species are represented in the Siege of Silverhaven's defender forces?
95. What is the only organization that has both criminal and humanitarian functions?
96. Which landmark was built on the site of an ancient elven settlement?
97. What spell is Throndor the Wise known for, and where is it still taught?
98. Who was the first non-dwarven runemaster, and who trained them?
99. What event happened in 440 AS that related to tiefling rights?
100. Which guild has the shortest name in the corpus, and what do they do?

---

## Expected Answer Characteristics

| Category | Ideal Search Mode | Key Retrieval Signals |
|----------|-------------------|-----------------------|
| Exact Entity Lookup | `text` | Proper nouns, exact names |
| Descriptive/Conceptual | `semantic` | Paraphrase, concept matching |
| Cross-Reference/Multi-Hop | `rff` | Multi-entity traversal |
| Category Aggregation | `rff` | Broad category terms |
| Negation/Distractor | `rff` + verification | Near-miss discrimination |
| Fuzzy/Natural Language | `semantic` | Paraphrase robustness |
| Temporal/Causal | `rff` | Event names + dates |
| Geographic/Spatial | `rff` | Place names + relationships |
| Relational/Comparative | `rff` | Multi-entity comparison |
| Edge Cases/Stress | `text` + `rff` | Rare terms, specific facts |

## Scoring Criteria

When evaluating RAG answers against these questions, consider:

1. **Retrieval accuracy** — Did the system return the correct source chunk(s)?
2. **Completeness** — Does the answer include all relevant information from the corpus?
3. **Hallucination rate** — Does the answer invent facts not present in the corpus?
4. **Cross-reference resolution** — Does the answer correctly follow entity links across documents?
5. **Mode appropriateness** — Which search mode performed best for which question type?
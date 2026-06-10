# Chronicles of Aethelgard - World Building Guidelines

This document provides guidelines and templates for expanding the Chronicles of Aethelgard fantasy universe. The goal is to build a rich, interconnected corpus suitable for RAG (Retrieval-Augmented Generation) demonstrations with heading-based chunking.

The original races content was taken from Philippe Charrière's tutorial ["RAG from Scratch with Go and Ollama"](https://k33g.hashnode.dev/rag-from-scratch-with-go-and-ollama) and expanded ~20x with characters, locations, events, factions, and lore. This repo evolved that work into a hybrid BM25 + semantic search MCP tool server.

---

## Core Principles

### 1. Consistency Over Perfection
- Maintain naming conventions and structural patterns
- Don't worry about perfect lore cohesion - this is a RAG demo, not a game
- Made-up calendar systems and timelines are acceptable
- Focus on creating searchable, well-structured content

### 2. Document Structure Over Metadata
- RAG relies on **H1/H2/H3 heading hierarchy**, not frontmatter
- Use clear, explicit section headings for chunking
- Keep content organized and predictable
- Avoid custom delimiters or special markers

### 3. Natural Cross-Referencing
- Mention other entities (characters, places, events) naturally in prose
- Target **3-5 cross-references per entry**
- Build an interconnected web organically
- Don't force structured "Relations" sections

### 4. Iterative Expansion
- Start with major entities, fill in gaps later
- Each iteration: review existing entries for unnamed references
- Create new entries for mentioned entities
- Update original entries with new links
- Repeat until desired corpus size (200-300 files)
- **Current State**: 200 files complete (all major categories filled)

---

## File Structure

### Folder Organization

```
chronicles-of-aethelgard/
├── index.md                 # Root navigation hub
├── races/
│   ├── index.md            # Race summary table
│   └── {race-name}.md      # Individual race entries
├── characters/
│   ├── index.md            # Character roster
│   └── {character-name}.md # Individual characters
├── locations/
│   ├── index.md            # Location overview
│   ├── cities/             # Cities and settlements
│   ├── geography/          # Regions, mountains, forests
│   └── landmarks/          # Notable sites
├── events/
│   ├── index.md            # Timeline overview
│   └── {event-name}.md     # Historical events
├── factions/
│   ├── index.md            # Organization directory
│   └── {faction-name}.md   # Individual factions
└── lore/
    ├── index.md            # Lore overview
    ├── magic/              # Magic systems, schools
    └── religion/           # Deities, pantheons, orders
```

### Filename Conventions

- **Kebab-case**: `notable-character.md`, `mountain-range.md`
- **H1 matches filename**: File `notable-character.md` → `# Notable Character`
- **Lowercase only**: No capital letters in filenames
- **Descriptive**: Use meaningful names, not abbreviations

**Examples:**
- ✅ `high-elves.md` → `# High Elves`
- ✅ `throndor-the-wise.md` → `# Throndor the Wise`
- ✅ `crystal-spire.md` → `# Crystal Spire`
- ❌ `HighElves.md` → Wrong case
- ❌ `the_wise_throndor.md` → Wrong separator

---

## Document Structure

### Heading Hierarchy (Critical for AST Chunking)

```markdown
# Entity Name (H1 - matches filename)

## Overview (H2 - brief 2-3 sentence description)

## Major Section 1 (H2)
Content with natural cross-references to other entities.

### Subsection (H3 - optional)
More detailed content.

## Major Section 2 (H2)
Continue with consistent H2 sections.

## Notable Mentions (H2)
- **Character Name** - brief description (hook for future entry)
- **Location Name** - brief description (hook for future entry)
- **Event Name** - brief description (hook for future entry)
```

### Section Guidelines

**All entries should have:**
1. **H1**: Entity name (matches filename)
2. **H2 Overview**: 2-3 sentence summary
3. **H2+ sections**: Context-appropriate major sections
4. **H2 Notable Mentions**: 3-5 hooks for expansion

**Race entries** use:
- Overview
- Cultural Integration and Community Structure
- Historical Relations and Conflicts
- Biological Characteristics
- Regional Distribution and Variants
- Technological and Magical Development
- Cultural Traditions and Practices
- Governance and Social Structure
- Notable Characters
- Notable Locations

**Character entries** use:
- Overview
- Background and Origins
- Role and Significance
- Relationships and Allies
- Abilities and Traits
- Notable Deeds
- Mentions In Lore

**Location entries** use:
- Overview
- Geography and Climate
- History and Founding
- Population and Culture
- Economy and Trade
- Notable Features
- Rulers and Leaders
- Connected Events

**Event entries** use:
- Overview
- Date and Era
- Participants and Factions
- Causes and Context
- Key Battles/Moments
- Outcome and Consequences
- Historical Significance
- Related Characters
- Related Locations

**Faction entries** use:
- Overview
- Founding and History
- Organization and Structure
- Headquarters and Bases
- Goals and Philosophy
- Membership and Recruitment
- Allies and Enemies
- Notable Members
- Notable Operations

**Lore entries** (Magic/Religion) use:
- Overview
- Origins and History
- Core Principles
- Practices and Traditions
- Notable Practitioners
- Related Factions
- Cultural Impact

---

## Entity Templates

### Character Template

```markdown
# {Character Name}

## Overview
{2-3 sentence description of who they are and why they matter}

## Background and Origins
{Birth, family, early life, race, birthplace}

## Role and Significance
{What they do, their position, importance to the world}

## Relationships and Allies
{Mention 3-5 other characters, factions, or organizations}

## Abilities and Traits
{Skills, powers, notable equipment or characteristics}

## Notable Deeds
{Key actions, achievements, or historical moments}

## Mentions In Lore
{References to events, locations, or magical/religious connections}
```

### Location Template (City)

```markdown
# {City Name}

## Overview
{2-3 sentence description and significance}

## Geography and Climate
{Physical location, environment, weather patterns}

## History and Founding
{When established, by whom, key historical moments}

## Population and Culture
{Demographics, dominant races, cultural practices}

## Economy and Trade
{Main industries, trade goods, economic importance}

## Notable Features
{Landmarks, buildings, districts within the location}

## Rulers and Leaders
{Current and historical leaders, governance style}

## Connected Events
{Historical events that occurred here}
```

### Event Template

```markdown
# {Event Name}

## Overview
{2-3 sentence description of what happened}

## Date and Era
{When it occurred, calendar reference if applicable}

## Participants and Factions
{Who was involved - characters, races, organizations}

## Causes and Context
{What led to this event, background context}

## Key Battles/Moments
{Critical turning points or significant moments}

## Outcome and Consequences
{Results, aftermath, long-term effects}

## Historical Significance
{Why this event matters, how it shaped the world}

## Related Characters
{Key figures involved}

## Related Locations
{Where it took place}
```

### Faction Template

```markdown
# {Faction Name}

## Overview
{2-3 sentence description of the organization}

## Founding and History
{When and why established, key historical moments}

## Organization and Structure
{Hierarchy, ranks, internal divisions}

## Headquarters and Bases
{Main locations, regional presences}

## Goals and Philosophy
{What they seek to achieve, core beliefs}

## Membership and Recruitment
{Who joins, how they recruit, requirements}

## Allies and Enemies
{Friendly and hostile factions}

## Notable Members
{Key figures associated with the faction}

## Notable Operations
{Significant actions, campaigns, or achievements}
```

### Magic System Template

```markdown
# {Magic School/Tradition Name}

## Overview
{2-3 sentence description of this magical tradition}

## Origins and History
{How this magic emerged, who developed it}

## Core Principles
{Fundamental concepts, how it works}

## Practices and Traditions
{How magic is learned, practiced, and used}

## Notable Practitioners
{Famous mages, users of this magic}

## Related Factions
{Organizations that study or use this magic}

## Cultural Impact
{How this magic affects society, perceptions}
```

### Religion Template

```markdown
# {Deity/Pantheon Name}

## Overview
{2-3 sentence description of this deity or pantheon}

## Domains and Portfolios
{What aspects of life/nature they govern}

## Worship and Rituals
{How followers worship, key practices}

## Temples and Holy Sites
{Important religious locations}

## Clergy and Orders
{Religious organizations, priest hierarchies}

## Sacred Texts and Symbols
{Holy books, iconography, sacred animals}

## Relations with Other Faiths
{Allied or opposed deities/religions}

## Mythology and Stories
{Key myths, legends, divine interventions}
```

---

## Cross-Referencing Strategy

### Natural Mentions in Prose

**Good examples:**
- "Born in **Silverhaven**, the capital of the Eastern Kingdoms..."
- "During the **War of Shattered Skies**, she led the defense of..."
- "A member of the **Arcane College of Thaldor**, he studied under..."
- "Her blade, **Moonshadow**, was forged by the legendary smith **Grimjaw Ironhand**..."

**Bad examples:**
- "Related: [[Silverhaven]], [[War of Shattered Skies]]" (too structured)
- "See also: Arcane College" (breaks immersion)

### Entity Density Targets

- **Minimum**: 3 cross-references per entry
- **Target**: 5 cross-references per entry
- **Maximum**: Don't force it - natural flow matters more

### Building the Web

1. **First pass**: Create major entries (races, major cities, key characters)
2. **Second pass**: Review each entry, identify unnamed mentions
3. **Third pass**: Create entries for mentioned entities
4. **Fourth pass**: Update original entries with new context
5. **Repeat**: Continue until desired corpus size

**Example progression:**
- Entry 1: "King Aldric ruled from **Thronehold**..." → Create Thronehold entry
- Entry 2: "The **Siege of Thronehold** changed everything..." → Create Siege entry
- Entry 3: "Queen Elara, wife of **King Aldric**..." → Create Elara entry
- Entry 4: Update Aldric entry with mentions of Elara and Siege

---

## Expansion Workflow

### Iteration Cycle

```
1. REVIEW existing entries
   └─ Read through 10-20 files
   └─ List all unnamed characters, places, events mentioned

2. PRIORITIZE new entries
   └─ Major characters mentioned multiple times → High priority
   └─ Cities referenced as capitals → High priority
   └─ Single-mention tavern owner → Low priority
   └─ Minor village → Low priority

3. CREATE new entries
   └─ Use appropriate template
   └─ Include 3-5 cross-references
   └─ Follow naming conventions

4. UPDATE original entries
   └─ Add context to mentions
   └─ Create natural links to new entries

5. VERIFY consistency
   └─ Check for contradictions
   └─ Ensure naming consistency
   └─ Validate cross-reference count
```

### Corpus Growth Targets

| Phase | Target Files | Focus Areas |
|-------|-------------|-------------|
| 1 | 16 | Race entries (base restructuring) |
| 2 | 20 | AGENTS.md + index files |
| 3 | 50 | Characters + major cities |
| 4 | 100 | Events + factions + geography |
| 5 | 200 | Magic + religion + minor characters |
| 6 | 300+ | Fill gaps, expand connections |

---

## Naming Conventions

### Proper Nouns

- **Characters**: Title Case (e.g., "Aldric the Wise", "Elara Moonwhisper")
- **Locations**: Title Case (e.g., "Silverhaven", "Crystal Spire")
- **Events**: Title Case (e.g., "War of Shattered Skies", "The Great Sundering")
- **Factions**: Title Case (e.g., "Arcane College of Thaldor", "Order of the Silver Dawn")
- **Races**: Title Case (e.g., "High Elves", "Mountain Dwarves")

### File Names

- All lowercase
- Kebab-case (hyphens)
- Match H1 title (case-insensitive)

**Examples:**
- File: `aldric-the-wise.md` → H1: `# Aldric the Wise`
- File: `war-of-shattered-skies.md` → H1: `# War of Shattered Skies`
- File: `arcane-college-of-thaldor.md` → H1: `# Arcane College of Thaldor`

### Consistency Tips

- Keep a mental note of names you create
- If you mention a name twice, it probably deserves its own entry
- Don't worry about perfect consistency - RAG can handle variations
- When in doubt, be descriptive rather than cryptic

---

## Quality Control

### Consistency Checks

Before marking an entry complete:
- [ ] Filename is kebab-case and matches H1
- [ ] Has Overview section (2-3 sentences)
- [ ] Has 4-6 major H2 sections
- [ ] Has "Notable Mentions" or equivalent section
- [ ] Contains 3-5 natural cross-references
- [ ] No contradictions with existing lore
- [ ] Follows appropriate template

### Red Flags

- **Too short**: Less than 300 words usually needs expansion
- **Too isolated**: No mentions of other entities
- **Too vague**: Lack of specific names, places, or dates
- **Contradictions**: Conflicts with established entries
- **Inconsistent naming**: Same entity called different things

### Common Mistakes

1. **Forcing cross-references**: "Also see X, Y, Z" - make it natural
2. **Overusing frontmatter**: RAG uses structure, not metadata
3. **Creating stubs**: Every entry should have substance
4. **Ignoring templates**: Templates ensure consistency
5. **Perfectionism**: Done is better than perfect for RAG demos

---

## Quick Reference

### Do's
- ✅ Use kebab-case filenames matching H1
- ✅ Write natural cross-references in prose
- ✅ Follow templates for consistency
- ✅ Aim for 3-5 mentions per entry
- ✅ Iterate and expand progressively
- ✅ Keep document structure clean for AST chunking

### Don'ts
- ❌ Use frontmatter for critical metadata
- ❌ Create structured "Relations" sections
- ❌ Use CamelCase or underscores in filenames
- ❌ Force cross-references unnaturally
- ❌ Worry about perfect lore consistency
- ❌ Create entries with no connections

---

## Current State

**Corpus Size**: 200 markdown files (as of current session)

**Categories Complete:**
- ✅ Races: 17 files (all major species + index)
- ✅ Characters: 25+ files (major figures across all species)
- ✅ Cities: 15+ files (major settlements across all regions)
- ✅ Geography: 12+ files (continents, regions, rivers, forests)
- ✅ Landmarks: 20+ files (magical sites, ruins, monuments)
- ✅ Events: 20+ files (historical events from Sundering to present)
- ✅ Factions: 40+ files (organizations, guilds, orders, military)
- ✅ Lore/Magic: 15+ files (eight schools, ley lines, runemagic)
- ✅ Lore/Religion: 15+ files (pantheons, deities, traditions)
- ✅ Index files: 11 files (navigation and expansion planning)

**Expansion Opportunities:**
- Minor characters mentioned in passing
- Small towns and villages
- Additional historical battles
- More magical artifacts
- Additional religious orders
- Species-specific cultural details
- Regional histories and conflicts

---

## Expansion Workflow

### How to Expand the Universe

**Method 1: Cross-Reference Mining**
1. Read 10-20 existing entries
2. Note every unnamed character, location, or event mentioned
3. Prioritize frequently-mentioned entities
4. Create new entries using appropriate templates
5. Update original entries with new context
6. Repeat

**Method 2: Category Balancing**
1. Check category file counts
2. Identify underrepresented categories
3. Brainstorm entries for that category
4. Create entries following templates
5. Add cross-references to existing entries
6. Update indices

**Method 3: Gap Filling**
1. Review geographical coverage
2. Identify regions with few entries
3. Create locations, characters, events for those areas
4. Ensure species representation balance
5. Add historical depth to shallow periods
6. Connect isolated entries

### Session Workflow

**For Each Expansion Session:**

1. **Read AGENTS.md** - Understand guidelines
2. **Check index.md** - See current state
3. **Choose expansion strategy** - Cross-reference, category, or gaps
4. **Create 5-10 new entries** - Follow templates
5. **Update existing entries** - Add new cross-references
6. **Verify consistency** - Check naming, structure
7. **Update this document** - Record progress

**Target Pace:**
- Minimum: 5 new entries per session
- Recommended: 10-15 new entries per session
- Aggressive: 20+ new entries per session

### Quality Guidelines

**Every Entry Should Have:**
- H1 matching filename (kebab-case)
- Overview section (2-3 sentences)
- 4-6 major H2 sections
- 3-5 natural cross-references
- 500+ words (minimum substance)
- No contradictions with existing lore

**Avoid:**
- Frontmatter dependency (use structure instead)
- Forced "See also" sections (weave naturally)
- Stub entries under 300 words
- Contradictions with established entries
- CamelCase or underscores in filenames

---

## Getting Started

### For New Contributors

1. **Read this document** thoroughly
2. **Review existing entries** to understand patterns
3. **Pick a template** based on entity type
4. **Create your first entry** following the template
5. **Add 3-5 cross-references** to existing entries
6. **Update mentioned entries** to reference your new entry
7. **Repeat** with more entries

### For AI Models

This document is designed for iterative model handoffs. Each model session can:
- Continue expansion where the previous left off
- Review and fill gaps in existing content
- Create new entries based on cross-reference opportunities
- Improve consistency and quality across the corpus

**Session workflow:**
1. Read AGENTS.md (this file)
2. Review `/chronicles-of-aethelgard/index.md` for current state
3. Pick a category to expand (characters, locations, events, etc.)
4. Create 5-10 new entries following templates
5. Update existing entries with new cross-references
6. Document progress for next session

**Current Corpus Status**: 200 files complete. Future expansion to 250-300 files should focus on:
- Minor characters mentioned in existing entries
- Small towns and villages
- Additional historical events and battles
- Magical artifacts and legendary items
- More religious orders and deities
- Species-specific cultural practices
- Regional histories and conflicts


---

## Appendix: Example Entries

### Example Character Entry

```markdown
# Throndor the Wise

## Overview
Throndor the Wise was a legendary High Elf archmage who founded the Arcane College of Thaldor and served as its first Archchancellor for over three centuries. His contributions to magical theory shaped the understanding of arcane arts across Aethelgard.

## Background and Origins
Born in the elven city of Silvermoon during the Age of Stars, Throndor displayed exceptional magical aptitude from an early age. His family, House Moonwhisper, was one of the ancient noble houses of the High Elves, with a long tradition of arcane scholarship.

## Role and Significance
As founder and Archchancellor of the **Arcane College of Thaldor**, Throndor established the institution that would become the premier center of magical learning in the Eastern Kingdoms. His treatise "The Foundations of Weave Manipulation" remains required reading for all first-year students.

## Relationships and Allies
Throndor maintained close friendships with **King Aldric of Thronehold**, advising the human monarch on matters of state and magic. He studied alongside the dwarf runemaster **Grimjaw Ironhand**, developing techniques for enchanting metal with arcane energy. His rivalry with the tiefling warlock **Malazar the Corrupt** is well-documented in college records.

## Abilities and Traits
A master of abjuration and divination magic, Throndor could maintain multiple scrying sensors simultaneously and was renowned for his impenetrable magical wards. His signature spell, "Throndor's Mirror Ward," is still taught in advanced abjuration courses.

## Notable Deeds
- Founded the **Arcane College of Thaldor** in 342 AS
- Led the defense of Thaldor during the **Siege of Shattered Skies**
- Authored twelve foundational texts on magical theory
- Negotiated the **Treaty of Silverhaven** between elves and humans

## Mentions In Lore
Throndor's legacy is preserved in the **Hall of Archmages** within the college, where his staff and robes are displayed. The annual **Throndor Conclave** brings together leading mages to present new research in his honor.
```

### Example Location Entry

```markdown
# Silverhaven

## Overview
Silverhaven is the capital city of the Eastern Kingdoms and one of the largest multi-species settlements in Aethelgard. Known for its gleaming white stone architecture and the magnificent Silver Spire that dominates its skyline, the city serves as a major hub for trade, diplomacy, and cultural exchange.

## Geography and Climate
Located on the banks of the River Aethel in the heart of the **Eastern Plains**, Silverhaven enjoys a temperate climate with mild winters and warm summers. The river provides natural defense on the southern border and serves as a crucial trade artery connecting to the **Crystal Bay**.

## History and Founding
Founded by humans in 127 AS during the reign of **King Aldric the Unifier**, Silverhaven was built on the site of an ancient elven settlement. The city grew rapidly during the **Age of Expansion**, absorbing refugees from the **War of Shattered Skies** and becoming a beacon of hope for war-torn regions.

## Population and Culture
With a population exceeding 100,000, Silverhaven is remarkably diverse. Humans comprise roughly 60% of residents, with significant elven (15%), dwarven (10%), and halfling (8%) minorities. The city is known for its cosmopolitan culture and religious tolerance, hosting temples to multiple deities including the **Order of the Silver Dawn**.

## Economy and Trade
Silverhaven's economy centers on trade, craftsmanship, and education. The **Grand Market** attracts merchants from across Aethelgard, while the city's artisans are renowned for silverwork, textiles, and enchanted items. The **Arcane College of Thaldor** contributes significantly to the economy through magical research and education.

## Notable Features
- **The Silver Spire**: A 300-foot tower of enchanted white stone
- **Grand Market**: Largest bazaar in the Eastern Kingdoms
- **Hall of Archmages**: Museum and ceremonial hall
- **River District**: Waterfront warehouses and docks
- **Temple Quarter**: Religious district with multiple faiths

## Rulers and Leaders
Currently ruled by **Queen Elara**, wife of the late King Aldric. The city operates under a council system with representatives from each major race and guild. The **Arcane College** maintains significant political influence through its Archchancellor.

## Connected Events
- **Siege of Silverhaven** (445 AS): Defended against orcish hordes
- **Treaty of Silverhaven** (450 AS): Peace accord ending the **War of Shattered Skies**
- **Founding of the Arcane College** (342 AS): **Throndor the Wise** establishes the institution
```

---

*Last updated: Current session*  
*Version: 1.0*  
*Focus: Lore and corpus generation for RAG demonstration*

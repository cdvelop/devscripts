# Badges SVG Migration Plan

## Overview
Migrate from HTML/CSS badges to SVG badges for better GitHub compatibility. The current HTML badges work locally but not on GitHub remote repositories.

## Requirements
- Create autonomous `badges.sh` script for SVG badge generation
- Maintain current double-label visual style (gray label + colored value)
- Generate single SVG file containing all badges
- Place SVG in `.github/` directory as `badges.svg`
- Support for n number of badges via parameters
- Native bash processing (Windows/Linux compatible)
- Testable with existing test framework

## Parameters Format
```bash
badges.sh "license:MIT:#007acc" "go:1.22:#00add8" "coverage:85%:#28a745"
```
- Format: `"label:value:color"`
- Separator: `:` (colon)
- Each badge as separate parameter
- Use colors as provided (no validation)
- Show errors for malformed parameters but continue processing others

## SVG Specifications
- **Height**: 20px (GitHub standard)
- **Font**: System default
- **Style**: Double-label (gray label + colored value)
- **Spacing**: Similar to current CSS (0.7ch equivalent)
- **File**: Single `badges.svg` in `.github/` directory

## Implementation Plan

### Phase 1: Core SVG Generation ✅
- [x] Create `badges.sh` script
- [x] Implement parameter parsing with `:` separator
- [x] Generate basic SVG structure
- [x] Implement double-label badge style
- [x] Handle multiple badges in single SVG

### Phase 2: Error Handling & Validation ✅
- [x] Add error handling for malformed parameters
- [x] Use functions.sh error/warning functions
- [x] Continue processing on individual badge errors
- [x] Validate minimum required parameters

### Phase 3: Testing ✅
- [x] Create comprehensive test suite
- [x] Test with `badgets.html` reference
- [x] Test SVG generation in `.github/` directory
- [x] Verify all current tests pass

### Phase 4: Integration ✅
- [x] Update `badgets.html` to use SVG link only
- [x] Integrate with `gobadge.sh` (complete migration)
- [x] Update existing workflows and tests
- [x] Fix SVG width calculation bug
- [x] Update test expectations for SVG approach

## Current Status: MIGRATION COMPLETE ✅
- All tests passing (badges.sh and gobadge.sh integration)
- Error handling working correctly
- SVG generation functional with proper width calculation
- Multiple badges supported
- Integration with existing workflow complete
- Documentation updated

## Migration Results
✅ **Successfully migrated from HTML/CSS to SVG badges**
- Better GitHub compatibility (SVG works in remote repositories)
- Maintained visual consistency (double-label design)
- Preserved all functionality from original HTML/CSS approach
- Added comprehensive error handling and testing
- Cross-platform native Bash implementation

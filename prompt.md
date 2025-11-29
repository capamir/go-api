# Task: Generate Complete Project Documentation

## Objective
Create comprehensive technical documentation for this project by analyzing the actual codebase. This documentation must be 100% accurate—based only on what exists in the code, not assumptions or best practices you think should be there.

---

## CRITICAL: Anti-Hallucination Rules

ONLY document what you can verify in the code:
- ✅ If a file exists → document it
- ✅ If a function/class exists → document it
- ✅ If a pattern is used → document it
- ❌ If you're unsure → mark as "[NEEDS VERIFICATION]"
- ❌ If it doesn't exist → DON'T mention it
- ❌ If you can't find proof in code → DON'T guess

Verification requirements:
- Include file paths for every claim
- Quote actual function/class/module names
- Reference actual configuration variables
- Cite actual error messages/codes
- Show actual command syntax

---

## Part 1: Project Structure & Architecture

### 1.1 Complete File Tree
List EVERY file in the project with its purpose:

project-root/
├── source-directory/
│ ├── module-a/
│ │ ├── file1 # Brief description of purpose
│ │ └── file2 # Brief description of purpose
│ └── module-b/
├── test-directory/
├── documentation/
├── configuration-files
└── build-scripts


Requirements:
- No files skipped or assumed
- Each file gets 1-line description based on actual code
- Use actual directory names from the project

### 1.2 Architecture Overview
Document the actual architecture found in code:
- Design patterns (Factory, Singleton, Observer, Strategy, etc.) - cite where used
- Architectural style (Client-Server, Layered, Event-Driven, MVC, etc.)
- Key abstractions (base classes, interfaces, protocols) - list actual names
- Component relationships (which parts depend on which)
- Data flow (how information moves through the system)

### 1.3 Technology Stack
List actual technologies found in code/config files:
- Programming language(s) and version(s)
- Frameworks and libraries (with versions if in lock files)
- External services/APIs
- Data storage systems
- Build/deployment tools

---

## Part 2: Core Functionality Documentation

For EVERY major feature/capability in the project:

### Feature: [Actual feature name from code]

Entry Points:
- Files: actual/path/to/implementation
- Functions/Classes: ActualClassName.method_name()
- Commands/Routes: actual command or endpoint
- Triggers: [How it's activated]

What It Does:
[User-facing description - what problem does it solve]

How It Works:
1. [Step-by-step flow citing actual function names]
2. [Which modules/classes are involved]
3. [Data transformations or state changes]

Input:
- Expected format/type
- Validation rules (from actual validators)
- Example: [from tests or actual code]

Processing:
- Core algorithm or business logic
- External calls (APIs, databases, services)
- Data transformations

Output:
- Return format/type
- How results are delivered
- Example: [from tests or actual code]

Error Handling:
- Possible errors (from actual error handling)
- Error codes/messages (actual text)
- Retry or fallback behavior

Configuration:
- Config variables: ACTUAL_VAR_NAME
- Config files involved
- Default values (from code)

Constraints:
- Limits (timeouts, sizes, counts - from constants)
- Supported types/formats
- Known limitations

---

## Part 3: Technical Systems

### 3.1 Configuration
- How config is loaded (actual files/functions)
- All config variables (list every one found)
- Priority/precedence order
- Default values (from code, not guessed)

### 3.2 Error Handling
- Error handling approach (centralized/distributed)
- Error types defined
- Error propagation flow
- User-facing messages (actual text)

### 3.3 Data Management
- What data is stored/cached (if any)
- Storage mechanism (files, database, memory, etc.)
- Data structures/schemas (actual definitions)
- Persistence strategy

### 3.4 External Integrations
For each external service/API:
- Service name and purpose
- Integration location: actual/file/path
- Authentication method (from code)
- Endpoints/operations used
- Request/response formats
- Error handling and retries

### 3.5 Testing
- Test framework(s) used
- Test organization (directory structure)
- How to run tests (actual commands)
- Coverage approach
- Test patterns observed

---

## Part 4: Development & Operations

### 4.1 Setup & Installation
- Prerequisites (from dependency files)
- Installation steps (actual commands)
- How to run (actual commands)
- Environment differences (dev/staging/prod)

### 4.2 Build & Deployment
- Build process (actual commands/scripts)
- Deployment method
- Environment requirements
- CI/CD configuration (if present)

### 4.3 Code Organization
- Module separation logic
- Naming conventions (observed patterns)
- Code structure principles
- Dependency management approach

---

## Validation Checklist

Before submitting, verify:
- [ ] Every file in project is listed
- [ ] Every claim cites actual file path or code reference
- [ ] All examples are from actual code/tests
- [ ] All error messages are actual text from code
- [ ] All config variables are from actual config files
- [ ] All limits/values are from actual constants
- [ ] No assumptions or "should have" statements
- [ ] No undocumented features mentioned
- [ ] Architecture matches actual code patterns

---

## Output Format

Single markdown file: docs/PROJECT_DOCUMENTATION.md

Structure:
1. Project Overview (what it is and does)
2. Project Structure & Architecture
3. Features & Functionality (one section per feature)
4. Technical Systems
5. Development & Operations
6. Configuration Reference

---

## Quality Standards

- Accuracy: 100% code-verified, 0% assumptions
- Completeness: Every feature, file, system documented
- Clarity: Technical but readable
- Structure: Consistent format throughout
- Maintainability: Easy to update when code changes

---

## Critical Reminder

This documentation is a faithful mirror of the codebase**—not a design document, roadmap, or wish list.

**Document only:
- What exists in the code
- How it actually works
- What it actually does

DO NOT document:
- What "should" exist
- What "could" be added
- What "best practices suggest"
- Features you assume exist

If something is unclear, incomplete, or missing from the code, explicitly state that rather than filling gaps with assumptions.

---

## Examples of Generic Placeholders

When citing code, use actual names:
- ✅ UserAuthenticator.validate_token()
- ❌ some_function()

When showing file paths, use actual paths:
- ✅ src/auth/authentication.py
- ❌ path/to/file

When listing config, use actual variable names:
- ✅ DATABASE_URL, API_KEY
- ❌ configuration_variable

The more specific and factual, the better.
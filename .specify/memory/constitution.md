&lt;!--
SYNC IMPACT REPORT

- Version change: none → 1.0.0
- Added Principles:
  - I. 严格的代码风格 (Strict Code Style)
  - II. 提交前静态检查 (Pre-Commit Static Analysis)
  - III. 全面的单元测试 (Comprehensive Unit Testing)
  - IV. 端到端质量保证 (End-to-End Quality Assurance)
- Removed sections:
  - Removed placeholder PRINCIPLE_5
  - Removed placeholder SECTION_2
  - Removed placeholder SECTION_3
- Templates requiring updates:
  - ✅ .specify/templates/plan-template.md
  - ✅ .specify/templates/tasks-template.md
  - ✅ .specify/templates/agent-file-template.md
  - ✅ .specify/templates/checklist-template.md
  - (✅ .specify/templates/spec-template.md - no changes needed)
--&gt;

# Xizexcample Constitution

## Core Principles

### I. 严格的代码风格 (Strict Code Style)
所有提交的代码都必须使用 `gofmt` 进行格式化，并遵循 `golint` 的规范。这是为了确保整个项目代码风格的一致性、可读性和可维护性。统一的风格可以减少代码审查中的噪音，让开发者更专注于逻辑本身。

### II. 提交前静态检查 (Pre-Commit Static Analysis)
在每次提交代码之前，必须执行并通过项目的公共 `lint` 检查。任何 `lint` 错误或警告都必须在提交前修复。这是一道关键的质量门禁，可以在早期发现潜在的错误、坏味道和不规范的写法，防止问题代码流入代码库。

### III. 全面的单元测试 (Comprehensive Unit Testing)
所有新增的功能模块、核心逻辑和公共函数都必须有配套的单元测试。测试需要覆盖正常的业务场景、边界条件和异常情况。单元测试是保证代码质量、功能正确性和支持未来重构的基石。

### IV. 端到端质量保证 (End-to-End Quality Assurance)
所有面向用户的新功能或重要流程变更，都必须编写并通页面级的端到端（E2E）测试。E2E 测试从用户的角度模拟真实操作，确保各个模块集成后整个系统的功能是完整和正确的，弥补单元测试无法覆盖的集成问题。

## Governance
本章程是项目的最高行为准则，所有代码贡献都必须严格遵守。对本章程的任何修订都必须经过团队评审，并更新版本号和修订日期。

**Version**: 1.0.0 | **Ratified**: 2025-10-10 | **Last Amended**: 2025-10-10
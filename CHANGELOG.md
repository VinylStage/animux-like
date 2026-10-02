# 변경 이력

이 프로젝트의 모든 주요 변경 사항은 이 파일에 기록됩니다.

형식은 [Keep a Changelog](https://keepachangelog.com/ko/1.1.0/)를 따르며,
이 프로젝트는 [Semantic Versioning](https://semver.org/lang/ko/)을 준수합니다.

---

## [v0.1.0] - 2026-10-02

### 추가됨

- 터미널 가상 펫 CLI 최초 릴리스
- 백그라운드 데몬 (`animux daemon`) — Unix Domain Socket 기반, XDG 표준 준수
- 펫 입양 명령어 (`animux adopt <종> <이름>`) — 펭귄, 고양이, 강아지 지원
- 펫 상호작용 명령어 — `feed`, `play`, `clean`, `status`
- 실시간 TUI 관찰 모드 (`animux show`) — Bubble Tea 기반, `f/p/c/q` 단축키 지원
- 다마고치 스타일 상태 관리 — 배고픔, 행복도, 청결도 시간 경과 감소
- 방치 시 병에 걸리는 메커니즘
- 구조화된 JSON 로그 (`log/slog`, DEBUG 레벨)
- XDG Base Directory 표준 준수 (상태 파일, 로그, 소켓)
- GoReleaser 기반 자동 배포 — APT(Gemfury), Snap Store, Homebrew Tap

[v0.1.0]: https://github.com/VinylStage/animux-like/releases/tag/v0.1.0

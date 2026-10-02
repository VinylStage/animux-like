# 기여 가이드

Animux에 관심을 가져주셔서 감사합니다! 버그 리포트, 기능 제안, 코드 기여 모두 환영합니다.

## 시작하기 전에

- 기존 [이슈](https://github.com/VinylStage/animux-like/issues)와 [PR](https://github.com/VinylStage/animux-like/pulls)을 먼저 확인하여 중복을 방지해 주세요.
- 큰 변경 사항은 PR 전에 먼저 이슈를 열어 논의해 주세요.
- 이 프로젝트는 [행동 강령](CODE_OF_CONDUCT.md)을 따릅니다.

## 개발 환경 설정

```bash
# 1. 저장소 포크 후 클론
git clone https://github.com/YOUR_USERNAME/animux-like.git
cd animux-like

# 2. 의존성 설치
go mod tidy

# 3. 빌드
go build -o animux

# 4. 테스트 실행
go test ./...
```

## 프로젝트 구조

```
animux-like/
├── main.go              # 진입점
├── cmd/                 # CLI 명령어 정의 (Cobra)
│   ├── root.go
│   ├── adopt.go
│   ├── feed.go
│   ├── play.go
│   ├── clean.go
│   ├── status.go
│   ├── show.go
│   └── daemon.go
├── daemon/              # 백그라운드 데몬 및 소켓 클라이언트
│   ├── daemon.go
│   └── client.go
├── pet/                 # 펫 상태 관리 및 종 정의
│   ├── state.go
│   └── species.go
├── ui/                  # ASCII 아트 및 Lipgloss UI
│   └── ascii.go
└── Formula/             # Homebrew Formula (자동 생성)
```

## 기여 절차

1. `main` 브랜치에서 새 브랜치를 생성합니다.
   ```bash
   git checkout -b feat/새기능이름
   ```

2. 변경 사항을 작성하고 테스트를 실행합니다.
   ```bash
   go test ./...
   go vet ./...
   ```

3. 커밋 메시지는 [Conventional Commits](https://www.conventionalcommits.org/) 형식을 따릅니다.
   ```
   feat: 새로운 동물 종 추가 (토끼)
   fix: 데몬 소켓 연결 오류 수정
   docs: README 설치 가이드 업데이트
   ```

4. PR을 열고 변경 사항을 설명해 주세요.

## 코딩 스타일

- Go 표준 포맷 도구를 사용합니다: `gofmt -w .`
- 새로운 동물 종 추가 시 `pet/species.go`에 ASCII 아트와 특성을 함께 정의해 주세요.
- 로그는 반드시 `log/slog`의 구조화된 로그를 사용합니다.

## 새로운 동물 종 추가하기

1. `pet/species.go`에 종 상수와 특성(배고픔 감소 속도 등)을 추가합니다.
2. `ui/ascii.go`에 해당 동물의 ASCII 아트를 상태별(기본, 먹기, 놀기)로 추가합니다.
3. 테스트를 작성하고 PR을 보내주세요!

# Animux 🐾

> 터미널에서 키우는 가상 펫 — Linux/macOS 네이티브 CLI 다마고치

**Animux**는 유닉스 철학을 따르는 터미널 가상 펫입니다. 백그라운드 데몬이 시간 경과에 따라 펫 상태를 실시간으로 관리하며, 표준 CLI 명령어와 실시간 TUI 관찰 모드를 통해 펫과 상호작용할 수 있습니다.

[![GitHub release](https://img.shields.io/github/v/release/VinylStage/animux-like)](https://github.com/VinylStage/animux-like/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/VinylStage/animux-like)](https://goreportcard.com/report/github.com/VinylStage/animux-like)

## 주요 기능

- 🐧 **다양한 동물 종류** — 펭귄, 고양이, 강아지 중 선택하여 입양
- ⏱️ **실시간 생애 주기** — 데몬이 백그라운드에서 배고픔, 행복도, 청결도를 지속 관리
- 🎨 **실시간 관찰 모드** — `animux show` 명령어로 귀여운 ASCII 아트 TUI 화면 표시
- 📁 **XDG 표준 준수** — 상태 파일, 로그, 소켓 모두 Linux 표준 디렉토리 사용
- 🔌 **데몬 아키텍처** — Unix Domain Socket 기반 경량 백그라운드 프로세스

## 설치

### Ubuntu / Debian (`apt`)

```bash
echo "deb [trusted=yes] https://apt.fury.io/vinylstage/ /" | sudo tee /etc/apt/sources.list.d/animux.list
sudo apt update
sudo apt install animux
```

### Snap Store (`snap`)

```bash
sudo snap install animux-like
```

### macOS / Linux (`brew`)

```bash
brew tap VinylStage/animux-like https://github.com/VinylStage/animux-like
brew install animux
```

### 직접 빌드

```bash
git clone https://github.com/VinylStage/animux-like.git
cd animux-like
go build -o animux
sudo mv animux /usr/local/bin/
```

## 빠른 시작

**1. 백그라운드 데몬 실행**

```bash
animux daemon &
```

**2. 펫 입양**

```bash
# 지원 종류: penguin, cat, dog
animux adopt penguin Pingu
```

**3. 펫 돌보기**

```bash
animux status   # 현재 상태 확인
animux feed     # 밥 주기
animux play     # 놀아주기
animux clean    # 청소하기
```

**4. 실시간 관찰 모드**

```bash
animux show
# 단축키: f(밥), p(놀기), c(청소), q(종료)
```

## 파일 저장 위치 (XDG 표준)

| 종류 | 경로 |
|------|------|
| 저장 파일 | `~/.local/share/animux/state.json` |
| 로그 | `~/.local/state/animux/animux.log` |
| Unix 소켓 | `$XDG_RUNTIME_DIR/animux.sock` |

## 기여하기

기여를 환영합니다! 자세한 내용은 [CONTRIBUTING.md](CONTRIBUTING.md)를 참고해 주세요.

## 라이선스

이 프로젝트는 [MIT 라이선스](LICENSE) 하에 배포됩니다.

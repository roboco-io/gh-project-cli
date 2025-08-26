# Pre-PRD: ghp-cli

## 1. 개요

### 1.1 제품명
ghp-cli (GitHub Project CLI Tool)

### 1.2 문제 정의
GitHub CLI(gh)는 GitHub의 많은 기능을 커맨드라인에서 효율적으로 사용할 수 있게 해주지만, GitHub Projects(beta)에 대한 지원이 제한적이다. 특히 다음과 같은 작업들이 불편하거나 불가능하다:
- 프로젝트 아이템의 일괄 관리
- 커스텀 필드 값 업데이트
- 프로젝트 뷰 필터링 및 관리
- 자동화 워크플로우 실행
- 프로젝트 간 아이템 이동

### 1.3 목표
GitHub Projects를 커맨드라인에서 완벽하게 제어할 수 있는 강력하고 직관적인 CLI 도구를 제공하여 개발자의 생산성을 향상시킨다.

## 2. 타겟 사용자

### 2.1 주요 사용자
- **개발자**: 터미널 환경에서 프로젝트 관리를 선호하는 개발자
- **DevOps 엔지니어**: CI/CD 파이프라인에서 프로젝트 관리 자동화가 필요한 엔지니어
- **프로젝트 매니저**: 대량의 이슈와 PR을 효율적으로 관리해야 하는 PM
- **오픈소스 메인테이너**: 여러 프로젝트를 동시에 관리하는 메인테이너

### 2.2 사용 시나리오
- 스프린트 계획 시 이슈 일괄 할당
- 자동화된 프로젝트 상태 업데이트
- 프로젝트 보고서 생성
- 크로스 레포지토리 프로젝트 관리

## 3. 핵심 기능

### 3.1 프로젝트 관리
```bash
# 프로젝트 목록 조회
ghp list [--org <org>] [--user <user>]

# 프로젝트 생성
ghp create <name> [--org <org>] [--template <template>]

# 프로젝트 삭제
ghp delete <project-id>

# 프로젝트 복제
ghp clone <project-id> <new-name>
```

### 3.2 아이템 관리
```bash
# 아이템 추가
ghp item add <project-id> --issue <issue-url>
ghp item add <project-id> --pr <pr-url>
ghp item add <project-id> --draft <title>

# 아이템 조회
ghp item list <project-id> [--filter <filter>]

# 아이템 업데이트
ghp item update <item-id> --field <field-name> --value <value>

# 아이템 삭제
ghp item delete <item-id>

# 아이템 이동
ghp item move <item-id> --to-project <project-id>
```

### 3.3 필드 관리
```bash
# 필드 목록 조회
ghp field list <project-id>

# 필드 생성
ghp field create <project-id> <field-name> --type <type>

# 필드 업데이트
ghp field update <field-id> --name <new-name> --options <options>

# 필드 삭제
ghp field delete <field-id>
```

### 3.4 뷰 관리
```bash
# 뷰 목록 조회
ghp view list <project-id>

# 뷰 생성
ghp view create <project-id> <view-name> --layout <layout>

# 뷰 필터 설정
ghp view filter <view-id> --filter <filter-expression>

# 뷰 삭제
ghp view delete <view-id>
```

### 3.5 자동화 및 워크플로우
```bash
# 워크플로우 실행
ghp workflow run <workflow-name> --project <project-id>

# 일괄 작업
ghp bulk update --file <csv-file>
ghp bulk import --source <source> --project <project-id>

# 템플릿 적용
ghp template apply <template-id> --project <project-id>
```

### 3.6 보고서 및 분석
```bash
# 프로젝트 통계
ghp stats <project-id>

# 번다운 차트 생성
ghp burndown <project-id> --sprint <sprint-id>

# 보고서 생성
ghp report <project-id> --format <format> --output <file>
```

## 4. 기술 요구사항

### 4.1 기술 스택
- **언어**: Go (gh CLI와의 일관성)
- **API**: GitHub GraphQL API v4
- **의존성**: 
  - github.com/cli/cli/v2
  - github.com/shurcooL/graphql
  - github.com/spf13/cobra

### 4.2 호환성
- gh CLI 2.0+ 필수
- macOS, Linux, Windows 지원
- GitHub Enterprise Server 지원

### 4.3 성능 요구사항
- 대량 작업 시 배치 처리 지원
- API 호출 최적화 (GraphQL 쿼리 최적화)
- 로컬 캐싱 지원

## 5. 사용자 경험

### 5.1 설치
```bash
# Homebrew로 설치 (macOS/Linux)
brew install ghp-cli

# 또는 바이너리 직접 다운로드
curl -L https://github.com/roboco-io/gh-project-cli/releases/latest/download/ghp-$(uname -s)-$(uname -m) -o /usr/local/bin/ghp
chmod +x /usr/local/bin/ghp
```

### 5.2 인증
- gh CLI의 인증 토큰 자동 사용
- 추가 권한 요구 시 자동 프롬프트

### 5.3 대화형 모드
```bash
# 대화형 프로젝트 선택
ghp select

# 대화형 필드 업데이트
ghp item update --interactive
```

### 5.4 출력 형식
- 테이블 (기본)
- JSON
- CSV
- YAML

## 6. 성공 지표

### 6.1 정량적 지표
- 월간 활성 사용자 1,000명 이상
- GitHub Marketplace 별점 4.5 이상
- 주요 기능 응답 시간 2초 이내
- 버그 리포트 대비 해결률 90% 이상

### 6.2 정성적 지표
- 사용자 피드백 긍정률 85% 이상
- 커뮤니티 기여도 (PR, 이슈 제출)
- 주요 오픈소스 프로젝트 채택

## 7. 로드맵

### Phase 1: MVP (2개월)
- 기본 프로젝트 CRUD 작업
- 아이템 추가/조회/업데이트
- 기본 필드 관리

### Phase 2: 고급 기능 (2개월)
- 뷰 관리
- 일괄 작업
- 대화형 모드

### Phase 3: 자동화 및 분석 (2개월)
- 워크플로우 지원
- 보고서 생성
- 통계 및 분석

### Phase 4: 엔터프라이즈 (3개월)
- GitHub Enterprise Server 지원
- 고급 권한 관리
- 감사 로그

## 8. 위험 요소 및 완화 방안

### 8.1 기술적 위험
- **GitHub API 제한**: Rate limiting 대응을 위한 캐싱 및 배치 처리
- **API 변경**: 버전 관리 및 하위 호환성 유지
- **성능 이슈**: 프로파일링 및 최적화

### 8.2 사용자 채택 위험
- **학습 곡선**: 직관적인 명령어 구조 및 풍부한 예제 제공
- **기존 도구와의 경쟁**: gh CLI extension으로 제공하여 진입 장벽 낮춤

## 9. 경쟁 분석

### 9.1 기존 솔루션
- **GitHub Web UI**: 완전한 기능, 하지만 자동화 불가
- **GitHub API 직접 사용**: 유연하지만 복잡함
- **gh CLI 기본 기능**: 제한적인 프로젝트 지원

### 9.2 차별화 요소
- 독립적인 CLI 도구로 gh와 충돌 없이 사용
- 직관적이고 일관된 명령어 체계
- 강력한 일괄 처리 및 자동화 기능
- 풍부한 보고서 및 분석 기능

## 10. 보안 고려사항

- GitHub 토큰 안전한 저장 및 관리
- 민감한 프로젝트 데이터 로컬 캐싱 시 암호화
- 감사 로그 지원
- 권한 기반 접근 제어

## 11. 문서화 계획

- README.md: 빠른 시작 가이드
- 상세 사용자 매뉴얼
- API 레퍼런스
- 예제 및 튜토리얼
- 기여 가이드

## 12. 오픈소스 전략

- MIT 라이선스
- GitHub에서 개발 진행
- 이슈 템플릿 및 PR 템플릿 제공
- 정기적인 릴리스 사이클
- 커뮤니티 기여 적극 수용
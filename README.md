# AnnoWiz CLI tool 

## 개요

로컬의 데이터를 서버로 업로드하기 위한 클라이언트 CLI 이다. 

## 빌드 방법

1. Go 패키지를 설치한다. https://go.dev/dl/ 
2. `$ make build` 명령으로 빌드하면 `bin` 디렉토리에 윈도우, 리눅스, 맥용 바이너리가 생성된다. 

## 기능

기능명은 CLI의 sub command이다

### login
- `AnnoWiz File Server`의 기본 계정은 admin/1234이다. 
- 로그인을 하면 Server는 (현재는) 72시간 동안 유효한 JWT(JSON Web Token)을 발행해준다.
- 로그인을 하면 사용자를 추가하거나 데이터를 업로드 할 수 있다. 

### config
- `--server` 플래그로 서버의 URL을 추가할 수 있다. 
- `http://` 또는 `https://` 를 붙이지 않으면 `http://` 로 간주한다.
- 포트를 붙여서 쓰지 않으면 기본 포트인 `:80`와 `:443`으로 연결한다. 
  
### adduser
- `--account` 플래그로 사용자를 추가할 수 있다. 예) `--account=user1:1111`
- 중복되는 사용자를 추가하면 에러를 회신한다.

### upload
- `--file` 플래그로 특정 파일 하나를 업로드 할 수 있다.
- `--dir` 플래그로 특정 디렉토리와 그 하위 디렉토리/파일을 업로드 할 수 있다. 

## 디렉토리 업로드 프로세스

1. 임시 디렉토리에 업로드를 하려는 디렉토리를 압축하여 하나의 임시 `.tar` 파일을 생성한다.
2. 서버로 파일을 업로드 완료하면 회신을 받고 임시 `.tar` 파일을 삭제한다.
3. 서버측에서는 받은 `.tar` 파일을 지정된 디렉토리에 압축을 풀고 `.tar` 파일을 삭제한다. 

## 백로그 

이후 추가를 고려하는 기능 리스트 

### 업로드 기능

- [ ] 데이터 카테고리 구분. 이미지, 자율주행, 라이다 등등
- [ ] 데이터를 저장할 서버 디렉토리 지정
- [ ] 사용자 추가, 삭제, 비밀번호 업데이트 


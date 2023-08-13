# BLOG PROJECT

## **프로젝트명 : 리액트로 바위치기**

<br>

![리액트로-바위치기](https://github.com/choigonyok/blog-project-frontend/assets/129271363/fb779c88-d2eb-42db-92a0-7e9f9885622b)

<br>

# **개요**

나만의 테크 블로그를 직접 개발해서 사용하자!

## **1-1 개발기간**

        23.06.01 ~ 23.06.20 (20일)

## **1-2 역할**

        최윤석 (경희대 컴퓨터공학과)

* #### [instagram](https://www.instagram.com/choigonyok)

* #### [blog](https://www.choigonyok.com)

> 설계, BE, FE, 배포 담당

## **1-3 결과물**

### [www.choigonyok.com](https://www.choigonyok.com)

## **1-4 리팩토링**

        기간 : 23.08.11 ~ 
        내용 : 
        - MVC패턴 적용
        - Go directory convention 적용
        - K8S 배포
        - Code refactoring

<br>

# **개발이유**

## **쉬워서**
   
블로그 개발은 너무 쉽다고들 한다. 오히려 잘 됐다고 생각했다. 첫 프로젝트인 만큼 과분하게 어려운 목표를 정하면 쉽게 지칠 수 있을 뿐더러 종국엔 프로젝트를 완성하지 못하는 상황이 생길 수 있다.

또 쉽다는 이유로 많은 개발자들이 블로그 개발에 도전하기 때문에, 진행하면서 겪는 어려움이나 문제들을 쉽게 서칭할 수 있을 것이다.

<br>

## **클론 코딩 하지 않을거라서**

클론 코딩의 정의가 사람들마다 나뉘는 것 같다.

인프런, 유튜브 등 강의를 보면서 A-Z까지 하라는대로 따라서 만드는 것

한 레퍼런스를 정하고 어떻게 기능들을 구현하였을까 고민하면서 같은 기능을 구현하는 것

둘 중 클론코딩의 정의가 뭐가 됐던 난 둘 다 하지 않을 것이다.

직접 들이박고 해결해내며 얻는 지식의 가치나 기쁨을 많이 느껴봤기 때문이다. 물론 문제에 직면했을 때, 빠르게 답을 찾아내 문제를 해결하는 것도 개발자의 중요한 덕목 중 하나이다.

그러나 지금 개발에 막 뛰어든 나로써는 혼자 고민하고 탐구하며 기본기를 탄탄히 쌓아가는 것이 더 중요하다.

<br>

## **앞으로 유용하게 쓸 거라서**

완성도있게 개발해두면 블로그에 대한 애정이 넘쳐서, 혹은 개발한 게 아까워서라도 블로그에 좋은 내용의 글을 더 담게 될 것이다. 블로그 개발이 학습에 대한 동기부여로 작용할 수 있다.

<br>

# **기술스택**

* BE : Go, Gin
* FE : React.js
* Publishing : CSS, HTML
* DB : MySQL
* Version Management : Git / Github
* Deployment : AWS EC2
* Web Server : Nginx
* Domain Name : AWS Route53
  
<br>

# **주요기능**

### 게시글 작성 / 수정 / 삭제
* uiw 마크다운 에디터 활용

### 로그인 및 쿠키 발급
* 로그인 시 admin 쿠키 발급
* admin 쿠키가 있어야만 게시글 작성 / 수정 / 삭제 가능

### 태그 버튼으로 태그 별 게시글 확인
* 메인 페이지에서 태그 버튼을 통해 원하는 카테고리의 글을 모아서 볼 수 있음

### 태그가 같은 관련 게시글 표시
* 보고 있는 게시글의 최하단에 같은 태그의 글들을 표시

### 게시글 별 댓글 / 답글 작성
* 게시글마다 댓글 작성, 달린 댓글에 답글 작성
* 댓글 / 답글은 작성할 때 입력한 PW로 삭제 가능
* ADMIN 사용자는 분홍색으로 아이디 구분
* 댓글을 삭제하면 댓글에 달린 답글들도 같이 삭제

### Today / Total 방문자 수 집계
* visitor 쿠키를 활용한 Today 방문자 수 집계
* Total 방문자는 누적, Today 방문자는 자정마다 초기화

<br>

# **관련 Posts**

* [[BLOG #1]블로그를 직접 만들어보자](https://choigonyok.com/post/1)

* [[BLOG #2]URL로 이미지 GET하기](https://choigonyok.com/post/8)

* [[BLOG #3]태그로 RELATED POST 관리하기](https://choigonyok.com/post/9)

* [[BLOG #4]여러 파일을 한 번에 업로드 하는 법](https://choigonyok.com/post/10)

* [[BLOG #5]댓글/답글 기능 구현하기](https://choigonyok.com/post/13)
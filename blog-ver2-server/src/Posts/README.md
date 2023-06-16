## **리액트로 바위치기**
---
<br/>

    이전 프로젝트에서 HTML과 CSS로 UI/UX만 간단히 만들고,
    DB CRUD부터 HTML 렌더링, 라우팅, 리다이렉션 등 모든 작업을 서버에서 관리하는 아키텍처로 개발
  
<br/>

BE / FE의 역할에 대해 정확히 알지 못했고, 첫 번째 프로젝트를 진행하면서 프론트엔드 학습의 필요성을 느꼈음.

<br/>

React.js를 학습하면서 동시에 바로 개발을 진행하여, 기존 개발되어있던 블로그의 기능들을 FE/BE로 분리하여 재개발하는 블로그 프로젝트 Ver.2 를 기획하게 됨.

<br/>

### **기술스택**

* BE : Go(language), Gin(web framework)
* FE : React.js(JS library), JS, HTML, CSS
* DB : MySQL
* Version Management : Git / Github

<br/>

### **학습목표**

* RESTful API 적용
* React 적용
* FE / BE 역할 이해





(6/6)
CORS 문제.
-> gin의 CORS 미들웨어 추가 (허용할 오리진 설정)

DB 읽어오기 문제.
-> SCAN()은 레퍼런스로 읽어와야함

JSON형식으로 데이터 송수신
ShouldBindJSON / Marshall
content type 헤더설정을 해줘야함 (application/json)->JSON형식으로 전송하기때문에

Reader로 요청본문 읽고, Writer로 응답본문 전송

useEffect() 사용

* 다른 언어로 작성된 클라이언트와 서버가 통신하는 걸 구현하고 직접 눈으로 확인하니까 너무 신기하고 재미있음
* 클론코딩이 어플리케이션을 분석해서 기능을 똑같이 구현하는게 아니라 누가 하는걸 보고 따라하면서 코딩하는 의미라면, 클론코딩 하지 않고 내가 직접 공부하고 생각하면서 구현하니까 느끼고 배우는 것이 많은 것 같음
* 1차 프로젝트에서 배웠던 Go는 겉핥기 수준이었고 이번 2차 프로젝트를 통해 Go에 대한 감은 더 확실히 잡은 것 같음
* 1차 프로젝트 개발 1주만에 끝나서 프론트/백 분리 개발도 쉽게 생각했는데, 생각보다 더 오래 걸릴 것 같음 그래도 마감기한은 1달을 넘기지 않도록 7월 중순 안에 마무리 해보자! 이후에는 디자이너와 협업을 통해 더 완성도 높은 어플리케이션을 개발해야지

+ useEffect에 대한 공부와, Git branch에 대한 공부 필요함. 깃브랜치를 쓰면 아마 서버/클라이언트 분리할 수 있지 않을까?

* 리액트만으로 SPA 블로그를 충분히 개발할 수 있다고 하는 것 같다. 그렇지만 내 목적은 프런트를 깊이 아는 것이 아니라 전체적인 웹 개발의 흐름을 잡고, 프론트보단 백에 더 초점이 맞춰져있어서, 굳이 서버/클라이언트 아키텍처로 개발하기 선택했다. 프런트는 앞으로 개인 프로젝트들을 해나갈 때 UI/UX적으로나 프런트적으로 불편을 느끼지 않을 정도로의 지식만 갖추기로 했다. 물론 그것 조차도 금방 해낼 수 있는 건 아니겠지만!


(6/7)

<h1>{responseData[0].Title}</h1> -> 인덱스를 꼭 써줘야함

DB 컬럼 추가 쿼리 -> ALTER TABLE `테이블명` ADD `컬럼명` 자료형
DB 레코드 추가 쿼리 -> 
 INSERT INTO [TABLE_NAME] (COL_NAM1, COL_NAME2...)
 VALUES (INPUT_VALUE1, INPUT_VALUE2...)

**내가 필요했던건 위가 아니라 아래;;**

DB 레코드 수정 쿼리 -> 
 UPDATE [TABLE_NAME] SET [COL_NAME1] = [VALUE1], ...... WHERE [CONDITION];

--------
 이미지 전송문제.
 ->리액트에서는 이미지를 경로만 쓰면 이미지를 출력하는 기능을 지원 안함
 하나하나 이미지 파일별로 import 해야하는데 그럼 동적으로 저장은 할 수 있어도 동적으로 이미지를 불러오는 게 불가능함
 ->글 하나 쓸때마다 웹서버 들어가서 이미지 import해주고 재배포 할거야? 말도 안됨
 파일째로 import하는 것도 불가능 그럼 방법은?

 1. 이미지를 base64로 변환해서 db에 저장하기 -> 데이터 사이즈가 기존 이미지 크기보다 커짐 (비효율적), base64로 인코딩,디코딩하는데 오버헤드가 너무 큼
 2. 이미지를 프런트엔드에서 저장/관리하기 -> 보안 위험 + 어차피 동적 import 문제는 해결되지 않음
 3. URL을 통한 파일 업로드/다운로드!!!!!!!!
   -> 가능한 이유? import는 컴파일 타임에 이미지를 가져오는데 url은 url이 요청될 때(런타임) 동적으로 이미지를 가져오기 때문에 import 없이도 이미지 출력이 가능하다!!!!!!



처음엔 그냥 어차피 같은 인스턴스 안에 프런트/백이 있으니까 같은 디렉토리를 공유하는 셈이고, 절대경로를 통해서 그대로 옮겨오면 안되나? 싶었음
근데 새롭게 알게된 지식....
서버는 컴퓨터다!라는 말을 이해만하고 적용을 못했다. 백과 프런트, DB를 다른 인스턴스에 분리해서 운영하는게 일반적 -> 보안을 위해서
그럼 인스턴스가 다르니까 프런트에서 백에 절대경로로 접근하는게 불가능! 다른 컴퓨터니까 당연히!
DB/백/프런트가 같은 인스턴스에서 포트만 다르게 쓰는 게 일반적인 건 줄 알았음

URL 업다운로드와 경로 업다운로드의 차이?
서버가 클라이언트로 파일을 업로드, 클라이언트는 해당 파일을 다운로드
클라이언트가 서버로 파일을 업로드, 서버는 클라이언트가 전송한 파일을 다운도르
이 개념인 줄 알았는데, 그게 아니라
업로드 : 클라이언트 -> 서버
다운로드 : 서버 -> 클라이언트
이 개념이 더 일반적으로 쓰이는 것 같다. 사용자 관점에서!

물론 서버에서 서버로 업로드 다운로드 다 할 수도 있는데, 굳이..?
**파일을 업로드하면 파일은 어디로 가는거야??**
처음엔 url로 파일을 업로드하면 파일이 웹 브라우저에 저장되는건가? 생각했음.
나는 이전 개념으로 알고있었기 때문에,
클라이언트----(업로드)--->>>(미상의 어떤 매개체)---(다운로드)--->>>서버
서버----(업로드)--->>>(미상의 어떤 매개체)---(다운로드)--->>>클라이언트
로 생각했고, 그 매개체가 브라우저인가, 아니면 뭐 OS 벤더에서 파일 업다운로드용 서버를 제공해주는 건가 했던 것임

URL은 리소스의 위치를 나타내는 주소! URL은 GET요청!! 이미지 파일등을 요청할 때 굳이 fetch 등으로 get method를 요청하지 않아도 url만 입력하면 요청이 됨. 대신 브라우저 url에 해당 요청이 다 보여서 보안적으로 좋진 않음. 간단한 파일 다운로드 정도?

url은 같은 네트워크 상에 있으면 다른 인스턴스여도 파일 다운로드가 가능함. url은 기본적으로 ip주소를 DNS를 통해 문자열로 변환한 것이기 때문! 결국 같은 네트워크 안에 있는 인스턴스면 스위치의 MAC table과 ARP, DNS를 통해서 서로 통신할 수 있기 때문에, GET요청을 하는 해당 인스턴스의 ip주소로 요청을 보내는 것과 마찬가지!! 
그래서 URL을 통해 같은 네트워크 안의 누구든 url을 통해 파일을 다운로드(GET요청)할 수 있음!!
물론 다운로드하려면 다운로드하는 파일이 있는 서버가 writer pointer로 응답을 해줘야겠지!! 안해주면 암만 url 요청 찍찍 보내도 아무것도 안옴
근데 백과 프런트를 각기 다른 인스턴스에 배포하면 같은 네트워크로 설정을 어떻게 하지? 네트워크 수업 내용을 적용하려면 스위치가 있어야하는데, AWS는 서버 가상화인데!? 가상 스위치가 있나!? ->AWS의 VPC 이용하면 됨!!

그럼 업로드는!? 업로드는 이미 구현해봤다!!!
첫번째 프로젝트에서 사용자가 업로드한 이미지파일을 multipart/form-data 형식으로 POST 요청받았다. 이건 해봐서 쉽게 구현할 수 있을 듯 하다.
아직 admin page는 구현하지 않아서 아마 하드코딩 데이터로 필요한 부분들을 다 구현한 뒤에 admin page를 마지막에 개발할 것 같다!

-> 이번 프로젝트엔 백, 프런트를 나누긴 했지만 인스턴스 하나에 개발을 했는데 다음 프로젝트때는 도커를 공부해서 각 서버를 가상화해봐야겠다.

마지막으로 

추가문제
1. 요청 url을 localhost:8080/로 하면 안됨 ->http://localhost:8080/로 해야함
2. 서버에서 파일을 open할 때 /Images/img9.jpeg 이나, ./Images/img9.jpeg로 하면 안됨 -> src/Images/img9.jpeg로 해야함 무슨 기준인지....
3. 에러처리를 위한 상태코드 설정을 보는게 디버깅에 도움이 많이 되었음 상태코드 500 = http.StatusInternalServerError

----리액트에서 URL 요청----

        <img className="postcard-image" alt="img" src="http://localhost:8080/Images/img9.jpeg"/>
 
---서버에서 GET요청 응답 구현----
        eg.GET("/", func (c *gin.Context){
                fmt.Println("접근!!!")
                file, err := os.Open("src/Images/img9.jpeg")
                if err != nil {
                    // 파일 열기에 실패한 경우 에러 처리
                    http.Error(c.Writer, "파일 오픈 실패", http.StatusInternalServerError)
                    return
                }
                defer file.Close()
                _, err = io.Copy(c.Writer, file)
                if err != nil {
                // 파일 전송에 실패한 경우 에러 처리
                http.Error(c.Writer, "Failed to send file", http.StatusInternalServerError)
                return
                }
            })


// 파일 정보 가져오기
fileInfo, err := file.Stat()
if err != nil {
  // 파일 정보 가져오기에 실패한 경우 에러 처리
  http.Error(c.Writer, "파일 정보 오픈 실패", http.StatusInternalServerError)
  return
}
c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name())
  c.Writer.Header().Set("Content-Type", "application/octet-stream")
  c.Writer.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
//이부분은 없어도 되는데 헤더에 그냥 정보 넣어주기용으로 넣어도 됨

-----동적으로 수정하기------

    eg.GET("/Images/:imgname", func (c *gin.Context){
            imgname := c.Param("imgname")
            file, err := os.Open("src/Images/"+imgname)
            ...
    })

//이렇게 수정해주면 /Images/ 뒤에 나오는 값을 파라미터로 받을 수 있고, 클라이언트는 원하는 파일을 동적으로 다운로드 받을 수 있다!!!

--------------------

이미지 (6/6-1)
DB에 저장되는 imgpath를 변경해서 DB의 imgpath를 읽어오면 바로 URL 요청을 할 수 있도록 수정했다
나중에 어드민 페이지에서 포스트 작성 후 DB로 POST 요청을 보낼 때 imgpath를 
http://localhost:8080/Images/{imagename}.jpeg
형식으로 저장하는 걸 구현하면 될 것 같다.

동영상 (6/6-1)
영상처럼 태그를 클릭하면 해당 태그가 있는 post의 정보를 서버에 요청
서버는 DB에서 data를 긁어와서 json형식으로 응답
이미지는 DB의 imgpath에 저장된 URL로 자동 GET 요청을 보내고
서버는 서버에 저장된 image 파일을 응답
하는 기능을 구현했다.

동영상 (6/6-2)
태그가 있는 post들만 사용자에게 보이도록 구현했다.

이어서는 처음엔 모든 작성된 글들이 최신 게시 순으로 전부 보이게 구현하면 될 것 같다.

+시험기간이라 깃 브랜치까지 동시에 공부하진 못하고 있어서 아직 GO파일을 커밋하지 못하고 있는데, 날아갈까봐 불안불안하다
시험 끝나면 깃 브랜치부터 어서 공부해야지
**+dev-blog랑 dev-blog ver.2 랑 이름 앞 부분이 같아서 dev-blog 태그 선택하면 dev-blog ver.2까지 같이 표시됨. 해결해야함**



(6/7)

<h1>{responseData[0].Title}</h1> -> 인덱스를 꼭 써줘야함

DB 컬럼 추가 쿼리 -> ALTER TABLE `테이블명` ADD `컬럼명` 자료형
DB 레코드 추가 쿼리 -> 
 INSERT INTO [TABLE_NAME] (COL_NAM1, COL_NAME2...)
 VALUES (INPUT_VALUE1, INPUT_VALUE2...)

**내가 필요했던건 위가 아니라 아래;;**

DB 레코드 수정 쿼리 -> 
 UPDATE [TABLE_NAME] SET [COL_NAME1] = [VALUE1], ...... WHERE [CONDITION];

--------
 이미지 전송문제.
 ->리액트에서는 이미지를 경로만 쓰면 이미지를 출력하는 기능을 지원 안함
 하나하나 이미지 파일별로 import 해야하는데 그럼 동적으로 저장은 할 수 있어도 동적으로 이미지를 불러오는 게 불가능함
 ->글 하나 쓸때마다 웹서버 들어가서 이미지 import해주고 재배포 할거야? 말도 안됨
 파일째로 import하는 것도 불가능 그럼 방법은?

 1. 이미지를 base64로 변환해서 db에 저장하기 -> 데이터 사이즈가 기존 이미지 크기보다 커짐 (비효율적), base64로 인코딩,디코딩하는데 오버헤드가 너무 큼
 2. 이미지를 프런트엔드에서 저장/관리하기 -> 보안 위험 + 어차피 동적 import 문제는 해결되지 않음
 3. URL을 통한 파일 업로드/다운로드!!!!!!!!
   -> 가능한 이유? import는 컴파일 타임에 이미지를 가져오는데 url은 url이 요청될 때(런타임) 동적으로 이미지를 가져오기 때문에 import 없이도 이미지 출력이 가능하다!!!!!!



처음엔 그냥 어차피 같은 인스턴스 안에 프런트/백이 있으니까 같은 디렉토리를 공유하는 셈이고, 절대경로를 통해서 그대로 옮겨오면 안되나? 싶었음
근데 새롭게 알게된 지식....
서버는 컴퓨터다!라는 말을 이해만하고 적용을 못했다. 백과 프런트, DB를 다른 인스턴스에 분리해서 운영하는게 일반적 -> 보안을 위해서
그럼 인스턴스가 다르니까 프런트에서 백에 절대경로로 접근하는게 불가능! 다른 컴퓨터니까 당연히!
DB/백/프런트가 같은 인스턴스에서 포트만 다르게 쓰는 게 일반적인 건 줄 알았음

URL 업다운로드와 경로 업다운로드의 차이?
서버가 클라이언트로 파일을 업로드, 클라이언트는 해당 파일을 다운로드
클라이언트가 서버로 파일을 업로드, 서버는 클라이언트가 전송한 파일을 다운도르
이 개념인 줄 알았는데, 그게 아니라
업로드 : 클라이언트 -> 서버
다운로드 : 서버 -> 클라이언트
이 개념이 더 일반적으로 쓰이는 것 같다. 사용자 관점에서!

물론 서버에서 서버로 업로드 다운로드 다 할 수도 있는데, 굳이..?
**파일을 업로드하면 파일은 어디로 가는거야??**
처음엔 url로 파일을 업로드하면 파일이 웹 브라우저에 저장되는건가? 생각했음.
나는 이전 개념으로 알고있었기 때문에,
클라이언트----(업로드)--->>>(미상의 어떤 매개체)---(다운로드)--->>>서버
서버----(업로드)--->>>(미상의 어떤 매개체)---(다운로드)--->>>클라이언트
로 생각했고, 그 매개체가 브라우저인가, 아니면 뭐 OS 벤더에서 파일 업다운로드용 서버를 제공해주는 건가 했던 것임

URL은 리소스의 위치를 나타내는 주소! URL은 GET요청!! 이미지 파일등을 요청할 때 굳이 fetch 등으로 get method를 요청하지 않아도 url만 입력하면 요청이 됨. 대신 브라우저 url에 해당 요청이 다 보여서 보안적으로 좋진 않음. 간단한 파일 다운로드 정도?

url은 같은 네트워크 상에 있으면 다른 인스턴스여도 파일 다운로드가 가능함. url은 기본적으로 ip주소를 DNS를 통해 문자열로 변환한 것이기 때문! 결국 같은 네트워크 안에 있는 인스턴스면 스위치의 MAC table과 ARP, DNS를 통해서 서로 통신할 수 있기 때문에, GET요청을 하는 해당 인스턴스의 ip주소로 요청을 보내는 것과 마찬가지!! 
그래서 URL을 통해 같은 네트워크 안의 누구든 url을 통해 파일을 다운로드(GET요청)할 수 있음!!
물론 다운로드하려면 다운로드하는 파일이 있는 서버가 writer pointer로 응답을 해줘야겠지!! 안해주면 암만 url 요청 찍찍 보내도 아무것도 안옴
근데 백과 프런트를 각기 다른 인스턴스에 배포하면 같은 네트워크로 설정을 어떻게 하지? 네트워크 수업 내용을 적용하려면 스위치가 있어야하는데, AWS는 서버 가상화인데!? 가상 스위치가 있나!? ->AWS의 VPC 이용하면 됨!!

그럼 업로드는!? 업로드는 이미 구현해봤다!!!
첫번째 프로젝트에서 사용자가 업로드한 이미지파일을 multipart/form-data 형식으로 POST 요청받았다. 이건 해봐서 쉽게 구현할 수 있을 듯 하다.
아직 admin page는 구현하지 않아서 아마 하드코딩 데이터로 필요한 부분들을 다 구현한 뒤에 admin page를 마지막에 개발할 것 같다!

-> 이번 프로젝트엔 백, 프런트를 나누긴 했지만 인스턴스 하나에 개발을 했는데 다음 프로젝트때는 도커를 공부해서 각 서버를 가상화해봐야겠다.

마지막으로 

추가문제
1. 요청 url을 localhost:8080/로 하면 안됨 ->http://localhost:8080/로 해야함
2. 서버에서 파일을 open할 때 /Images/img9.jpeg 이나, ./Images/img9.jpeg로 하면 안됨 -> src/Images/img9.jpeg로 해야함 무슨 기준인지....
3. 에러처리를 위한 상태코드 설정을 보는게 디버깅에 도움이 많이 되었음 상태코드 500 = http.StatusInternalServerError

----리액트에서 URL 요청----

        <img className="postcard-image" alt="img" src="http://localhost:8080/Images/img9.jpeg"/>

---서버에서 GET요청 응답 구현----
        eg.GET("/", func (c *gin.Context){
                fmt.Println("접근!!!")
                file, err := os.Open("src/Images/img9.jpeg")
                if err != nil {
                    // 파일 열기에 실패한 경우 에러 처리
                    http.Error(c.Writer, "파일 오픈 실패", http.StatusInternalServerError)
                    return
                }
                defer file.Close()
                _, err = io.Copy(c.Writer, file)
                if err != nil {
                // 파일 전송에 실패한 경우 에러 처리
                http.Error(c.Writer, "Failed to send file", http.StatusInternalServerError)
                return
                }
            })


// 파일 정보 가져오기
fileInfo, err := file.Stat()
if err != nil {
  // 파일 정보 가져오기에 실패한 경우 에러 처리
  http.Error(c.Writer, "파일 정보 오픈 실패", http.StatusInternalServerError)
  return
}
c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name())
  c.Writer.Header().Set("Content-Type", "application/octet-stream")
  c.Writer.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
//이부분은 없어도 되는데 헤더에 그냥 정보 넣어주기용으로 넣어도 됨

-----동적으로 수정하기------

    eg.GET("/Images/:imgname", func (c *gin.Context){
            imgname := c.Param("imgname")
            file, err := os.Open("src/Images/"+imgname)
            ...
    })

//이렇게 수정해주면 /Images/ 뒤에 나오는 값을 파라미터로 받을 수 있고, 클라이언트는 원하는 파일을 동적으로 다운로드 받을 수 있다!!!

--------------------

이미지 (6/6-1)
DB에 저장되는 imgpath를 변경해서 DB의 imgpath를 읽어오면 바로 URL 요청을 할 수 있도록 수정했다
나중에 어드민 페이지에서 포스트 작성 후 DB로 POST 요청을 보낼 때 imgpath를 
http://localhost:8080/Images/{imagename}.jpeg
형식으로 저장하는 걸 구현하면 될 것 같다.

동영상 (6/6-1)
영상처럼 태그를 클릭하면 해당 태그가 있는 post의 정보를 서버에 요청
서버는 DB에서 data를 긁어와서 json형식으로 응답
이미지는 DB의 imgpath에 저장된 URL로 자동 GET 요청을 보내고
서버는 서버에 저장된 image 파일을 응답
하는 기능을 구현했다.

동영상 (6/6-2)
태그가 있는 post들만 사용자에게 보이도록 구현했다.

이어서는 처음엔 모든 작성된 글들이 최신 게시 순으로 전부 보이게 구현하면 될 것 같다.

+시험기간이라 깃 브랜치까지 동시에 공부하진 못하고 있어서 아직 GO파일을 커밋하지 못하고 있는데, 날아갈까봐 불안불안하다
시험 끝나면 깃 브랜치부터 어서 공부해야지
**+dev-blog랑 dev-blog ver.2 랑 이름 앞 부분이 같아서 dev-blog 태그 선택하면 dev-blog ver.2까지 같이 표시됨. 해결해야함**



onClick event 핸들러에 파라미터를 전달하려면,

onClick={cardClickHandler(value)}
이게 아니라!
onClick={()=>{cardClickHandler(value)}}
이렇게 화살표함수를 사용해야한다.

------------------------------------------------------------------------------------
## (6/7)
포스트 글을 읽고 맨 밑에 보고있는 글과 같은 태그를 가진 글들을 모아서 RELATED POST로 보여주는 기능을 구현하려고 했는데 글 전체를 다 불러오는 건 가능한, 현재 보고있는 글만빼는 건 잘 안됐음 리액트의 FILTER를 이용! (6/7 이미지 -1)

INCLUEDS 붙여보기 -> 안됨 (6/7 이미지 -2)
CONSOLE.LOG(POSTID); CONSOLE.LOG(JSONARRAY[0].ID); 로 콘솔 확인해보기 (6/7 이미지 -3) -> 둘다 10이 잘 나오는데 왜 필터링이 안되지??? 하다가 콘솔상에서 10의 색이 다른 걸 발견...! (6/7 이미지 -4)
그 위에 콘솔에 POST 응답받은 본문을 살펴보니 (6/7 이미지 -5) ID만 보라색! 그 때 깨달았다. ID는 INT 형식으로 전달받았다는 거 (6/7 이미지 -2) -> GO로 작성된 서버 코드의 구조체 서버에서 선언은 INT타입으로 해도 JSON으로 마샬링(인코딩)해서 전달하면 다 키-값 형태의 STRING으로 변환되는 것인 줄 알았음. 타입을 그대로 갖고오는 줄 몰랐음!!! 그래서 알고 난 이후에 (6/7 이미지 -7)처럼 ID 키의 값을 STRING으로 변환해주었더니 (6/7 이미지 -8) 같은 흰색으로 콘솔에 출력되는 걸 확인할 수 있었다. 이후에 (6/7 이미지 -9) 필터링을 해서 USESTATE를 초기화해주고, 리렌더링 되면서 CARD 컴포넌트에는 USESTATE로 초기화된 필터링 POST들의 DATA가 전달된다! 그럼 원하는 기능이 잘 구현되는 걸 볼 수 있다. (6/7 영상 -1)

처음 루트로 라우팅되면 홈페이지에서 전체 POST가 나오도록 구현했다. 원래는 렌더링하면 아무 POST도 나오지 않고 태그 버튼을 눌러야만 해당 태그에 맞는 글들이 나왔었는데, 그 이유는 (6/7 이미지 -10) 버튼 클릭 핸들러에만 JSON 형식으로 태그 데이터와 함께 POST요청을 보내도록 구현해뒀기 때문이다. (6/7 이미지 -11) 처음에 실행되는 POST요청은 USESTATE의 디폴트 값으로 요청이 되고, 이후에는 (USESTATE 값 변경으로 인한)리렌더링 여부와 관계없이 []에 들어있는 POSTDATA 값에 변화가 있을 때만 USEEFFECT가 재실행 된다->POST 요청을 보낸다. SETPOSTDATA() 가 클릭 핸들러에서만 동작했기 때문에 태그 버튼을 클릭하는 이벤트가 발생해야만 POST요청이 서버에 전송되었던 것이고, 내가 원하는 기능은 처음 접속했을 때 아무것도 누르지 않아도 전체 POST들을 볼 수 있게 하는 것이다.

USEEFFECT는 처음 컴포넌트 실행시에는 기본적으로 자동 실행되기 때문에, POSTDATA에 DEFAULT값으로 무엇이 들어있는지가 중요하다. 원하는 기능을 구현하려면 DEFAULT 값으로 모든 태그가 POSTDATA에 들어있어야한다. 근데 DEFAULT값은 USESTATE를 선언할 때 정적으로 입력해주어야하는데, 태그는 블로그의 글을 작성/삭제하는 것에 따라서 태그가 늘어날수도/줄어들수도 있잖아!? 그렇게 동적으로 바뀐 태그가 있다면 몇몇 태그는 DEFAULT 값에서 누락되게 되고 -> 결국 첫 화면에 모든 POST를 전부 보여주지 못하고 몇개의 POST는 누락될 수 있다. (6/7 이미지 -13) 그래서 POSTDATA의 DEFAULT 값으로 키 => TAGNAME, 값 => "ALL" 이라는 JSON 형식의 데이터를 할당했다. 이제 이 데이터는 POST요청으로 본문에 담겨 서버로 전송되고 서버에서 이를 잘 처리하는 로직을 구현하면 될 것이다.

GO 서버에서 TAGDATA 타입의 DATA를 선언해주었다 TAGDATA는

TYPE TAGDATA STRUCT {
TAGS STRING `JSON:"TAGNAME"`
    }
구조체이고 JSON형식의 TAGNAME키를 보면 TAGS가 값을 받아오도록 미리 선언해주었다.

(6/7 이미지 -14 안씀) (6/7 이미지 -15)

SHOULDBINDJSON으로 JSON 형태의 데이터를 변수 DATA에 저장한다. 그리고 DATA에 저장된 키(TAGNAME)의 값이 "ALL"이면 값을 "ALL"에서 빈 문자열 ""로 바꿔준다. GO는 STRING에서 ""을 0이 아니라 NIL(=NULL)로 취급한다

이 뒤에는

DB.QUERY("SELECT ID,TAG,TITLE,BODY,DATETIME,IMGPATH FROM POST WHERE TAG LIKE `%"+DATA.TAGS+"%`")
쿼리문을 통해 POST들의 데이터를 읽어온다. DATA.TAGS는 NIL 상태이기 때문에 **WHERE TAG LIKE %%**는 사실상 없는 것과 마찬가지인 상태가 된다. 물론 키값이 "ALL"이 아니라 다른 특정 태그가 들어오면, 이 쿼리를 통해 해당 태그가 들어간 POST들을 찾게된다.

이렇게 찾아온 데이터를 JSON형식으로 변환하고,

RESPONSE, ERR := JSON.MARSHAL(POSTDATA)
응답이 JSON 형식으로 간다고 헤더로 알려주고, 본문에 RESPONSE를 담아 응답하면

 C.WRITER.HEADER().SET("CONTENT-TYPE", "APPLICATION/JSON")
 C.WRITER.WRITE(RESPONSE)
(6/7 영상 2) 영상처럼 홈페이지 루트로 라우팅했을 때 전체 POST를 다 볼 수 있다. ALL 태그를 누른 것과 안누른 것을 구분하기 위해서 (6/7 이미지 -16) 홈페이지로 라우팅했을 땐 아무 태그가 안눌린 상태이기 때문에 태그버튼 위의 제목을 CHOIGONYOK으로 지정해두었고 (6/7 이미지 -17) 이후에 ALL 버튼을 누르면 제목이 ALL로 되도록 구현했다. 표시되는 POST는 ALL 태그를 누르나 안누르나 전체 POST로 동일하지만, 기본 홈페이지로 접속했을 땐 제목으로 ALL보단 CHOIGONYOK이 렌더링되는 게 더 적절한 것 같아서 이렇게 구현하게 되었다.

-----리액트 RELATED POST로 현재 보고있는 글은 빼고 보여주기------ CONSOLE.LOG(POSTID); CONSOLE.LOG(STRING(JSONARRAY[0].ID)); CONSOLE.LOG(JSONARRAY.FILTER((POST) => STRING(POST.ID) !== POSTID)); CONSOLE.LOG(JSONARRAY.FILTER(JSONARRAY => JSONARRAY.ID !== POSTID)); SETRELATEDPOSTDATA(JSONARRAY);

------------------------------------------------------------------------------------
## (6/8)
깃 커밋하다가 코드충돌이 나서 작업 몇시간 한 걸 그대로 다 날릴뻔했다.
다행히 VScode의 타임라인 기능과 기존 로컬에 있던 코드를 미리 복사해서 백업해둔 덕분에
원래 있던 작업영역은 아예 삭제하고 git에서 clone을 새롭게 만들어서 git 레포와 remote 연동하고 git pull 하고 백업코드 복사해서 다시 푸시해서 다행히 살렸음
날아간 거 알았을 때 진짜 식은땀이 났음
긴장 잘 안하는데 오랜만에 완전 식겁함

태그를 유동적으로 쓸 수 있게 동적화 작업을 함
글을 작성했는데, 글에 새로운 태그가 있으면 자동으로 태그 버튼이 생기고,
한 post만 갖고있던 태그에서 post가 삭제되면 더이상 그 태그를 가지고있는 post가 없는거니까
자동으로 그 태그도 태그 버튼에서 사라지도록 구현함
그러다보니 DB에서 모든 게시글의 태그를 매번 읽어와서 버튼을 생성해야하는데,
아직 마무리를 다 못해서
같은 태그의 버튼이 중복해서 여러개 생긴 상황임
중복되는 태그만 하나로 묶어주면 태그버튼 동적화 작업은 끝!
나중에 admin 페이지에서도 따로 태그버튼을 관리할 필요는 없어졌다~

GO에서 slice 중복 제거하는 법
->map을 이용하면 간편!!

```go
ret := []TagButtonData{}
		m := make(map[TagButtonData]int)
		for i, v := range posttagdata {
			if _, ok := m[v]; !ok{
				m[v] = i
				ret = append(ret, v)
			}
		}
```

태그버튼 중복 제거도 구현 완료!



* GO DB 작성할 때 스트링을 건드려야하는 작업이면
"대신 `백틱 쓰는 거 잊지말기!


* GO에는 상대경로 개념이 없다.
os.create()에 절대경로 넣어야함

로직
클라이언트에서 post data 따로, 이미지파일 따로 post 파라미터 통해서 두번 받아오고
이미지파일은 go에서 create하고 받은 파일 open해서
open한 파일 create한 파일에 io.copy로 복사
둘다 defer로 닫아주기

문제 : 
1. json은 json형식으로, 파일은 multipart/form-data로 받아야하는데, 둘 다 한 번에 하려다가 오류
2. 리액트에서 post할 때 헤더에 컨텐트타입 명시해줘야하는데 헤더 쓰는 법 몰라서 오류
3.  const formData = new FormData();
    formData.append("file", img); 이거 몰라서 오류
4. e.target.files[0] 이거 몰라서 오류
5. post요청 두개 비동기 실행되서 오류
6. image가 하나씩만 읽힘

-> 여러개를 동시에 선택하면 e.target.value 때문에 계속 최신 상태로 초기화가 되고, 결국 맨 마지막 파일만 읽어들임
const imgHandler = (e) => {
    setIMG(e.target.files[0]);
  };

  에서

  const imgHandler = (e) => {
    setIMG(e.target.files);
  };

  로 변경

  그럼 서버에서도 여러 파일을 읽을 수 있게 구현해야함


 c.SaveUploadedFile()? 
 c.String() ?->응답 전송?


리액트에서 보낸 여러개 파일을 go로 받는법

 서버사이드

 d, err := c.MultipartForm()
			if err != nil {
				fmt.Println("111111111111111111111111111111111111")
			}
			fff := d.File["file"]
			for _, v := range fff {
				fmt.Println(v.Filename)
				c.SaveUploadedFile(v, "/Users/yunsuk/blog-ver2-server/src/IMAGES/"+v.Filename)
			}
			c.Writer.WriteHeader(http.StatusOK)


클라이언트 사이드

 useEffect(() => {
    if (!mounted.current) {
      mounted.current = true;  //컴포넌트 실행시에 useEffect는 실행되지 않고, []안의 값이 바뀔때만 실행되도록 함
    } else {
      const formData = new FormData(); //formdata로 만들기
      for (let i = 0; i < img.length; i++) { 
        ->img useState에 들어있는 여러개의 파일을 하나씩 읽어서 formdata에 넣는 거
        formData.append("file", img[i]);  
        -> 여기서 "file" 이 이름이 서버사이드에서 파싱할 수 있는 key가 됨!!!!!!
      }
      axios
        .post("http://localhost:8080/postdata/img", formData, {
          "Content-type": "multipart/form-data",
        })
        .then((response) => {
          // 응답 데이터 수신
          console.log("POST2 Success");
          navigate("/");
        })
        .catch((error) => {
          console.error(error);
        });
    }
  }, [unlock]);


  업로드된 파일 저장할 때

  loc, err := os.Create("/Users/yunsuk/blog-ver2-server/src/IMAGES/"+f.Filename)
	if err != nil {
		fmt.Println("FOLDER OPEN ERROR")
	}
	defer loc.Close()
  _, err = io.Copy(loc, file)
			if err != nil {
				fmt.Println("COPY ERROR")
	}

  -> 이거 대신에

  c.SaveUploadedFile(v, "/Users/yunsuk/blog-ver2-server/src/IMAGES/"+v.Filename)

  이거 쓰면됨. 이건 파일이 multipartform일 때만 가능, 즉 업로드된 파일에 쓸 수 있는 함수

  * db에 저장할 imgpath는 썸네일이미지 하나만 있으면 됨
  -> 본문에 들어가는 이미지들은 애초에 글을 쓸 때 링크를 걸어놓고 본문을 쓰고, 본문이 불러와질때 자체적으로 get요청을 해서 이미지를 렌더링하니까!
  * 근데 처음에 글 쓸때 이미지를 사용자에게 받는 건 썸네일+본문이미지 전체를 받아서 서버에 이미지를 저장해둬야함
  * 그럼 여러개의 파일을 동시에 업로드해야하는데, 여러개 이미지를 input file 태그로 받는 건 multiple 속성을 통해 가능 (코드첨부하기)
  * 근데 그 여러개 선택된 이미지파일을 클라이언트에서 정제해서 서버로 보내줘야하고 서버도 여러개의 이미지 파일을 한 번에 받아서 저장해야함
  * + 데이터 전송은 저장하기 버튼 누르면 한꺼번에 데이터가 서버로 업로드 되어야하는데, 이미지는 multiform 형식이고, 나머지는 json 형식임
  

# 서버사이드

  f, err := c.FormFile("file")
			if err != nil {
				http.Error(c.Writer, err.Error(), http.StatusBadRequest)
				return
			}
			file, err := f.Open()
			if err != nil {
				fmt.Println("IMG OPEN ERROR")
        }
			defer file.Close()
			loc, err := os.Create("../IMAGES/"+f.Filename)
			if err != nil {
				fmt.Println("FOLDER OPEN ERROR")
        }
			defer loc.Close()
			_, err = io.Copy(loc, file)
      if err != nil {
				fmt.Println("COPY ERROR")
			}
			c.Writer.WriteHeader(http.StatusOK)



  -> 


  d, _ := c.MultipartForm()
			fff := d.File["FileList"]    
      ------>>>>(왜 filelist로 적었는지 이미지 첨부하기:멀티파트이미지7)
			for _, v := range fff {
				c.SaveUploadedFile(v, "/Users/yunsuk/blog-ver2-server/src/IMAGES/"+v.Filename)
			}
      c.Writer.WriteHeader(http.StatusOK)





 d, _ := c.MultipartForm()
			fff := d.File["0"]    
      ------>>>>(왜 0으로 적었는지 이유 설명첨부하기:멀티파트이미지8)
			for _, v := range fff {
				c.SaveUploadedFile(v, "/Users/yunsuk/blog-ver2-server/src/IMAGES/"+v.Filename)
			}
      c.Writer.WriteHeader(http.StatusOK)

      
## 서버 최종코드

eg.POST("/postdata/:datatype", func (c *gin.Context){
		datatype := c.Param("datatype")
		if datatype == "post"{
			var data RecieveData
			if err := c.ShouldBindJSON(&data); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			var idnum int
			r, _ := db.Query("SELECT id FROM post order by id desc limit 1")
			for r.Next(){
				r.Scan(&idnum)
			}
			idnum += 1
			fmt.Println(strconv.Itoa(idnum))
			_, err = db.Query(`INSERT INTO post (id, tag, datetime, title, body) values (`+strconv.Itoa(idnum)+`, "`+data.Tag+`","`+data.Datetime+`","`+data.Title+`","`+data.Body+`")`)
			if err != nil {
				fmt.Println("DB POST ADD ERROR")
			}

		} else if datatype == "img" {
			d, err := c.MultipartForm()
			if err != nil {
				fmt.Println("111111111111111111111111111111111111")
			}

			//방금 만들어진 post의 id 확인
			r, _ := db.Query("SELECT id FROM post order by id desc limit 1")
			var idnum int
			for r.Next(){
				r.Scan(&idnum)
			}

			fff := d.File["file"]
			for _, v := range fff {
				fmt.Println(v.Filename)
				c.SaveUploadedFile(v, "/Users/yunsuk/blog-ver2-server/src/IMAGES/"+v.Filename)
			}
			
			//DB에는 대표이미지(썸네일)만 저장
			_ , err = db.Query(`UPDATE post SET imgpath = "http://localhost:8080/IMAGES/`+fff[0].Filename+`" where id = `+strconv.Itoa(idnum))
			if err != nil {
				fmt.Println("POST IMAGE ATTACH ERROR")
			}
			c.Writer.WriteHeader(http.StatusOK)
		}
	})











# 클라이언트 사이드


const postHandler = (e) => {
        **-----------> img를 한 postdata에 묶어서 보냄**
    const postdata = {
      title: titleText,
      tag: tagText,
      datetime: dateText,
      img: img,
      body: bodyText,
      };
    console.log(postdata);
    axios
      .post("http://localhost:8080/post", postdata) 
      **->JSON형태라 이미지파일 전송 오류생김**
      .then((response) => {
        // 응답 데이터 수신
        console.log("POST Success")
        navigate("/");
      })
      .catch((error) => {
        console.error(error);
      });
  };



  const postHandler = (e) => {
    const postdata = [
      {
        title: titleText,
        tag: tagText,
        datetime: dateText,
        img: img,
        body: bodyText,
      },
    ];
    console.log(postdata);
    axios
      .post("http://localhost:8080/post", postdata,{headers: { "Content-type": "multipart/form-data" }})
      **----->헤더 붙이니까 다른 데이터가 전송이 안됨**
      .then((response) => {
        // 응답 데이터 수신
        console.log("POST Success");
        navigate("/");
      })
      .catch((error) => {
        console.error(error);
      });
  };






const postHandler = (e) => {
    const postdata =
      {
        title: titleText,
        tag: tagText,
        datetime: dateText,
        // img: img, **--> 이미지 빼고 헤더 빼고 보내니까 잘 전송됨**
        body: bodyText,
      };
    console.log(postdata);
    axios
      .post("http://localhost:8080/post", postdata)
      .then((response) => {
        // 응답 데이터 수신
        console.log("POST Success");
        navigate("/");
      })
      .catch((error) => {
        console.error(error);
      });
  };


**-->>그래서 이미지파일 따로, json 따로 각자 헤더설정해서 두개의post요청을 보냄**
**--> 근데 막 저장이 잘 됐다가 안됐다가 함**
  const postHandler = (e) => {
    const postdata =
      {
        title: titleText,
        tag: tagText,
        datetime: dateText,
        body: bodyText,
      };
    console.log(postdata);
    axios
      .post("http://localhost:8080/post/post", postdata)
      .then((response) => {
        // 응답 데이터 수신
        console.log("POST1 Success");
      })
      .catch((error) => {
        console.error(error);
      });
      axios
      .post("http://localhost:8080/post/img", img, {"Content-type" : "multi/form-data"})
      .then((response) => {
        // 응답 데이터 수신
        console.log("POST2 Success");
        navigate("/");
      })
      .catch((error) => {
        console.error(error);
      });

      
  };



**---->>> post요청이 동기적으로 작동해서 이미지가 먼저 요청되면 오류가 나는 거였음(DB쿼리문 어떻게 짰는지 설명)**
**그래서 unlock useState를 하나 만들고 더 늦게 요청해야하는 이미지 POST 요청을 useEffect에 넣음, 첫번째 요청이 잘 끝나고 imgpath가 null인 상태로 db에 저장된 후 응답받으면 unlock을 set하고 그럼 unlock이 바뀌니까 두번째 Post 요청이 실행**

  useEffect(() => {
      axios
        .post("http://localhost:8080/postdata/img", formData, {
          "Content-type": "multipart/form-data",
        })
        .then((response) => {
          // 응답 데이터 수신
          console.log("POST2 Success");
          navigate("/");
        })
        .catch((error) => {
          console.error(error);
        });
  }, [unlock]);

  const postHandler = (e) => {
    const postdata = {
      title: titleText,
      tag: tagText,
      datetime: dateText,
      body: bodyText,
    };
    axios
      .post("http://localhost:8080/postdata/post", postdata)
      .then((response) => {
        setUnLock(true);
      })
      .catch((error) => {
        console.error(error);
      });
  };



**근데 useEffect는 처음 컴포넌트 실행할 때도 한 번 실행돼서, 오류가 생기길래 컴포넌트가 마운트 되도 useEFFect의 []값이 바뀌기 전까진 실행되지 않도록 구현**
   useEffect(() => {
    if (!mounted.current) {  **---->이부분**
      mounted.current = true;
    } else {
      axios
        .post("http://localhost:8080/postdata/img", formData, {
          "Content-type": "multipart/form-data",
        })
        .then((response) => {
          // 응답 데이터 수신
          console.log("POST2 Success");
          navigate("/");
        })
        .catch((error) => {
          console.error(error);
        });
    }
  }, [unlock]);



  
  **이미지가 대표이미지 하나일 땐 모든 게 잘 작동함**
  **이렇게 잘 처리했는데, 아직 여러 파일을 전송받아야하는 문제는 해결 안됨 그래서 formdata라는 걸 활용**
   useEffect(() => {
    if (!mounted.current) {
      mounted.current = true;
    } else {
      const formData = new FormData();
      formData.append("file", img);
      axios
        .post("http://localhost:8080/postdata/img", formData, {
          "Content-type": "multipart/form-data",
        })
        .then((response) => {
          // 응답 데이터 수신
          console.log("POST2 Success");
          navigate("/");
        })
        .catch((error) => {
          console.error(error);
        });
    }
  }, [unlock]);
  

  **작동안됨. img는 여러개의 키-값을 가진 GO로 따지면 map이었고,**
  **여러 데이터를 보내려면 각 키에 접근해서 값을 formData에 넣어줘야함**
  **최종코드**


  useEffect(() => {
    if (!mounted.current) {
      mounted.current = true;
    } else {
      const formData = new FormData();
      for (let i = 0; i < img.length; i++) {
        formData.append("file", img[i]);
      }
      axios
        .post("http://localhost:8080/postdata/img", formData, {
          "Content-type": "multipart/form-data",
        })
        .then((response) => {
          // 응답 데이터 수신
          console.log("POST2 Success");
          navigate("/");
        })
        .catch((error) => {
          console.error(error);
        });
    }
  }, [unlock]);

  const postHandler = (e) => {
    const postdata = {
      title: titleText,
      tag: tagText,
      datetime: dateText,
      body: bodyText,
    };
    console.log(postdata);
    axios
      .post("http://localhost:8080/postdata/post", postdata)
      .then((response) => {
        setUnLock(true);
      })
      .catch((error) => {
        console.error(error);
      });
  };


  -> 결론. 게시글 작성 시 여러 이미지 데이터를 한 번에 넣어서 서버에 전달하고 서버는 첫번째 이미지를 썸네일용으로 쓰고
  나머지 이미지는 서버에 저장되어있다가 본문에서 URL로 GET요청을 하면 때마다 응답해서 본문 이미지도 구현할 수 있게 됨
  (멀티파트이미지1~6) 설명
















**->delete용 전체 게시글 data 불러오기 코드 재사용**
  eg.GET("/post/:postid", func (c *gin.Context){
		postid := c.Param("postid")
		var r2 *sql.Rows
		if postid == "all" {
			r2, err = db.Query("SELECT id, tag,title,body,datetime,imgpath FROM post")
			if err != nil {
				log.Fatalln("ID FINDING ERROR!!")
			}
	} else {
		r2, err = db.Query("SELECT id, tag,title,body,datetime,imgpath FROM post where id = "+postid)
		if err != nil {
			log.Fatalln("ID FINDING ERROR!!")
		}
	}
		
		postdata := []SendData{}
		var temp SendData
		for r2.Next(){
			r2.Scan(&temp.Id,&temp.Tag,&temp.Title,&temp.Body,&temp.Datetime, &temp.ImagePath)
			temp.Datetime = strings.TrimSuffix(temp.Datetime, " 00:00:00")
			temp.Datetime = strings.ReplaceAll(temp.Datetime,"-", "/")
			temp.Tag = strings.ToUpper(temp.Tag)
			temp.Title = strings.ToUpper(temp.Title)
			temp.Body = strings.ToUpper(temp.Body)
			temp.ImagePath = strings.ToUpper(temp.ImagePath)
			
			postdata = append(postdata, temp)
			
		}
		
		// JSON 응답 생성
		response, err := json.Marshal(postdata)
		if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
		}
		
		// JSON 헤더 설정 및 응답 전송
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Write(response)
	

	})
// describe block: first param is tested target, second param is anonymous callback function to test
describe('Blog', () => {
  beforeEach(()=> { 
    cy.visit('/')
  })
  // it block: individual testm, first param is title of test, second param is anonymous callback function to test
  it('passes', () => {
    cy.get('input').should("have.length", 5).then(($btn)=>{
      cy.wrap($btn).click({multiple: true}) 
    })
    cy.get('[cypress="tag"]').should("not.have.value", 'disabled').then(($btn2)=>{
      cy.wrap($btn2).click({multiple: true})
    })
  })
})

/* 
Cypress는 비동기적으로 작동해서 일반적인 자바스크립트 문법처럼 
const button = cy.get("button")
button.click() 이런 식의 문법 사용이 불가능하다.
대신 yield를 사용해서 큐에 하나씩 커맨드를 집어넣어두고 

DOM에 임의 속성을 추가하고 그걸 Cypress에서 사용할 수 있음
CSS이름이나 속성같은 건 쉽게 바뀌고 예기치 않은 테스트 실패가 발생할 수 있기 때문
ex)
<div className="container"></div>
를
<div cypress="container" className="container"></div>
이런식으로 바꿔서 사용하고, cy.js 파일에는
cy.get('[cypress="tag"]').should("not.have.value", 'disabled').then(($btn2)=>{
  cy.wrap($btn2).click({multiple: true})
})
이런식으로 적용할 수 있음

cypress/support/commands.js에 커맨드를 정의하면 cy.COMMAND로 사용자 정의 커맨드를 사용할 수 있음

Cypress는
./node_modules/.bin/cypress run
커맨드를 통해서 터미널에서 정의된 cy.js파일대로 테스트를 수행할 수 있음
CI/CD 파이프라인에 E2E 테스트를 포함시킬 때 이 명령어를 사용하면 유용할 듯
*/
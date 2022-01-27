# gotrade
GO 언어 공부용 프로젝트
해당 프로젝트는 공부용 프로젝트로 되도록 외부 lib를 사용하지 않고 구현해 보는것의 의의를 둠.

# 작업하면서 느낀점
## Golang답게 하려면 Channel, go rutine을 최대한 많이...

# 구현된 내용 
## Golang 언어 공부용 프로젝트로 Go 언어의 특성을 살렸다고 볼수는 없을듯...

1. client
 - IClient
 - UpbitRestClient (IClient 구현체)
2. coin
 - Candle 객체
 - Coin 객체
3. config
 - 추후 ENV이용할 계획이며 현재는 테스트 버전인 관계로 설정값 하드코딩
4. database
 - IDatabase
 - FileDatabase (IDatabase 구현체로 내용은 그냥 파일 읽기/쓰기)
5. market
 - IMarket
 - MarketUpbit (IMarket 구현체)
   : Market 객체가 client 객체를 소유하고 Monitor객체에서 마켓을 통해서 정보를 가져온다.
6. monitor
 - Monitor 객체는 각 Market객체에서 전달될 코인시세를 취합하여 전략(IStrategy)검증 작업을 수행
7. pubusb
 - Go lang channel기반 publish/subscribe (기본구조)
8. strategy
 - IStrategy
 - 이평, 상승추세, 급등 샘플 (테스트 소스이므로 손절가, 매입가 거의 3% 정도로 넣음. 실제 매매에 도움 안될듯..)
9. tree
 - FileDatabase 객체에서 사용할 cache용도
 - Simple btree
10. txlog
 - ITxTransfer
 - TransferConsole (콘솔로 로그 전송)
 - TransferFileDB (4번 FileDatabase로 로그 전동)
11. util
 - TimeBase Expire Cache 간단버전
 - 환형 큐 간단버전
 - 기타 유틸함수들..
12. marketDirector
 - MarketDirector 프로그램의 메인이 되는 객체로 마켓들과 모니터, 전략등을 관리, 디렉팅

###
https://www.dumels.com/ 에서 UML Class diagram 확인 가능

# 초기기능

# 채팅방 참여하기

- 일단 서버쪽에 미리 채팅방을 만들어놔야 한다.
-

# 내가 참여중인 채팅방 목록 불러오기, 검색

# 나가있다가 다시 들어왔을 때 안 읽은 메세지만 구분해서 보여주는 기능

# 최대 커넥션 개수

# 기본 메세지 읽기

# 기본 메세지 보내기 (일대일, 단체)

# 채팅방 나가기

# https://remix.run/docs/en/main/start/tutorial

# 스트림이란?

https://docs.nats.io/nats-concepts/jetstream/streams

메세지 저장소!
Streams consume normal NATS subjects,
any message published on those subjects will be captured in the defined storage system.

그렇지만 리소스를 가장 많이 잡아먹기 때문에, 아주 오래도록 쌓여있는 메세지 같은 건
규칙을 설정해서 자동으로 지워버리게끔 할 수도 있다.

# NATS의 장점: RPC 문제, DB 동기화 오류 (RDB와 레디스 ELK ..)

https://nats.io/blog/mongodb-nats-connector/
손쉽게 해결!

Change Data Capture
One solution to this problem is Change Data Capture (CDC). Whenever there is a data change event in MongoDB (such as an insertion, an update or a deletion), we can capture it and write it to NATS JetStream. From there, any number of consumers can subscribe to the stream, process the change events, and react accordingly. Each consumer processes the data and sends an acknowledged (ACK) signal to NATS, to confirm the correct handling. In case of retryable errors, consumers can send a not-acknowledged (NAK) signal, prompting NATS to redeliver the message. Although this introduces eventual consistency, as the post may be found on MongoDB but not yet on Redis and Elastic, it eventually reaches all systems.

Now the question is, how do we capture data changes on MongoDB and publish them on NATS JetStream? A vital piece of the puzzle is still missing.

그럼 mysql, redis 와의 연동은 지원해주나?

https://debezium.io/
그리고 NATS CLI 툴로 가능하다고 하네??

사용자가 마지막으로 접근한 시각 / 이벤트 발행시각 비교해주기
이벤트는 정적으로 저장 (저장기간은 당일 밤 12시까지만. 이후엔 삭제)

=> 기존 폴링시스템에서 시퀀스넘버 대신에 시각으로만 하면안되나?
... //쿠키에다가 마지막 접근시각 넣어줄 수 있자나..

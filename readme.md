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

# NATS 장점

요소 간의 통신을 간소화한다.
확장하기 쉽고,
NATS 서버쪾이 바뀐다고 해서 클라이언트를 매번 재시작할 이유가 없다
적응성이 높음
클라우드 온프레미스 플랫폼을 혼용한 경우에도 통신 가능

# 사용사례

Cloud Messaging

Services (microservices, service mesh)

Event/Data Streaming (observability, analytics, ML/AI)

Command and Control

IoT and Edge

Telemetry / Sensor Data / Command and Control

Augmenting or Replacing Legacy Messaging Systems

# 내장된 개발패턴

Streams and Services through built-in publish/subscribe, request-reply, and load-balanced queue subscriber patterns. Dynamic request permissioning and request subject obfuscation is supported.

# 통신에서 보장하는 것

At most once, at least once, and exactly once is available in JetStream.

# 멀티테넌시 지원해줌. 분산된 보안

NATS supports true multi-tenancy and decentralized security through accounts and defining shared streams and services.

# 보안관련

NATS supports TLS, NATS credentials, NKEYS (NATS ED25519 keys), username and password, or simple token.

# 메세지 보존과 지속성

Supports memory, file, and database persistence.
Messages can be replayed by time, count, or sequence number,
and durable subscriptions are supported.
With NATS streaming, scripts can archive old log segments to cold storage.

메세지를 보관하는 기능을 제공하기 때문에 필요할 때 꺼내볼 수 있는것 같다. (카톡 안읽은 메세지처럼..)

# 페일오버

Core NATS supports full mesh clustering with self-healing features to provide high availability to clients. NATS streaming has warm failover backup servers with two modes (FT and full clustering). JetStream supports horizontal scalability with built-in mirroring.

# 배포

크고 유연한 배포가 가능하다. 어떤 형태로든 배포할 수 있고 확장 가능하다.

# 모니터링

프로메테우스와 그라파나 대시보드를 지원하며, 경고 기능도 지원한다.
nats-top 같은 모니터링 도구를 제공해준다.
사이드카 배포도 가능하고, 간단한 연결-뷰 모델도 지원한다

# 관리

보안과 운영을 독립화한다.
운영중일때도 설정파일을 바꿀 수 있다.

# 통합

웹소켓, 카프카브리지, 레디스 커넥터,엘라스틱서치, 프로메테우스, .. HTTP, 등을 지원한다

# 보장성

<한번 보내고 잊어버림. 저장안한다>
At most once QoS: Core NATS offers an at most once quality of service. If a subscriber is not listening on the subject (no subject match), or is not active when the message is sent, the message is not received. This is the same level of guarantee that TCP/IP provides. Core NATS is a fire-and-forget messaging system. It will only hold messages in memory and will never write messages directly to disk.

<제트스트림 기능을 활성화를 화면, 저장 가능. 이 방법말고도 ack, seq 번호로도 같은 기능 구현 가능>
At-least / exactly once QoS: If you need higher qualities of service (at least once and exactly once), or functionalities such as persistent streaming, de-coupled flow control, and Key/Value Store, you can use NATS JetStream, which is built in to the NATS server (but needs to be enabled). Of course, you can also always build additional reliability into your client applications yourself with proven and scalable reference designs such as acks and sequence numbers.

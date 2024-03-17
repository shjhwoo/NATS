import {
  connect,
  StringCodec,
  Empty,
} from "https://cdn.jsdelivr.net/npm/nats.ws@1.10.0/esm/nats.js";

let natClient;
let sc = new StringCodec();

$(function () {
  initialize();

  //로그인
  $(document).on("click", ".login", function (e) {
    sendLoginRequest($("input").val());
  });

  $(document).on("keyup", "input", function (e) {
    if (e.keyCode && e.keyCode == 13) {
      sendLoginRequest($("input").val());
    }
  });
});

async function initialize() {
  //서버측에서는 모든메세지 가져와라는 플래그를 구독중, 해당 요청을 받아서 스트림에 대한 컨슈머를 생성하고
  //메세지를 다 빼오고, 응답으로 보내준다.
  await ConnectNats();

  console.log("init: ", location.href);
}

async function ConnectNats() {
  console.log("Connecting to NATS server");
  natClient = await connect({
    servers: ["ws://localhost:8080"],
  });
  console.log("Connected to " + natClient.getServer());
}

function sendLoginRequest(username) {
  console.log("login request", username);
  natClient
    .request(`login`, sc.encode(username), { timeout: 5000 })
    .then((m) => {
      console.log("login success", m);

      const lastAccess = sc.decode(m.data);
      document.cookie = `NATSLOGIN_${username}=${lastAccess}`;
      location.href = "/client4/main.html";
    })
    .catch((err) => {
      console.error(err);
      console.log(`problem with request: ${err.message}`);
    });
}

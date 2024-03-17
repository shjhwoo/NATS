import {
  connect,
  StringCodec,
  Empty,
} from "https://cdn.jsdelivr.net/npm/nats.ws@1.10.0/esm/nats.js";

let natClient;
let sc = new StringCodec();

$(function () {
  initialize();
});

const myId = document.cookie
  .split(";")
  .filter((cookie) => {
    return cookie.includes("NATSLOGIN_");
  })[0]
  .split("=")[0]
  .split("_")[1];

async function initialize() {
  //서버측에서는 모든메세지 가져와라는 플래그를 구독중, 해당 요청을 받아서 스트림에 대한 컨슈머를 생성하고
  //메세지를 다 빼오고, 응답으로 보내준다.
  await ConnectNats();

  console.log("init: ", location.href);

  initChatRoom();

  //채팅방
  $(document).on("click", "button.msg", function (e) {
    e.preventDefault();
    const message = $("textarea").val();
    console.log(message);
    sendMessage(message);
    $("textarea").val("");
  });

  $(document).on("keyup", "textarea", function (e) {
    if (e.keyCode && e.keyCode == 13) {
      const message = $(this).val();
      console.log(message);
      sendMessage(message);
      $("textarea").val("");
    }
  });

  function sendMessage(message) {
    natClient.publish(
      `msg.${myId}`,
      sc.encode(`{"user":"${myId}","text":"${message.trimEnd()}"}`)
    );
  }
}

async function ConnectNats() {
  console.log("Connecting to NATS server");
  natClient = await connect({
    servers: ["ws://localhost:8080"],
  });
  console.log("Connected to " + natClient.getServer());
}

function initChatRoom() {
  natClient.publish(
    `msg.${myId}`,
    sc.encode(`{"user":"${myId}","text":"connected"}`)
  );

  //모든 메세지를 가져와라는 요청을 publish
  natClient
    .request(`meta.ALL`, Empty, { timeout: 5000 }) //msg로 하니까 내용 못받고 스트림과 직통해서... 안되더라. 이름 다르게 해야 함!!반드시!!
    .then((m) => {
      console.log(sc.decode(m.data));
      const messageList = JSON.parse(sc.decode(m.data));

      $("span").text(messageList.length);

      for (let i = 0; i < messageList.length; i++) {
        let message = messageList[i].trim();
        if (message == "") continue;

        message = JSON.parse(message);

        if (message.user != myId) {
          $("ul").append(
            `<li class="youMsg">${
              message.text == "connected"
                ? `${message.user}님이 입장하셨습니다`
                : `${message.user}: ${message.text}`
            }</li>`
          );
        } else {
          if (message.text != "connected") {
            $("ul").append(`<li class="meMsg">${message.text}</li>`);
          }
        }
      }
    })
    .catch((err) => {
      console.log(`problem with request: ${err.message}`);
    });

  const sub = natClient.subscribe("msg.>");
  (async () => {
    for await (const msg of sub) handle(msg);
  })();

  const handle = (msg) => {
    console.log(sc.decode(msg.data));
    const message = JSON.parse(sc.decode(msg.data));

    if (message.user != myId) {
      $("ul").append(
        `<li class="youMsg">${
          message.text == "connected"
            ? `${message.user}님이 입장하셨습니다`
            : `${message.user}: ${message.text}`
        }</li>`
      );
    } else {
      if (message.text != "connected") {
        $("ul").append(`<li class="meMsg">${message.text}</li>`);
      }
    }
  };
}

// Finally drain the connection which will handle any outstanding
// messages before closing the connection.
//natClient.drain(); 쓰지마

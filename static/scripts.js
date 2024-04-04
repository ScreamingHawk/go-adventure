const addMessage = (msg, isUser) => {
  const msgList = document.getElementById("msgList");
  const p = document.createElement("p");
  p.className = isUser ? "msg-user" : "msg-ai";
  p.textContent = msg;
  msgList.appendChild(p);
};

document.addEventListener("DOMContentLoaded", () => {
  const msgForm = document.getElementById("msgForm");
  const msgBox = document.getElementById("msgBox");
  const sendBtn = document.getElementById("sendBtn");

  msgForm.addEventListener("submit", (event) => {
    event.preventDefault();
    const msg = msgBox.value;
    sendBtn.disabled = true;
    msgBox.value = "";

    addMessage(msg, true);

    fetch("/send", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ message: msg }),
    })
      .then((response) => {
        if (response.ok) {
          return response.json();
        }
        throw new Error("Request failed.");
      })
      .then((data) => {
        console.log(data);
        addMessage(data.message, false);
      })
      .finally(() => {
        sendBtn.disabled = false;
      });
  });
});
